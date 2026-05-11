package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgres(dbURL string) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", dbURL)
}