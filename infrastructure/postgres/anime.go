package postgres

import (
	"context"
	"schoperation/crunchyrollanimestatus/domain/anime"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
)

type AnimeDao struct {
	db *pgx.Conn
}

func NewAnimeDao(db *pgx.Conn) AnimeDao {
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
	AnimeId          int       `db:"anime_id" goqu:"skipinsert,skipupdate"`
	SeriesId         string    `db:"series_id"`
	SlugTitle        string    `db:"slug_title"`
	Title            string    `db:"title"`
	LastUpdated      time.Time `db:"last_updated"`
	SeasonIdentifier string    `db:"season_identifier"`
}

func (dao AnimeDao) GetAllMinimal() ([]anime.MinimalAnimeDto, error) {
	sql, args, err := goqu.
		Select(&minimalAnimeModel{}).
		From("anime").
		WithDialect("postgres").
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, sqlBuilderError("anime", err)
	}

	rows, err := dao.db.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, couldNotRetrieveError("anime", err)
	}

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[minimalAnimeModel])
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
		WithDialect("postgres").
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, sqlBuilderError("anime", err)
	}

	rows, err := dao.db.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, couldNotRetrieveError("anime", err)
	}

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[animeModel])
	if err != nil {
		return nil, couldNotRetrieveError("anime", err)
	}

	dtos := make([]anime.AnimeDto, len(models))
	for i, model := range models {
		dtos[i] = anime.AnimeDto{
			AnimeId:          model.AnimeId,
			SeriesId:         model.SeriesId,
			SlugTitle:        model.SlugTitle,
			Title:            model.Title,
			LastUpdated:      model.LastUpdated,
			SeasonIdentifier: model.SeasonIdentifier,
		}
	}

	return dtos, nil
}

func (dao AnimeDao) InsertAll(dtos []anime.AnimeDto) error {
	if len(dtos) == 0 {
		return nil
	}

	models := make([]animeModel, len(dtos))
	for i, dto := range dtos {
		models[i] = dao.animeDtoToModel(dto)
	}

	sql, args, err := goqu.
		Insert("anime").
		Rows(models).
		WithDialect("postgres").
		Prepared(false).
		ToSQL()
	if err != nil {
		return sqlBuilderError("anime", err)
	}

	_, err = dao.db.Exec(context.Background(), sql, args...)
	if err != nil {
		return couldNotCreateError("anime", err)
	}

	return nil
}

func (dao AnimeDao) Update(dto anime.AnimeDto) error {
	sql, args, err := goqu.
		Update("anime").
		Set(dao.animeDtoToModel(dto)).
		Where(
			goqu.C("anime_id").Eq(dto.AnimeId),
		).
		WithDialect("postgres").
		Prepared(false).
		ToSQL()
	if err != nil {
		return sqlBuilderError("anime", err)
	}

	_, err = dao.db.Exec(context.Background(), sql, args...)
	if err != nil {
		return couldNotUpdateError("anime", err)
	}

	return nil
}

func (dao AnimeDao) animeDtoToModel(dto anime.AnimeDto) animeModel {
	return animeModel{
		AnimeId:          dto.AnimeId,
		SeriesId:         dto.SeriesId,
		SlugTitle:        dto.SlugTitle,
		Title:            dto.Title,
		LastUpdated:      dto.LastUpdated,
		SeasonIdentifier: dto.SeasonIdentifier,
	}
}
