package sqlite

import (
	"schoperation/crunchyrollanimestatus/domain/anime"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
)

type PosterDao struct {
	db *goqu.Database
}

func NewPosterDao(db *goqu.Database) PosterDao {
	return PosterDao{
		db: db,
	}
}

type posterModel struct {
	AnimeId     int    `db:"anime_id" goqu:"skipupdate"`
	ImageTypeId int    `db:"image_type_id"`
	Url         string `db:"url"`
	Encoded     string `db:"encoded"`
}

func (dao PosterDao) GetAllByAnimeIds(animeIds []int) ([]anime.ImageDto, error) {
	var models []posterModel
	err := dao.db.
		Select(&posterModel{}).
		From("poster").
		Where(
			goqu.C("anime_id").In(animeIds),
		).
		WithDialect(Dialect).
		Prepared(true).
		Executor().
		ScanStructs(&models)
	if err != nil {
		return nil, couldNotRetrieveError("poster", err)
	}

	if len(models) < len(animeIds) {
		return nil, couldNotRetrieveAllError("poster", len(animeIds), len(models))
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

	_, err := dao.db.
		Insert("poster").
		Rows(models).
		OnConflict(goqu.DoNothing()).
		WithDialect(Dialect).
		Prepared(false).
		Executor().
		Exec()
	if err != nil {
		return couldNotCreateError("poster", err)
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
		WithDialect(Dialect).
		Prepared(false).
		ToSQL()
	if err != nil {
		return sqlBuilderError("poster", err)
	}

	_, err = dao.db.Exec(sql, args...)
	if err != nil {
		return couldNotUpdateError("poster", err)
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
