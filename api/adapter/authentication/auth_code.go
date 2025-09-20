package authentication

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"
)

type AuthenticationCode struct{}

func NewAuthenticationCode() output_port.AuthCode {
	return &AuthenticationCode{}
}

func (ac *AuthenticationCode) Generate4DigitCode() string {
	num, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		panic(fmt.Sprintf("Generate4DigitCode: %v", err))
	}

	return fmt.Sprintf("%04d", num.Int64())
}
