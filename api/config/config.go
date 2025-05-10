package config

import (
	"log"
	"os"
	"strings"
	"time"
)

var (
	awsAccessKeyID            string
	awsSecretAccessKey        string
	awsRegion                 string
	env                       string
	frontendURL               string
	postCodeJPToken           string
	s3Bucket                  string
	emailFrom                 string
	sigKey                    string // JWTトークンの署名
	jst                       *time.Location
	stripeAPIKey              string
	stripeEndpointSecret      string // StripeのWebhookのエンドポイントシークレット
	videoCloudFrontURL        string
	videoCloudFrontKeyID      string
	videoCloudFrontPrivateKey string
)

func init() {
	sigKey = os.Getenv("SIG_KEY")
	if sigKey == "" {
		log.Print("SIG_KEY environment variable is empty")
	}
	env = os.Getenv("ENV")
	if env == "" {
		log.Print("ENV environment variable is empty")
	}
	awsAccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	if awsAccessKeyID == "" {
		log.Print("AWS_ACCESS_KEY_ID environment variable is empty")
	}
	awsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	if awsSecretAccessKey == "" {
		log.Print("AWS_SECRET_ACCESS_KEY environment variable is empty")
	}
	awsRegion = os.Getenv("AWS_REGION")
	if awsRegion == "" {
		log.Print("AWS_REGION environment variable is empty")
	}
	s3Bucket = os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		log.Print("S3_BUCKET environment variable is empty")
	}
	emailFrom = os.Getenv("EMAIL_FROM")
	if emailFrom == "" {
		log.Print("EMAIL_FROM environment variable is empty")
	}
	postCodeJPToken = os.Getenv("POST_CODE_JP_TOKEN")
	if postCodeJPToken == "" {
		log.Print("POST_CODE_JP_TOKEN environment variable is empty")
	}
	frontendURL = os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		log.Print("FRONTEND_URL environment variable is empty")
		frontendURL = "http://example.com"
	}
	if j, err := time.LoadLocation("Asia/Tokyo"); err != nil {
		log.Print("Failed to load location")
	} else {
		jst = j
	}
	stripeEndpointSecret = os.Getenv("STRIPE_ENDPOINT_SECRET")
	if stripeEndpointSecret == "" {
		log.Print("STRIPE_ENDPOINT_SECRET environment variable is empty")
	}
	stripeAPIKey = os.Getenv("STRIPE_API_KEY")
	if stripeAPIKey == "" {
		log.Print("STRIPE_API_KEY environment variable is empty")
	}
	videoCloudFrontURL = os.Getenv("VIDEO_CLOUD_FRONT_URL")
	if videoCloudFrontURL == "" {
		log.Print("VIDEO_CLOUD_FRONT_URL environment variable is empty")
	}
	videoCloudFrontKeyID = os.Getenv("VIDEO_CLOUD_FRONT_KEY_ID")
	if videoCloudFrontKeyID == "" {
		log.Print("VIDEO_CLOUD_FRONT_KEY_ID environment variable is empty")
	}
	videoCloudFrontPrivateKey = os.Getenv("VIDEO_CLOUD_FRONT_PRIVATE_KEY")
	if videoCloudFrontPrivateKey == "" {
		log.Print("VIDEO_CLOUD_FRONT_PRIVATE_KEY environment variable is empty")
	}
}

func IsDevelopment() bool {
	return env == "development"
}

func IsTest() bool {
	return env == "test"
}

func IsGitLabCI() bool {
	return env == "gitlab-ci"
}

func IsAWSConfigFilled() bool {
	return awsAccessKeyID != "" && awsSecretAccessKey != "" && awsRegion != ""
}

func AWSAccessKeyID() string {
	return awsAccessKeyID
}

func AWSSecretAccessKey() string {
	return awsSecretAccessKey
}

func AWSRegion() string {
	return awsRegion
}

func SigKey() string {
	return sigKey
}

func EmailFrom() string {
	return emailFrom
}

func S3Bucket() string {
	return s3Bucket
}

func PostCodeJPToken() string {
	return postCodeJPToken
}

func FrontendURL() string {
	return frontendURL
}

func StripeEndpointSecret() string {
	return stripeEndpointSecret
}

func StripeAPIKey() string {
	return stripeAPIKey
}
func JST() *time.Location {
	return jst
}

func VideoCloudFrontURL() string {
	return videoCloudFrontURL
}

func VideoCloudFrontKeyID() string {
	return videoCloudFrontKeyID
}

func VideoCloudFrontPrivateKey() string {
	return strings.Replace(videoCloudFrontPrivateKey, "\\n", "\n", -1)
}
