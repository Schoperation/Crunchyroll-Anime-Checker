package postgres

import (
	"context"
	"schoperation/crunchyrollanimestatus/domain/anime"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
)

type PosterDao struct {
	db *pgx.Conn
}

func NewPosterDao(db *pgx.Conn) PosterDao {
	return PosterDao{
		db: db,
	}
}

type posterModel struct {
	AnimeId     int    `db:"anime_id"`
	ImageTypeId int    `db:"image_type_id"`
	Url         string `db:"url"`
	Encoded     string `db:"encoded"`
}

func (dao PosterDao) GetAllByAnimeId(animeId int) ([]anime.ImageDto, error) {
	sql, args, err := goqu.
		Select(&posterModel{}).
		From("poster").
		Where(
			goqu.C("anime_id").Eq(animeId),
		).
		WithDialect("postgres").
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, sqlBuilderError("poster", err)
	}

	rows, err := dao.db.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, couldNotRetrieveError("poster", err)
	}

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[posterModel])
	if err != nil {
		return nil, couldNotRetrieveError("poster", err)
	}

	dtos := make([]anime.ImageDto, len(models))
	for i, model := range models {
		dtos[i] = anime.ImageDto{
			AnimeId:       model.AnimeId,
			ImageType:     model.ImageTypeId,
			SeasonNumber:  0,
			EpisodeNumber: 0,
			Url:           model.Url,
			Encoded:       model.Encoded,
		}
	}

	return dtos, nil
}

func (dao PosterDao) InsertAll(dtos []anime.ImageDto) error {
	if len(dtos) == 0 {
		return nil
	}

	models := make([]posterModel, len(dtos))
	for i, dto := range dtos {
		models[i] = dao.imageDtoToModel(dto)
	}

	sql, args, err := goqu.
		Insert("poster").
		Rows(models).
		WithDialect("postgres").
		Prepared(false).
		ToSQL()
	if err != nil {
		return sqlBuilderError("poster", err)
	}

	_, err = dao.db.Exec(context.Background(), sql, args...)
	if err != nil {
		return couldNotUpdateError("poster", err)
	}

	return nil
}

func (dao PosterDao) Update(dto anime.ImageDto) error {
	sql, args, err := goqu.
		Update("poster").
		Set(dao.imageDtoToModel(dto)).
		Where(
			goqu.C("anime_id").Eq(dto.AnimeId),
			goqu.C("image_type_id").Eq(dto.ImageType),
		).
		WithDialect("postgres").
		Prepared(false).
		ToSQL()
	if err != nil {
		return sqlBuilderError("poster", err)
	}

	_, err = dao.db.Exec(context.Background(), sql, args...)
	if err != nil {
		return couldNotCreateError("poster", err)
	}

	return nil
}

func (dao PosterDao) imageDtoToModel(dto anime.ImageDto) posterModel {
	return posterModel{
		AnimeId:     dto.AnimeId,
		ImageTypeId: dto.ImageType,
		Url:         dto.Url,
		Encoded:     dto.Encoded,
	}
}
