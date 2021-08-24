package db

import (
	"github.com/bagus2x/recovy/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func OpenPostgres(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", cfg.DatabaseConnection())
	if err != nil {
		return nil, err
	}

	return db, nil
}
