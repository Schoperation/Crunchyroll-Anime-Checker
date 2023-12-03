package postgres

import (
	"schoperation/crunchyrollanimestatus/domain/anime"
	"time"

	"github.com/doug-martin/goqu/v9"
)

type AnimeDao struct {
	db *goqu.Database
}

func NewAnimeDao(db *goqu.Database) AnimeDao {
	return AnimeDao{
		db: db,
	}
}

type minimalAnimeModel struct {
	AnimeId     int       `db:"anime_id"`
	SeriesId    string    `db:"series_id"`
	SlugTitle   string    `db:"slug_title"`
	LastUpdated time.Time `db:"last_updated"`
}

type animeModel struct {
}

func (dao AnimeDao) GetAllMinimal() (map[string]anime.MinimalAnimeDto, error) {
	sql, args, err := goqu.Select(&minimalAnimeModel{}).From("anime").Prepared(true).ToSQL()
	if err != nil {
		return nil, sqlBuilderError("anime")
	}

	var models []minimalAnimeModel
	err = dao.db.ScanStructs(&models, sql, args)
	if err != nil {
		return nil, couldNotRetrieveError("anime")
	}

	dtos := make(map[string]anime.MinimalAnimeDto, len(models))
	for _, model := range models {
		dtos[model.SeriesId] = anime.MinimalAnimeDto{
			AnimeId:     model.AnimeId,
			SeriesId:    model.SeriesId,
			SlugTitle:   model.SlugTitle,
			LastUpdated: model.LastUpdated,
		}
	}

	return dtos, nil
}
