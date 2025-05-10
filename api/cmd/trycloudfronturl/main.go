package main

import (
	"fmt"

	"github.com/go-resty/resty/v2"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/aws"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/cache"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/file"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/config"
)

func main() {
	var awsCli *aws.Cli
	if config.IsAWSConfigFilled() {
		awsCli = aws.NewCli()
	}

	cache := cache.GetInstance()
	fileDriver := file.NewFileDriver(awsCli, cache)

	m3u8FileName := "mov_hts-samp001.m3u8"

	// video/conv/testulid1/mov_hts-samp001conv_presigned.m3u8というファイルを既にs3にアップロード済
	m3u8FileURL, _, err := fileDriver.CreateVideoPreSignedURLForGet("testulid1", m3u8FileName)
	if err != nil {
		fmt.Printf("Error while create presigned url for get: %s", err)
		return
	}

	// Resty クライアントを作成
	client := resty.New()

	// Presigned URL からファイルを取得
	resp, err := client.R().Get(m3u8FileURL)
	if err != nil {
		fmt.Printf("Error while downloading file: %s", err)
		return
	}

	// ステータスコードが 200 OK であることを確認
	if resp.StatusCode() != 200 {
		fmt.Printf("Error while downloading file, status code: %d", resp.StatusCode())
		return
	}

	// ダウンロードしたデータ（m3u8ファイルの内容）を取得
	data := resp.Body()

	// ここで 'data' を使って必要な処理を行う
	// 例えば、m3u8ファイルの内容を文字列として表示
	fmt.Println("Downloaded m3u8 file content:")
	fmt.Println(string(data))

	fmt.Println(fileDriver.CreateVideoPreSignedURLForGet("ulid", "mov_hts-samp001.m3u8"))
}
