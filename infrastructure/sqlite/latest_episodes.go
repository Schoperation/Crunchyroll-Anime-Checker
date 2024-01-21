package sqlite

import (
	"schoperation/crunchyrollanimestatus/domain/anime"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
)

type LatestEpisodesDao struct {
	db *goqu.Database
}

func NewLatestEpisodesDao(db *goqu.Database) LatestEpisodesDao {
	return LatestEpisodesDao{
		db: db,
	}
}

type latestEpisodesModel struct {
	AnimeId          int    `db:"anime_id" goqu:"skipupdate"`
	LocaleId         int    `db:"locale_id"`
	LatestSubSeason  int    `db:"latest_sub_season"`
	LatestSubEpisode int    `db:"latest_sub_episode"`
	LatestSubTitle   string `db:"latest_sub_title"`
	LatestDubSeason  int    `db:"latest_dub_season"`
	LatestDubEpisode int    `db:"latest_dub_episode"`
	LatestDubTitle   string `db:"latest_dub_title"`
}

func (dao LatestEpisodesDao) GetAllByAnimeIds(animeIds []int) ([]anime.LatestEpisodesDto, error) {
	var models []latestEpisodesModel
	err := dao.db.
		Select(&latestEpisodesModel{}).
		From("latest_episodes").
		Where(
			goqu.C("anime_id").In(animeIds),
		).
		WithDialect(Dialect).
		Prepared(true).
		Executor().
		ScanStructs(&models)
	if err != nil {
		return nil, couldNotRetrieveError("latest_episodes", err)
	}

	if len(models) < len(animeIds) {
		return nil, couldNotRetrieveAllError("latest_episodes", len(animeIds), len(models))
	}

	dtos := make([]anime.LatestEpisodesDto, len(models))
	for i, model := range models {
		dtos[i] = anime.LatestEpisodesDto{
			AnimeId:          model.AnimeId,
			LocaleId:         model.LocaleId,
			LatestSubSeason:  model.LatestSubSeason,
			LatestSubEpisode: model.LatestSubEpisode,
			LatestSubTitle:   model.LatestSubTitle,
			LatestDubSeason:  model.LatestDubSeason,
			LatestDubEpisode: model.LatestDubEpisode,
			LatestDubTitle:   model.LatestDubTitle,
		}
	}

	return dtos, nil
}

func (dao LatestEpisodesDao) InsertAll(dtos []anime.LatestEpisodesDto) error {
	if len(dtos) == 0 {
		return nil
	}

	models := make([]latestEpisodesModel, len(dtos))
	for i, dto := range dtos {
		models[i] = dao.latestEpisodesDtoToModel(dto)
	}

	_, err := dao.db.
		Insert("latest_episodes").
		Rows(models).
		WithDialect(Dialect).
		Prepared(false).
		Executor().
		Exec()
	if err != nil {
		return couldNotCreateError("latest_episodes", err)
	}

	return nil
}

func (dao LatestEpisodesDao) latestEpisodesDtoToModel(dto anime.LatestEpisodesDto) latestEpisodesModel {
	return latestEpisodesModel{
		AnimeId:          dto.AnimeId,
		LocaleId:         dto.LocaleId,
		LatestSubSeason:  dto.LatestSubSeason,
		LatestSubEpisode: dto.LatestSubEpisode,
		LatestSubTitle:   dto.LatestSubTitle,
		LatestDubSeason:  dto.LatestDubSeason,
		LatestDubEpisode: dto.LatestDubEpisode,
		LatestDubTitle:   dto.LatestDubTitle,
	}
}
