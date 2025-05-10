package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/config"
)

type Cli struct {
	awsRegion string
}

func NewCli() *Cli {
	return &Cli{
		awsRegion: config.AWSRegion(),
	}
}

func (cli *Cli) CreateSession() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cli.awsRegion),
	})
	if err != nil {
		return nil, err
	}

	return sess, nil
}
