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
	AnimeId          int       `db:"anime_id" goqu:"skipinsert"`
	SeriesId         string    `db:"series_id"`
	SlugTitle        string    `db:"slug_title"`
	Title            string    `db:"title"`
	LastUpdated      time.Time `db:"last_updated"`
	SeasonIdentifier string    `db:"season_identifier"`
}

func (dao AnimeDao) GetAllMinimal() ([]anime.MinimalAnimeDto, error) {
	sql, args, err := goqu.Select(&minimalAnimeModel{}).From("anime").WithDialect("postgres").Prepared(true).ToSQL()
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

func (dao AnimeDao) InsertAll(dtos []anime.AnimeDto) error {
	sql, args, err := goqu.Insert("anime").Rows(dao.dtosToModels(dtos)).WithDialect("postgres").Prepared(false).ToSQL()
	if err != nil {
		return sqlBuilderError("anime", err)
	}

	_, err = dao.db.Exec(context.Background(), sql, args...)
	if err != nil {
		return couldNotCreateError("anime", err)
	}

	return nil
}

func (dao AnimeDao) dtosToModels(dtos []anime.AnimeDto) []animeModel {
	models := make([]animeModel, len(dtos))
	for i, dto := range dtos {
		models[i] = animeModel{
			AnimeId:          dto.AnimeId,
			SeriesId:         dto.SeriesId,
			SlugTitle:        dto.SlugTitle,
			Title:            dto.Title,
			LastUpdated:      dto.LastUpdated,
			SeasonIdentifier: dto.SeasonIdentifier,
		}
	}

	return models
}
