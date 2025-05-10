package email

import (
	"fmt"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"
)

type EmailMock struct {
}

func NewEmailDriverMock() output_port.Email {
	return &EmailMock{}
}

func (e EmailMock) Send(mailAddresses []string, subject, body, htmlBody string) error {
	fmt.Println("Email mock send")
	fmt.Println(mailAddresses)
	fmt.Println(subject)
	fmt.Println(body)
	fmt.Println(htmlBody)
	return nil
}
