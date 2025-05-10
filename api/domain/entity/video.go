package entity

type Video struct {
	ID            string
	Ranking       int
	VideoURL      string
	ThumbnailURL  string
	TweetURL      string
	DownloadCount int
	LikeCount     int
	Comment       []Comment
	CreatedAt     string
}

type Comment struct {
	ID string
	Comment string
	likeCount int
	createdAt string
}