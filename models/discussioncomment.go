package models

type DiscussionComment struct {
	ID           int64
	DiscussionID int64
	Commentator  User
	Description  string
	CreatedAt    int64
	UpdatedAt    int64
}
