package models

type StarredPodcast struct {
	ID        int64
	Podcast   Podcast
	User      User
	CreatedAt int64
}
