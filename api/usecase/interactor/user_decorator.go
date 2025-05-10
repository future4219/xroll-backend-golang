package interactor

import (
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/input_port"
)

type AuthorizationUserUseCaseDecorator struct {
	inner input_port.IUserUseCase
}

func NewAuthorizationUserUseCase(inner input_port.IUserUseCase) input_port.IUserUseCase {
	return &AuthorizationUserUseCaseDecorator{inner: inner}
}

func (a AuthorizationUserUseCaseDecorator) Authenticate(token string) (string, error) {
	return a.inner.Authenticate(token)
}

func (a AuthorizationUserUseCaseDecorator) AuthenticateForUpdateEmail(token string) (string, error) {
	return a.inner.AuthenticateForUpdateEmail(token)
}

func (a AuthorizationUserUseCaseDecorator) AuthenticateForUpdatePassword(token string) (string, error) {
	return a.inner.AuthenticateForUpdatePassword(token)
}

func (a AuthorizationUserUseCaseDecorator) Create(create input_port.UserCreate) (entity.User, error) {
	return a.inner.Create(create)
}

func (a AuthorizationUserUseCaseDecorator) CreateUserWithDetail(user entity.User) error {
	return a.inner.CreateUserWithDetail(user)
}

func (a AuthorizationUserUseCaseDecorator) Delete(myself entity.User, userID string) (entity.User, error) {

	return a.inner.Delete(myself, userID)
}

func (a AuthorizationUserUseCaseDecorator) FindByID(myself entity.User, userID string) (entity.User, error) {

	return a.inner.FindByID(myself, userID)
}

func (a AuthorizationUserUseCaseDecorator) FindByIDByAdmin(myself entity.User, userID string) (entity.User, error) {

	return a.inner.FindByIDByAdmin(myself, userID)
}

func (a AuthorizationUserUseCaseDecorator) Login(email, password string) (entity.User, string, error) {
	return a.inner.Login(email, password)
}

func (a AuthorizationUserUseCaseDecorator) Search(myself entity.User, query, userType string, skip int, limit int) ([]entity.User, int, error) {

	return a.inner.Search(myself, query, userType, skip, limit)
}

func (a AuthorizationUserUseCaseDecorator) SendResetPasswordMail(email string) error {
	return a.inner.SendResetPasswordMail(email)
}

func (a AuthorizationUserUseCaseDecorator) Update(myself entity.User, update input_port.UserUpdate) (entity.User, error) {

	return a.inner.Update(myself, update)
}
