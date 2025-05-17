package input_port

type ITwitterUseCase interface {
	GetVideoByURL(tweetURL string) (string, error)
}
