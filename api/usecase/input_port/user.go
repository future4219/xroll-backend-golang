package input_port

import (
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
)

type UserCreate struct {
	ID    string
	StudentID string
	IdmUniv   string
	IdmBus    string
	UserType  string
}

type UserUpdate struct {
	ID    string
	StudentID string
	IdmUniv   string
	IdmBus    string
	UserType  string
}

type UserUpdatePassword struct {
	UserID      string
	NewPassword string
}

type IUserUseCase interface {
	Authenticate(token string) (string, error)
	AuthenticateForUpdateEmail(token string) (string, error)
	AuthenticateForUpdatePassword(token string) (string, error)
	Create(UserCreate) (entity.User, error)
	CreateUserWithDetail(user entity.User) error
	Delete(myself entity.User, userID string) (entity.User, error)
	FindByID(myself entity.User, userID string) (entity.User, error)
	FindByIDByAdmin(myself entity.User, userID string) (entity.User, error)
	Login(email, password string) (entity.User, string, error)
	Search(myself entity.User, query, userType string, skip int, limit int) ([]entity.User, int, error)
	SendResetPasswordMail(email string) error
	Update(myself entity.User, update UserUpdate) (entity.User, error)
}
