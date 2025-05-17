package twitter

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"io/ioutil"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/interactor"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"
)

type Twitter struct{}

func NewTwitter() output_port.Twitter {
	return &Twitter{}
}

// GetVideoURL は tweet の URL から Twidropper を使って動画URLを抽出する
func (t *Twitter) GetVideoByURL(tweetURL string) (string, error) {
	endpoint := os.Getenv("TWIDROPPER_ENDPOINT")
	if endpoint == "" {
		return "", fmt.Errorf("TWIDROPPER_ENDPOINT が設定されていません")
	}

	data := url.Values{}
	data.Set("url", tweetURL)
	data.Set("submitBtn", "submit")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.PostForm(endpoint, data)
	if err != nil {
		return "", fmt.Errorf("POST失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTPステータスエラー: %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("レスポンス読み込み失敗: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyBytes)))
	if err != nil {
		return "", fmt.Errorf("HTML解析失敗: %w", err)
	}

	videoSrc, exists := doc.Find("video").Attr("src")
	if !exists {
		return "", fmt.Errorf("%w: twitter video", interactor.ErrKind.NotFound)
	}

	return videoSrc, nil
}
