package input_port

import "gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"

type ThreadSearch struct {
	Limit  int
	Offset int
}

type IThreadUseCase interface {
	Search(search ThreadSearch) ([]entity.Thread, error)
	Create(video entity.Thread) (entity.Thread, error)
	FindByID(id string) (entity.Thread, error)
	FindByIDs(ids []string) ([]entity.Thread, error)
	Like(videoID string) error
	Comment(videoID string, comment string) (error)
}
