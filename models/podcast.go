package models

type Podcast struct {
	ID          int64  `db:"id"`
	Author      User   `db:"author"`
	Picture     string `db:"picture"`
	Title       string `db:"title"`
	Description string `db:"description"`
	File        string `db:"file"`
	CreatedAt   int64  `db:"created_at"`
	UpdatedAt   int64  `db:"updated_at"`
}
