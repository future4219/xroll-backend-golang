package interactor

import (
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/input_port"
)

type AuthorizationGofileUseCaseDecorator struct {
	inner input_port.IGofileUseCase
}

func NewAuthorizationGofileUseCase(inner input_port.IGofileUseCase) input_port.IGofileUseCase {
	return &AuthorizationGofileUseCaseDecorator{inner: inner}
}

func (a AuthorizationGofileUseCaseDecorator) Create(user entity.User, gofile input_port.GofileCreate) (entity.GofileVideo, error) {
	return a.inner.Create(user, gofile)
}

func (a AuthorizationGofileUseCaseDecorator) FindByUserID(user entity.User) ([]entity.GofileVideo, error) {
	return a.inner.FindByUserID(user)
}

func (a AuthorizationGofileUseCaseDecorator) FindByID(user entity.User, id string) (entity.GofileVideo, error) {
	return a.inner.FindByID(user, id)
}

func (a AuthorizationGofileUseCaseDecorator) FindByUserIDShared(user entity.User, targetUserID string) ([]entity.GofileVideo, error) {
	return a.inner.FindByUserIDShared(user, targetUserID)
}

func (a AuthorizationGofileUseCaseDecorator) UpdateIsShareVideo(user entity.User, id string, isShare bool) error {
	return a.inner.UpdateIsShareVideo(user, id, isShare)
}
