package output_port

import (
	"context"
	"net/http"
)

type Twitter interface {
	GetVideoByURL(tweetURL string) (string, error)
	GetThumbnailByURL(tweetURL string) (string, error)
	FetchTwimgStream(ctx context.Context, twimgURL string) (string, *http.Response, error)
}
