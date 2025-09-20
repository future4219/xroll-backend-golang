package constructor

import (
	"time"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/validation"
)

type NewRegisterVerificationCreateArgs struct {
	RegisterVerificationID   string
	Email                    string
	Password                 string
	HashedPassword           string
	HashedAuthenticationCode string
	ExpiresAt                time.Time
}

func NewRegisterVerificationCreate(arg NewRegisterVerificationCreateArgs) (entity.RegisterVerification, error) {
	err := validation.ValidateEmail(arg.Email, false)
	if err != nil {
		return entity.RegisterVerification{}, err
	}
	if err = validation.ValidatePassword(arg.Password); err != nil {
		return entity.RegisterVerification{}, err
	}

	return entity.RegisterVerification{
		RegisterVerificationID:   arg.RegisterVerificationID,
		Email:                    arg.Email,
		HashedPassword:           arg.HashedPassword,
		HashedAuthenticationCode: arg.HashedAuthenticationCode,
		ExpiresAt:                arg.ExpiresAt,
	}, nil
}
