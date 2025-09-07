package input_port

import "gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"

type IGofileUseCase interface {
	Create(user entity.User, gofile GofileCreate) (entity.GofileVideo, error)
	FindByUserID(user entity.User) ([]entity.GofileVideo, error)
	FindByID(user entity.User, id string) (entity.GofileVideo, error)
	FindByUserIDShared(user entity.User, targetUserID string) ([]entity.GofileVideo, error)
	UpdateIsShareVideo(user entity.User, id string, isShare bool) error
	Delete(user entity.User, id string) error
}

type GofileCreate struct {
	Name     string
	GofileID string
	TagIDs   []string
	// User情報
	UserID      *string
	GofileToken *string
}
