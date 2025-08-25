package input_port

import "gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"

type IGofileUseCase interface {
	Create(gofile GofileCreate) (entity.GofileVideo, error)
	FindByUserID(userID string) ([]entity.GofileVideo, error)
	FindByID(id string) (entity.GofileVideo, error)
	FindByUserIDShared(userID string) ([]entity.GofileVideo, error)
	UpdateIsShareVideo(id string, isShare bool) error
}

type GofileCreate struct {
	Name     string
	GofileID string
	TagIDs   []string
	// User情報
	UserID      *string
	GofileToken string
}
