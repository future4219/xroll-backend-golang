package output_port

import "gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"

type GofileRepository interface {
	Create(gofile entity.GofileVideo) error
	FindByID(id string) (entity.GofileVideo, error)
	FindByUserID(userID string) ([]entity.GofileVideo, error)
	FindByUserIDShared(userID string) ([]entity.GofileVideo, error)
	Update(gofile entity.GofileVideo) error
	Delete(id string) error
}
