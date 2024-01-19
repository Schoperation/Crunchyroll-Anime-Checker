package sqlite

import (
	"database/sql"
	"schoperation/crunchyrollanimestatus/domain/anime"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
)

type LatestEpisodesDao struct {
	db *sql.DB
}

func NewLatestEpisodesDao(db *sql.DB) LatestEpisodesDao {
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

func (dao LatestEpisodesDao) GetAllByAnimeId(animeId int) ([]anime.LatestEpisodesDto, error) {
	sql, args, err := goqu.
		Select(&latestEpisodesModel{}).
		From("latest_episodes").
		Where(
			goqu.C("anime_id").Eq(animeId),
		).
		WithDialect(GoquDialect).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, sqlBuilderError("latest_episodes", err)
	}

	rows, err := dao.db.Query(sql, args...)
	if err != nil {
		return nil, couldNotRetrieveError("latest_episodes", err)
	}

	models, err := scanRows[latestEpisodesModel](rows)
	if err != nil {
		return nil, couldNotRetrieveError("latest_episodes", err)
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

	sql, args, err := goqu.
		Insert("latest_episodes").
		Rows(models).
		WithDialect(GoquDialect).
		Prepared(false).
		ToSQL()
	if err != nil {
		return sqlBuilderError("latest_episodes", err)
	}

	_, err = dao.db.Exec(sql, args...)
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
