package sqlite

import (
	"schoperation/crunchyrollanimestatus/domain/anime"
	"schoperation/crunchyrollanimestatus/domain/core"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
)

type ThumbnailDao struct {
	db *goqu.Database
}

func NewThumbnailDao(db *goqu.Database) ThumbnailDao {
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

func (dao ThumbnailDao) GetAllByAnimeIds(animeIds []int) ([]anime.ImageDto, error) {
	var models []thumbnailModel
	err := dao.db.
		Select(&thumbnailModel{}).
		From("thumbnail").
		Where(
			goqu.C("anime_id").In(animeIds),
		).
		WithDialect(Dialect).
		Prepared(true).
		Executor().
		ScanStructs(&models)
	if err != nil {
		return nil, couldNotRetrieveError("thumbnail", err)
	}

	if len(models) < len(animeIds) {
		return nil, couldNotRetrieveAllError("thumbnail", len(animeIds), len(models))
	}

	dtos := make([]anime.ImageDto, len(models))
	for i, model := range models {
		dtos[i] = anime.ImageDto{
			AnimeId:       model.AnimeId,
			ImageType:     core.ImageTypeThumbnail.Int(),
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

	_, err := dao.db.
		Insert("thumbnail").
		Rows(models).
		WithDialect(Dialect).
		Prepared(false).
		Executor().
		Exec()
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
		WithDialect(Dialect).
		Prepared(false).
		ToSQL()
	if err != nil {
		return sqlBuilderError("thumbnail", err)
	}

	_, err = dao.db.Exec(sql, args...)
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
