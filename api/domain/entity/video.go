package entity

import "time"

type Video struct {
	ID            string
	Ranking       int
	VideoURL      string
	ThumbnailURL  string
	TweetURL      *string
	DownloadCount int
	LikeCount     int
	Comments       []Comment
	CreatedAt     time.Time
}

type Comment struct {
	ID        string
	Comment   string
	LikeCount int
	CreatedAt time.Time
}
