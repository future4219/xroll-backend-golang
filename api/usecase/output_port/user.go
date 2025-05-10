package output_port

import (
	"errors"
	"time"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
)

var (
	TokenScopeGeneral                 = "general"
	TokenGeneralExpireDuration        = 24 * time.Hour // 1 day
	TokenScopeUpdateEmail             = "updateEmail"
	TokenEmailUpdateExpireDuration    = 24 * time.Hour // 1 day
	TokenScopeUpdatePassword          = "updatePassword"
	TokenChangePasswordExpireDuration = 24 * time.Hour // 1 day
	ErrUnknownScope                   = errors.New("unknown scope")
	ErrTokenExpired                   = errors.New("token expired")
	ErrTokenIssuedFutureTime          = errors.New("token issued future time")
	ErrTokenScopeInvalid              = errors.New("token scope invalid")
)

type UserRepository interface {
	Create(user entity.User) error
	CreateWithTx(tx interface{}, user entity.User) error
	Delete(userID string) error
	FindByID(userID string) (entity.User, error)
	FindByEmail(email string) (entity.User, error)
	ListByEmails(emails []string) ([]entity.User, error)
	FindByLoginID(loginID string) (entity.User, error)
	ListByLoginIDs(loginIDs []string) ([]entity.User, error)
	Search(Query string, userType string, Skip int, Limit int) ([]entity.User, int, error)
	Update(entity.User) error
	FindMaxLoginID() (string, error)
	GetAdminUser() ([]entity.User, error)
}

type UserAuth interface {
	Authenticate(token string) (string, error)
	AuthenticateForUpdateEmail(token string) (string, error)
	AuthenticateForUpdatePassword(token string) (string, error)
	HashPassword(password string) (string, error)
	IssueUserToken(user entity.User, issuedAt time.Time) (string, error)
	IssueUserTokenForUpdateEmail(user entity.User, issuedAt time.Time) (string, error)
	IssueUserTokenForUpdatePassword(user entity.User, issuedAt time.Time) (string, error)
	GenerateInitialPassword(length int) (string, error)
}
