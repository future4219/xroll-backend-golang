package file

import (
	"bufio"
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudfront/sign"

	"github.com/go-resty/resty/v2"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/config"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entconst"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	awsDriver "gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/aws"
)

type File struct {
	awsCli      *awsDriver.Cli
	cacheDriver output_port.Cache
}

func NewFileDriver(awsCli *awsDriver.Cli, cacheDriver output_port.Cache) output_port.FileDriver {
	return &File{
		awsCli:      awsCli,
		cacheDriver: cacheDriver,
	}
}

func (f File) CopyFile(srcID, dstID string) error {
	sess, err := f.awsCli.CreateSession()
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	svc := s3.New(sess)
	bucket := aws.String(config.S3Bucket())
	copySource := aws.String(config.S3Bucket() + "/" + srcID)
	key := aws.String(dstID)

	if _, err := svc.CopyObject(&s3.CopyObjectInput{
		CopySource: copySource,
		Bucket:     bucket,
		Key:        key,
	}); err != nil {
		return fmt.Errorf("copy object: %w", err)
	}

	return nil
}

func (f File) CreatePreSignedURLForGet(filepath string) (string, error) {
	sess, err := f.awsCli.CreateSession()
	if err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}

	svc := s3.New(sess)
	key := aws.String(encodeFilepath(filepath))
	bucket := aws.String(config.S3Bucket())

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: bucket,
		Key:    key,
	})
	url, err := req.Presign(24 * time.Hour)
	if err != nil {
		return "", fmt.Errorf("create pre signed url for get: %w", err)
	}

	return url, nil
}

func (f File) createVideoPreSignedURLForGet(filepath string) (string, error) {
	sess, err := f.awsCli.CreateSession()
	if err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}

	svc := s3.New(sess)
	key := aws.String(encodeFilepath(filepath))
	bucket := aws.String(config.S3Bucket())

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: bucket,
		Key:    key,
	})
	url, err := req.Presign(48 * time.Hour)
	if err != nil {
		return "", fmt.Errorf("create pre signed url for get: %w", err)
	}

	return url, nil
}

// CreateVideoPreSignedURLForGet
// 1. m3u8ファイルを取得
// 2. m3u8ファイルの内容を読み込み、TSファイルの行を見つけたら、presigned URLに変換
// 3. 新しいm3u8ファイルの内容を構築
// 4. 新しいm3u8ファイルをPUT
// 5. 新しいm3u8ファイルのpresigned URLを返す
func (f File) CreateVideoPreSignedURLForGet(key, fileName string) (url string, status entconst.FileStatus, err error) {
	value, found := f.cacheDriver.Get(key + "/" + fileName)
	if found {
		url, isString := value.(string)
		if !isString {
			return "", entconst.FileStatusFailed, fmt.Errorf("cache value is not string: %v", value)
		}
		return url, entconst.FileStatusSuccess, nil
	}
	fileNameParts := strings.Split(fileName, ".")
	fileNameWithoutLast := fileNameParts[:len(fileNameParts)-1]
	fileNameBeforePeriod := strings.Join(fileNameWithoutLast, ".")

	res, err := f.CreatePreSignedURLForGet("video/conv/" + key + "/" + fileNameBeforePeriod + "conv.m3u8")
	if err != nil {
		return "", entconst.FileStatusFailed, fmt.Errorf("create pre signed url for get: %w", err)
	}

	// Resty クライアントを作成
	client := resty.New()

	// Presigned URL からファイルを取得
	resp, err := client.R().Get(res)
	if err != nil {
		return "", entconst.FileStatusFailed, fmt.Errorf("get file: %w", err)
	}
	switch resp.StatusCode() {
	case 200:
	// pass
	case 404:
		// TODO: コンバートの進行具合が必要な場合は、lambdaの変換の進行度を何らかの方法で返す処理を作る
		// コンバート中か確認する
		// originalファイルがあれば変換中だと暫定する
		res, err := f.CreatePreSignedURLForGet("video/original/" + key + "/" + fileName)
		if err != nil {
			return "", entconst.FileStatusFailed, fmt.Errorf("create pre signed url for get: %w", err)
		}
		resp, err := client.R().Head(res)
		if err != nil {
			return "", entconst.FileStatusFailed, fmt.Errorf("get file: %w", err)
		}
		if resp.StatusCode() != 200 {
			return "", entconst.FileStatusInProgress, nil
		}
	default:
		return "", entconst.FileStatusFailed, fmt.Errorf("status code is not 200: %d", resp.StatusCode())
	}

	data := resp.Body()
	m3u8Content := string(data)

	// 新しいm3u8ファイルの内容を構築
	var newM3u8Content bytes.Buffer
	scanner := bufio.NewScanner(strings.NewReader(m3u8Content))
	for scanner.Scan() {
		line := scanner.Text()

		// TSファイルの行を見つけたら、presigned URLに変換
		if strings.HasSuffix(line, ".ts") {
			res, err := f.createCloudFrontPresignedURLForGet("video/conv/" + key + "/" + line)
			if err != nil {
				return "", entconst.FileStatusFailed, fmt.Errorf("create pre signed url for get: %w", err)
			}
			newM3u8Content.WriteString(res + "\n")
		} else {
			newM3u8Content.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return "", entconst.FileStatusFailed, fmt.Errorf("scan m3u8 file: %w", err)
	}

	m3u8Data := newM3u8Content.Bytes()

	updateURL, err := f.CreatePreSignedURLForPut("video/conv/" + key + "/" + fileNameBeforePeriod + "conv_presigned.m3u8")
	if err != nil {
		return "", entconst.FileStatusFailed, fmt.Errorf("create pre signed url for put: %w", err)
	}

	_, err = client.R().
		SetHeader("Content-Type", "application/vnd.apple.mpegurl").
		SetBody(m3u8Data).
		Put(updateURL)
	if err != nil {
		return "", entconst.FileStatusFailed, fmt.Errorf("put file: %w", err)
	}

	getURL, err := f.createVideoPreSignedURLForGet("video/conv/" + key + "/" + fileNameBeforePeriod + "conv_presigned.m3u8")
	if err != nil {
		return "", entconst.FileStatusFailed, fmt.Errorf("create pre signed url for get: %w", err)
	}
	f.cacheDriver.Set(key+"/"+fileName, getURL)
	return getURL, entconst.FileStatusSuccess, nil
}

// LoadPrivateKeyFromString はPEM形式の秘密鍵文字列から*rsa.PrivateKeyを生成します
func LoadPrivateKeyFromString(privateKeyStr string) (*rsa.PrivateKey, error) {
	// PEM形式のデータをデコード
	block, _ := pem.Decode([]byte(privateKeyStr))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the key")
	}
	// 秘密鍵をパース
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA private key: %v", err)
	}

	return privKey, nil
}

func (f File) createCloudFrontPresignedURLForGet(filepath string) (string, error) {
	// 秘密鍵を読み込み
	privateKeyStr := config.VideoCloudFrontPrivateKey()

	privateKey, err := LoadPrivateKeyFromString(privateKeyStr)
	if err != nil {
		return "", fmt.Errorf("failed to load private key: %v", err)
	}
	// サイナーを作成 ※KeyID
	signer := sign.NewURLSigner(config.VideoCloudFrontKeyID(), privateKey)

	resourceURL := fmt.Sprintf("%s/%s", config.VideoCloudFrontURL(), encodeFilepath(filepath))

	// プリサインURLを生成
	return signer.Sign(resourceURL, time.Now().AddDate(0, 0, 2)) // キャッシュが１日であることを踏まえて2日の有効期限
}

func (f File) CreatePreSignedURLForPut(filepath string) (string, error) {
	sess, err := f.awsCli.CreateSession()
	if err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}

	svc := s3.New(sess)
	key := aws.String(encodeFilepath(filepath))
	bucket := aws.String(config.S3Bucket())

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: bucket,
		Key:    key,
	})
	url, err := req.Presign(5 * time.Minute)
	if err != nil {
		return "", fmt.Errorf("create pre signed url for put: %w", err)
	}

	return url, nil
}

func (f File) DeleteFileWithPath(filepath string) error {
	sess, err := f.awsCli.CreateSession()
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}

	svc := s3.New(sess)
	key := aws.String(filepath)
	bucket := aws.String(config.S3Bucket())

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		return err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: bucket,
		Key:    key,
	})

	if err != nil {
		return err
	}

	return nil
}

func (f File) DeleteDirectoryWithPath(filepath string) error {
	sess, err := f.awsCli.CreateSession()
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}

	svc := s3.New(sess)
	bucket := aws.String(config.S3Bucket())
	prefix := aws.String(filepath)

	err = svc.ListObjectsV2Pages(&s3.ListObjectsV2Input{
		Bucket: bucket,
		Prefix: prefix,
	}, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		objects := make([]*s3.ObjectIdentifier, 0, len(page.Contents))
		for _, obj := range page.Contents {
			objects = append(objects, &s3.ObjectIdentifier{
				Key: obj.Key,
			})
		}
		if _, err = svc.DeleteObjects(&s3.DeleteObjectsInput{
			Bucket: bucket,
			Delete: &s3.Delete{
				Objects: objects,
			},
		}); err != nil {
			return false
		}

		return true
	})

	if err != nil {
		return err
	}
	input := &s3.ListObjectsV2Input{
		Bucket:  bucket,
		Prefix:  prefix,
		MaxKeys: aws.Int64(1),
	}

	result, err := svc.ListObjectsV2(input)
	if err != nil {
		return err
	}
	if len(result.Contents) > 0 {
		return fmt.Errorf("directory is not empty: %s", filepath)
	}

	return nil
}

func (f File) DeleteVideoByKey(key string) error {
	prefixPaths := []string{"video/conv/", "video/original/"}
	for _, prefixPath := range prefixPaths {
		if err := f.DeleteDirectoryWithPath(prefixPath + key); err != nil {
			return err
		}
	}
	return nil
}

func (f File) UploadCsv(filepath string, data []byte) error {
	sess, err := f.awsCli.CreateSession()
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}

	svc := s3.New(sess)
	key := aws.String(encodeFilepath(filepath))
	bucket := aws.String(config.S3Bucket())

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket:             bucket,
		Key:                key,
		Body:               bytes.NewReader(data),
		ContentType:        aws.String("text/csv"),
		ContentDisposition: aws.String("attachment; filename=memberList.csv"),
	})

	err = req.Send()
	if err != nil {
		return fmt.Errorf("upload csv: %w", err)
	}

	return nil
}

func encodeFilepath(filepath string) string {
	segments := strings.Split(filepath, "/")
	for i, segment := range segments {
		segments[i] = url.PathEscape(segment)
	}
	return strings.Join(segments, "/")
}
