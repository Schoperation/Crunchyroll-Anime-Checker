package postgres

import (
	"github.com/jackc/pgx/v5"
)

type AnimeDao struct {
	conn *pgx.Conn
}

func NewAnimeDao(conn *pgx.Conn) AnimeDao {
	return AnimeDao{
		conn: conn,
	}
}
