package gofileAPI

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	const baseURL = "https://api.gofile.io/contents/"

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
	const baseURL = "https://api.gofile.io/contents/"
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
	const baseURL = "https://api.gofile.io/contents/"

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
