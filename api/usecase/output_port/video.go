package output_port

import (
	"time"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
)

type VideoRepository interface {
	Search(search VideoSearch) ([]entity.Video, error)
	Create(video entity.Video) error
	CreateBulk(videos []entity.Video) error
	FindByID(id string) (entity.Video, error)
	FindByIDs(ids []string) ([]entity.Video, error)
	Update(video entity.Video) error
	CreateComment(videoID string, comment entity.Comment) error
}

type VideoSearch struct {
	Limit  int
	Offset int
	Start  time.Time
	End    time.Time
}
