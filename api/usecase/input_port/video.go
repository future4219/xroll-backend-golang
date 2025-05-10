package input_port

import "gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"

type VideoSearch struct {
	Limit  int
	Offset int
}

type IVideoUseCase interface {
	Search(search VideoSearch) ([]entity.Video, error)
	Create(video entity.Video) (entity.Video, error)
	CreateBulk(videos []entity.Video) error
	FindByID(id string) (entity.Video, error)
	FindByIDs(ids []string) ([]entity.Video, error)
	Like(videoID string) error
	Comment(videoID string, comment string) (error)
}
