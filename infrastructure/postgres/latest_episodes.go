package postgres

import (
	"context"
	"schoperation/crunchyrollanimestatus/domain/anime"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
)

type LatestEpisodesDao struct {
	db *pgx.Conn
}

func NewLatestEpisodesDao(db *pgx.Conn) LatestEpisodesDao {
	return LatestEpisodesDao{
		db: db,
	}
}

type latestEpisodesModel struct {
	AnimeId          int       `db:"anime_id" goqu:"skipupdate"`
	LocaleId         int       `db:"locale_id"`
	LastUpdated      time.Time `db:"last_updated"`
	LatestSubSeason  int       `db:"latest_sub_season"`
	LatestSubEpisode int       `db:"latest_sub_episode"`
	LatestSubTitle   string    `db:"latest_sub_title"`
	LatestDubSeason  int       `db:"latest_dub_season"`
	LatestDubEpisode int       `db:"latest_dub_episode"`
	LatestDubTitle   string    `db:"latest_dub_title"`
}

func (dao LatestEpisodesDao) GetAllByAnimeId(animeId int) ([]anime.LatestEpisodesDto, error) {
	sql, args, err := goqu.
		Select(&latestEpisodesModel{}).
		From("latest_episodes").
		Where(
			goqu.C("anime_id").Eq(animeId),
		).
		WithDialect("postgres").
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, sqlBuilderError("latest_episodes", err)
	}

	rows, err := dao.db.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, couldNotRetrieveError("latest_episodes", err)
	}

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[latestEpisodesModel])
	if err != nil {
		return nil, couldNotRetrieveError("latest_episodes", err)
	}

	dtos := make([]anime.LatestEpisodesDto, len(models))
	for i, model := range models {
		dtos[i] = anime.LatestEpisodesDto{
			AnimeId:          model.AnimeId,
			LocaleId:         model.LocaleId,
			LastUpdated:      model.LastUpdated,
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
