package interactor

import (
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/input_port"
)

type AuthorizationVideoUseCaseDecorator struct {
	inner input_port.IVideoUseCase
}

func NewAuthorizationVideoUseCase(inner input_port.IVideoUseCase) input_port.IVideoUseCase {
	return &AuthorizationVideoUseCaseDecorator{inner: inner}
}

func (a AuthorizationVideoUseCaseDecorator) Search(search input_port.VideoSearch) ([]entity.Video, error) {
	return a.inner.Search(search)
}

func (a AuthorizationVideoUseCaseDecorator) Create(video entity.Video) (entity.Video, error) {
	return a.inner.Create(video)
}

func (a AuthorizationVideoUseCaseDecorator) CreateBulk(videos []entity.Video) error {
	return a.inner.CreateBulk(videos)
}

func (a AuthorizationVideoUseCaseDecorator) FindByID(id string) (entity.Video, error) {
	return a.inner.FindByID(id)
}

func (a AuthorizationVideoUseCaseDecorator) FindByIDs(ids []string) ([]entity.Video, error) {
	return a.inner.FindByIDs(ids)
}

func (a AuthorizationVideoUseCaseDecorator) Like(videoID string) error {
	return a.inner.Like(videoID)
}

func (a AuthorizationVideoUseCaseDecorator) Comment(videoID string, comment string) error {
	return a.inner.Comment(videoID, comment)
}
