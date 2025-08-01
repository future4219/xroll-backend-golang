package entity

import "time"

type Thread struct {
	ID        string
	Title     string
	LikeCount int
	Comments  []ThreadComment
	CreatedAt time.Time
}

type ThreadComment struct {
	ID         string
	ThreadID   string
	ThreaderID string
	Comment    string
	LikeCount  int
	CreatedAt  time.Time
}
