package sqlite

import (
	"database/sql"
	"schoperation/crunchyrollanimestatus/domain/anime"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
)

type AnimeDao struct {
	db *sql.DB
}

func NewAnimeDao(db *sql.DB) AnimeDao {
	return AnimeDao{
		db: db,
	}
}

type minimalAnimeModel struct {
	AnimeId     int       `db:"anime_id"`
	SeriesId    string    `db:"series_id"`
	LastUpdated time.Time `db:"last_updated"`
}

type animeModel struct {
	AnimeId     int       `db:"anime_id" goqu:"skipinsert,skipupdate"`
	SeriesId    string    `db:"series_id"`
	SlugTitle   string    `db:"slug_title"`
	Title       string    `db:"title"`
	IsSimulcast bool      `db:"is_simulcast"`
	LastUpdated time.Time `db:"last_updated"`
}

func (dao AnimeDao) GetAllMinimal() ([]anime.MinimalAnimeDto, error) {
	sql, args, err := goqu.
		Select(&minimalAnimeModel{}).
		From("anime").
		WithDialect(GoquDialect).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, sqlBuilderError("anime", err)
	}

	rows, err := dao.db.Query(sql, args...)
	if err != nil {
		return nil, couldNotRetrieveError("anime", err)
	}

	models, err := scanRows[animeModel](rows)
	if err != nil {
		return nil, couldNotRetrieveError("anime", err)
	}

	dtos := make([]anime.MinimalAnimeDto, len(models))
	for i, model := range models {
		dtos[i] = anime.MinimalAnimeDto{
			AnimeId:     model.AnimeId,
			SeriesId:    model.SeriesId,
			LastUpdated: model.LastUpdated,
		}
	}

	return dtos, nil
}

func (dao AnimeDao) GetAllByAnimeIds(animeIds []int) ([]anime.AnimeDto, error) {
	if len(animeIds) == 0 {
		return nil, nil
	}

	sql, args, err := goqu.
		Select(&animeModel{}).
		From("anime").
		Where(
			goqu.C("anime_id").In(animeIds),
		).
		WithDialect(GoquDialect).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, sqlBuilderError("anime", err)
	}

	rows, err := dao.db.Query(sql, args...)
	if err != nil {
		return nil, couldNotRetrieveError("anime", err)
	}
	defer rows.Close()

	models, err := scanRows[animeModel](rows)
	if err != nil {
		return nil, couldNotRetrieveError("anime", err)
	}

	if len(animeIds) != len(models) {
		return nil, couldNotRetrieveAllError("anime", len(animeIds), len(models))
	}

	dtos := make([]anime.AnimeDto, len(models))
	for i, model := range models {
		dtos[i] = anime.AnimeDto{
			AnimeId:     model.AnimeId,
			SeriesId:    model.SeriesId,
			SlugTitle:   model.SlugTitle,
			Title:       model.Title,
			IsSimulcast: model.IsSimulcast,
			LastUpdated: model.LastUpdated,
		}
	}

	return dtos, nil
}

func (dao AnimeDao) InsertAll(dtos []anime.AnimeDto) ([]anime.MinimalAnimeDto, error) {
	if len(dtos) == 0 {
		return nil, nil
	}

	models := make([]animeModel, len(dtos))
	for i, dto := range dtos {
		models[i] = dao.animeDtoToModel(dto)
	}

	sql, args, err := goqu.
		Insert("anime").
		Rows(models).
		Returning("anime_id", "series_id", "last_updated").
		WithDialect(GoquDialect).
		Prepared(false).
		ToSQL()
	if err != nil {
		return nil, sqlBuilderError("anime", err)
	}

	rows, err := dao.db.Query(sql, args...)
	if err != nil {
		return nil, couldNotCreateError("anime", err)
	}
	defer rows.Close()

	minimalModels, err := scanRows[minimalAnimeModel](rows)
	if err != nil {
		return nil, couldNotRetrieveError("minimal anime", err)
	}

	if len(minimalModels) != len(dtos) {
		return nil, couldNotRetrieveAllError("minimal anime", len(dtos), len(minimalModels))
	}

	minimalDtos := make([]anime.MinimalAnimeDto, len(minimalModels))
	for i, model := range minimalModels {
		minimalDtos[i] = anime.MinimalAnimeDto{
			AnimeId:     model.AnimeId,
			SeriesId:    model.SeriesId,
			LastUpdated: model.LastUpdated,
		}
	}

	return minimalDtos, nil
}

func (dao AnimeDao) Update(dto anime.AnimeDto) error {
	sql, args, err := goqu.
		Update("anime").
		Set(dao.animeDtoToModel(dto)).
		Where(
			goqu.C("anime_id").Eq(dto.AnimeId),
		).
		WithDialect(GoquDialect).
		Prepared(false).
		ToSQL()
	if err != nil {
		return sqlBuilderError("anime", err)
	}

	_, err = dao.db.Exec(sql, args...)
	if err != nil {
		return couldNotUpdateError("anime", err)
	}

	return nil
}

func (dao AnimeDao) animeDtoToModel(dto anime.AnimeDto) animeModel {
	return animeModel{
		AnimeId:     dto.AnimeId,
		SeriesId:    dto.SeriesId,
		SlugTitle:   dto.SlugTitle,
		Title:       dto.Title,
		IsSimulcast: dto.IsSimulcast,
		LastUpdated: dto.LastUpdated,
	}
}
