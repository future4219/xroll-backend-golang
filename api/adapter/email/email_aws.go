package email

import (
	"fmt"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/config"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"

	// TODO SDKのバージョンをv2にする。互換がおそらくないため、いろいろ気を付ける。
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"

	awsDriver "gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/aws"
)

type Email struct {
	awsCli *awsDriver.Cli
}

func NewEmailDriver(awsCli *awsDriver.Cli) output_port.Email {
	return &Email{
		awsCli: awsCli,
	}
}

// Send Amazon SESを用いたメール送信
func (e Email) Send(mailAddresses []string, subject, body, htmlBody string) error {
	sess, err := e.awsCli.CreateSession()
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}

	svc := ses.New(sess)
	sender := config.EmailFrom()

	pointerMailAddresses := make([]*string, len(mailAddresses))
	for i, v := range mailAddresses {
		value := v
		pointerMailAddresses[i] = &value
	}

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: pointerMailAddresses,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(sender),
	}

	_, err = svc.SendEmail(input)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case
				ses.ErrCodeMessageRejected,
				ses.ErrCodeMailFromDomainNotVerifiedException,
				ses.ErrCodeConfigurationSetDoesNotExistException:
				return awsErr
			default:
				return fmt.Errorf("unexpected aws error: %w", awsErr)
			}
		} else {
			return err
		}
	}
	return nil
}
