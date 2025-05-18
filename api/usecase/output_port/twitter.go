package output_port

type Twitter interface {
	GetVideoByURL(tweetURL string) (string, error)
	GetThumbnailByURL(tweetURL string) (string, error)
}
