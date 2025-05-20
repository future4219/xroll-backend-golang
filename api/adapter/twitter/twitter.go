package twitter

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"

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
	form := url.Values{}
	form.Set("test", tweetURL)

	req, err := http.NewRequest("POST", "https://awakest.net/twitter-video-downloader/", bytes.NewBufferString(form.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Referer", "https://awakest.net/twitter-video-downloader/")
	req.Header.Set("Origin", "https://awakest.net")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// .mp4のURLを探す
	re := regexp.MustCompile(`https://video\.twimg\.com/[^"'<>\\\s]+\.mp4`)
	matches := re.FindAllString(string(body), -1)

	if len(matches) == 0 {
		return "", fmt.Errorf("%w: twitter video", interactor.ErrKind.NotFound)
	}

	return matches[0], nil
}

func (t *Twitter) GetThumbnailByURL(tweetURL string) (string, error) {
	form := url.Values{}
	form.Set("page", tweetURL)
	form.Set("ftype", "all")
	form.Set("ajax", "1")

	req, err := http.NewRequest("POST", "https://twmate.com/", strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("リクエスト作成失敗: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Referer", "https://twmate.com/")
	req.Header.Set("Origin", "https://twmate.com")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("POST失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTPエラー: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("レスポンス読み込み失敗: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return "", fmt.Errorf("HTMLパース失敗: %w", err)
	}

	src, exists := doc.Find(".thumb-container img").First().Attr("src")
	if !exists {
		return "", fmt.Errorf("サムネイル画像が見つかりませんでした")
	}
	return src, nil
}
