package twitter

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
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

func (t *Twitter) FetchTwimgStream(ctx context.Context, srcURL string) (string, *http.Response, error) {
	if srcURL == "" {
		return "", nil, fmt.Errorf("url is required")
	}

	uParsed, err := url.Parse(srcURL)
	if err != nil {
		return "", nil, fmt.Errorf("invalid url: %w", err)
	}
	if !strings.EqualFold(uParsed.Host, "video.twimg.com") {
		return "", nil, fmt.Errorf("unsupported host: %s (expect video.twimg.com)", uParsed.Host)
	}

	// ファイル名推定
	filename := path.Base(uParsed.Path)
	if filename == "" || filename == "." || filename == "/" {
		filename = "video.mp4"
	}
	if !strings.HasSuffix(strings.ToLower(filename), ".mp4") {
		filename += ".mp4"
	}

	const ua = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120 Safari/537.36"
	const referer = "https://twitter.com/"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, srcURL, nil)
	if err != nil {
		return "", nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Referer", referer)

	cli := &http.Client{
		Timeout: 0, // タイムアウトは ctx で管理
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.Header.Set("User-Agent", ua)
			r.Header.Set("Referer", referer)
			return nil
		},
	}

	resp, err := cli.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("download request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		_ = resp.Body.Close()
		return "", nil, fmt.Errorf("download failed: %s, body=%s", resp.Status, string(b))
	}

	return filename, resp, nil
}
