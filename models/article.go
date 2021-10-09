package models

type Article struct {
	ID          int64
	Author      User
	Picture     string
	Title       string
	Description string
	Category    string
	CreatedAt   int64
	UpdatedAt   int64
}
