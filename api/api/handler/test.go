package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

func ProxyGofileVideo(gofileDirectLink string) error {

	// URLをデコード
	decodedURL, err := url.QueryUnescape(gofileDirectLink)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "URLのデコードに失敗しました")
	}

	// GETリクエスト送信
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(decodedURL)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Gofileから動画取得に失敗しました: "+err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return echo.NewHTTPError(http.StatusBadGateway, "Gofileからのレスポンスが異常です")
	}

	// ★ ファイル保存（例：./downloads/video.mp4）
	savePath := "./video.mp4" // 任意のパスに変更可
	outFile, err := os.Create(savePath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ファイル保存に失敗しました: "+err.Error())
	}
	defer outFile.Close()

	// ファイルに保存しながら、レスポンスにも流す
	// multiWriter := io.MultiWriter(outFile, c.Response().Writer)

	// c.Response().Header().Set("Content-Type", "video/mp4")
	// c.Response().Header().Set("Content-Disposition", "inline")
	// c.Response().Header().Set("Accept-Ranges", "bytes") // 追加
	// c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	// c.Response().WriteHeader(http.StatusOK)

	// _, err = io.Copy(multiWriter, resp.Body)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, "動画転送中にエラーが発生しました: "+err.Error())
	// }

	fmt.Println("動画保存完了:", savePath)
	return nil
}

func main() {
	// テスト用のURLを指定
	gofileDirectLink := "https://example.com/path/to/video.mp4"

	// 動画をプロキシして保存
	if err := ProxyGofileVideo(gofileDirectLink); err != nil {
		fmt.Println("エラー:", err)
	} else {
		fmt.Println("動画のプロキシと保存が成功しました")
	}
}
