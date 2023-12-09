package postgres

import (
	"context"
	"schoperation/crunchyrollanimestatus/domain/anime"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
)

type ThumbnailDao struct {
	db *pgx.Conn
}

func NewThumbnailDao(db *pgx.Conn) ThumbnailDao {
	return ThumbnailDao{
		db: db,
	}
}

type thumbnailModel struct {
	AnimeId       int    `db:"anime_id"`
	SeasonNumber  int    `db:"season_number"`
	EpisodeNumber int    `db:"episode_number"`
	Url           string `db:"url"`
	Encoded       string `db:"encoded"`
}

func (dao ThumbnailDao) GetAllByAnimeId(animeId int) ([]anime.ImageDto, error) {
	sql, args, err := goqu.
		Select(&thumbnailModel{}).
		From("thumbnail").
		Where(
			goqu.C("anime_id").Eq(animeId),
		).
		WithDialect("postgres").
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, sqlBuilderError("thumbnail", err)
	}

	rows, err := dao.db.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, couldNotRetrieveError("thumbnail", err)
	}

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[thumbnailModel])
	if err != nil {
		return nil, couldNotRetrieveError("thumbnail", err)
	}

	dtos := make([]anime.ImageDto, len(models))
	for i, model := range models {
		dtos[i] = anime.ImageDto{
			AnimeId:       model.AnimeId,
			ImageType:     anime.ImageTypeThumbnail.Int(),
			SeasonNumber:  model.SeasonNumber,
			EpisodeNumber: model.SeasonNumber,
			Url:           model.Url,
			Encoded:       model.Encoded,
		}
	}

	return dtos, nil
}

func (dao ThumbnailDao) InsertAll(dtos []anime.ImageDto) error {
	if len(dtos) == 0 {
		return nil
	}

	models := make([]thumbnailModel, len(dtos))
	for i, dto := range dtos {
		models[i] = dao.imageDtoToModel(dto)
	}

	sql, args, err := goqu.
		Insert("thumbnail").
		Rows(models).
		WithDialect("postgres").
		Prepared(false).
		ToSQL()
	if err != nil {
		return sqlBuilderError("thumbnail", err)
	}

	_, err = dao.db.Exec(context.Background(), sql, args...)
	if err != nil {
		return couldNotCreateError("thumbnail", err)
	}

	return nil
}

func (dao ThumbnailDao) Delete(dto anime.ImageDto) error {
	sql, args, err := goqu.
		Delete("thumbnail").
		Where(
			goqu.C("anime_id").Eq(dto.AnimeId),
			goqu.C("season_number").Eq(dto.SeasonNumber),
			goqu.C("episode_number").Eq(dto.EpisodeNumber),
		).
		WithDialect("postgres").
		Prepared(false).
		ToSQL()
	if err != nil {
		return sqlBuilderError("thumbnail", err)
	}

	_, err = dao.db.Exec(context.Background(), sql, args...)
	if err != nil {
		return couldNotDeleteError("thumbnail", err)
	}

	return nil
}

func (dao ThumbnailDao) imageDtoToModel(dto anime.ImageDto) thumbnailModel {
	return thumbnailModel{
		AnimeId:       dto.AnimeId,
		SeasonNumber:  dto.SeasonNumber,
		EpisodeNumber: dto.EpisodeNumber,
		Url:           dto.Url,
		Encoded:       dto.Encoded,
	}
}
