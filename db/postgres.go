package db

import (
	"database/sql"

	"github.com/bagus2x/recovy/config"
	_ "github.com/lib/pq"
)

func OpenPostgres(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DatabaseConnection())
	if err != nil {
		return nil, err
	}

	return db, nil
}
