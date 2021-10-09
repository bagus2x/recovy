package models

type Webinar struct {
	ID          int64
	Author      User
	Picture     string
	Title       string
	Description string
	Category    string
	StartDate   int64
	LastDate    int64
	Time        string
	CreatedAt   int64
	UpdatedAt   int64
}
