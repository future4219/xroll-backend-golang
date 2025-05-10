package interactor

import (
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entconst"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/input_port"
)

type AuthorizationFileUseCaseDecorator struct {
	inner input_port.IFileUseCase
}

func (a AuthorizationFileUseCaseDecorator) IssuePreSignedURLForPutVideo(user entity.User, fileName string) (url string, key string, err error) {

	return a.inner.IssuePreSignedURLForPutVideo(user, fileName)
}

func (a AuthorizationFileUseCaseDecorator) IssuePreSignedURLForPut(user entity.User) (url string, key string, err error) {

	return a.inner.IssuePreSignedURLForPut(user)
}

func (a AuthorizationFileUseCaseDecorator) IssuePreSignedURLForGetVideo(user entity.User, fileName, id string) (url string, status entconst.FileStatus, err error) {

	return a.inner.IssuePreSignedURLForGetVideo(user, fileName, id)
}

func NewAuthorizationFileUseCase(inner input_port.IFileUseCase) input_port.IFileUseCase {
	return &AuthorizationFileUseCaseDecorator{inner: inner}
}
