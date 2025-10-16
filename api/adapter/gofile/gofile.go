package gofileAPI

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"
)

type GofileAPIDriver struct{}

func NewGofileAPI() output_port.GofileAPIDriver {
	return &GofileAPIDriver{}
}

// Gofileの https://api.gofile.io/contents/{contentId} を使う
func (t *GofileAPIDriver) GetContent(gofileID string, gofileToken string) (output_port.GofileGetContentRes, error) {
	var zero output_port.GofileGetContentRes
	baseURL := os.Getenv("GOFILE_API_ENDPOINT") + "/contents/"

	// タイムアウト付きコンテキスト
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+gofileID, nil)
	if err != nil {
		return zero, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+gofileToken)

	// クライアント（必要ならリトライやTransport調整はここで）
	httpClient := &http.Client{Timeout: 60 * time.Second}

	resp, err := httpClient.Do(req)
	if err != nil {
		return zero, fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return zero, fmt.Errorf("gofile api status %d: %s", resp.StatusCode, string(b))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return zero, fmt.Errorf("read body: %w", err)
	}

	var out output_port.GofileGetContentRes
	if err := json.Unmarshal(body, &out); err != nil {
		return zero, fmt.Errorf("unmarshal: %w", err)
	}
	if out.Status != "ok" {
		return zero, fmt.Errorf("gofile status not ok: %s", out.Status)
	}

	return out, nil
}

// 直リンクを発行（POST /contents/{contentId}/directlinks）
// ※戻り値は単体の GofileDirectLink に変更
func (t *GofileAPIDriver) IssueDirectLink(contentID, gofileToken string) (output_port.GofileDirectLink, error) {
	baseURL := os.Getenv("GOFILE_API_ENDPOINT") + "/contents/"
	var zero output_port.GofileDirectLink

	// レスポンス用（data 直下がそのまま GofileDirectLink の形）
	type issueResp struct {
		Status string                       `json:"status"`
		Data   output_port.GofileDirectLink `json:"data"`
	}

	// ★ タイムアウト無し（ContextはBackground固定）
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, baseURL+contentID+"/directlinks", nil)
	if err != nil {
		return zero, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+gofileToken)

	// ★ クライアントのTimeoutもゼロ（＝無制限）
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return zero, fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return zero, fmt.Errorf("gofile api status %d: %s", resp.StatusCode, string(b))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return zero, fmt.Errorf("read body: %w", err)
	}

	var out issueResp
	if err := json.Unmarshal(body, &out); err != nil {
		return zero, fmt.Errorf("unmarshal: %w", err)
	}
	if out.Status != "ok" {
		return zero, fmt.Errorf("gofile status not ok: %s", out.Status)
	}
	return out.Data, nil
}

func (t *GofileAPIDriver) GetDirectLinks(contentID, gofileToken string) (map[string]output_port.GofileDirectLink, error) {
	baseURL := os.Getenv("GOFILE_API_ENDPOINT") + "/contents/"

	// レスポンス用ローカル型（必要十分だけ定義）
	type directLinksResp struct {
		Status string `json:"status"`
		Data   struct {
			DirectLinks map[string]output_port.GofileDirectLink `json:"directLinks"`
		} `json:"data"`
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+contentID+"/directlinks", nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+gofileToken)

	httpClient := &http.Client{Timeout: 20 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gofile api status %d: %s", resp.StatusCode, string(b))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	var out directLinksResp
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	if out.Status != "ok" {
		return nil, fmt.Errorf("gofile status not ok: %s", out.Status)
	}

	return out.Data.DirectLinks, nil
}

// 返却用：Gofileのアップロードレスポンス（必要十分のフィールド）
// Upload: 受け取った io.Reader を multipart/form-data でそのままGofileにストリーム転送
// - ctx: タイムアウト含むコンテキスト（呼び出し側で管理）
// - filename: アップロード時のファイル名
// - folderID: 同一フォルダに積みたい時に指定（空なら未指定）
// - r: 動画などの実データのストリーム
func (t *GofileAPIDriver) Upload(ctx context.Context, filename, folderID string, r io.Reader) (output_port.GofileUploadData, error) {
	var zero output_port.GofileUploadData

	endpoint := os.Getenv("GOFILE_UPLOAD_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://upload-ap-tyo.gofile.io/uploadfile"
	}
	const ua = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120 Safari/537.36"

	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)

	// multipart 生成は並行で書き込み
	go func() {
		defer func() {
			_ = mw.Close()
			_ = pw.Close()
		}()

		if folderID != "" {
			if err := mw.WriteField("folderId", folderID); err != nil {
				_ = pw.CloseWithError(fmt.Errorf("write folderId: %w", err))
				return
			}
		}
		part, err := mw.CreateFormFile("file", filename)
		if err != nil {
			_ = pw.CloseWithError(fmt.Errorf("create form file: %w", err))
			return
		}
		if _, err := io.Copy(part, r); err != nil {
			_ = pw.CloseWithError(fmt.Errorf("copy stream: %w", err))
			return
		}
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, pr)
	if err != nil {
		return zero, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("User-Agent", ua)

	cli := &http.Client{
		// タイムアウトは ctx 側で管理する想定（必要なら Transport 調整）
		Timeout: 0,
	}
	resp, err := cli.Do(req)
	if err != nil {
		return zero, fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return zero, fmt.Errorf("upload failed: %s, body=%s", resp.Status, string(body))
	}

	var up output_port.UploadResult
	if err := json.Unmarshal(body, &up); err != nil {
		return zero, fmt.Errorf("unmarshal upload: %w (body=%s)", err, string(body))
	}
	if up.Status != "ok" {
		return zero, fmt.Errorf("gofile status=%s body=%s", up.Status, string(body))
	}
	return up.Data, nil
}
