package interactor

import (
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/input_port"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"
)

type TwitterUseCase struct {
	Twitter output_port.Twitter
}

func NewTwitterUseCase(twitter output_port.Twitter) input_port.ITwitterUseCase {
	return &TwitterUseCase{
		Twitter: twitter,
	}
}

func (u *TwitterUseCase) GetVideoByURL(url string) (string, error) {
	tweet, err := u.Twitter.GetVideoByURL(url)
	if err != nil {
		return "", err
	}
	return tweet, nil
}
