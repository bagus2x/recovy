package models

import "database/sql"

type User struct {
	ID        int64          `db:"id"`
	Name      string         `db:"name"`
	Email     string         `db:"email"`
	Picture   sql.NullString `db:"picture"`
	Password  string         `db:"password"`
	CreatedAt int64          `db:"created_at"`
	UpdatedAt int64          `db:"updated_at"`
	DeletedAt sql.NullInt64  `db:"deleted_at"`
}
