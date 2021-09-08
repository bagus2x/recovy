package models

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64         `db:"id"`
	Name      string        `db:"name"`
	Email     string        `db:"email"`
	Picture   string        `db:"picture"`
	Password  string        `db:"password"`
	CreatedAt int64         `db:"created_at"`
	UpdatedAt int64         `db:"updated_at"`
	DeletedAt sql.NullInt64 `db:"deleted_at"`
}

func (p *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.Password = string(bytes)

	return nil
}

func (p *User) ComparePasswords(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(password))
	return err == nil
}
