package handler

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo/v4"
)

type GofileHandler struct{}

func NewGofileHandler() *GofileHandler {
	return &GofileHandler{}
}

func (g *GofileHandler) ProxyGofileVideo(c echo.Context) error {
	rawURL := c.QueryParam("url")
	if rawURL == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "URLパラメータが必要です (?url=...)")
	}

	decodedURL, err := url.QueryUnescape(rawURL)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "URLのデコードに失敗しました")
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(decodedURL)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Gofileから動画取得に失敗しました: "+err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return echo.NewHTTPError(http.StatusBadGateway, "Gofileからのレスポンスが異常です")
	}

	// Content-Typeが無い場合のフォールバック
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "video/mp4"
	}

	// ブラウザでインライン再生用のヘッダ
	c.Response().Header().Set("Content-Type", contentType)
	c.Response().Header().Set("Content-Disposition", "inline")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().WriteHeader(http.StatusOK)

	// ストリーミング転送
	_, err = io.Copy(c.Response().Writer, resp.Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "動画転送中にエラーが発生しました: "+err.Error())
	}

	return nil
}
