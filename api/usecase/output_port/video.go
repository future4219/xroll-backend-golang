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

type VideoSearchOrderBy string

const (
	VideoSearchOrderByRanking   VideoSearchOrderBy = "ranking"
	VideoSearchOrderByCreatedAt VideoSearchOrderBy = "created_at"
)

func (videoSearchOrderBy VideoSearchOrderBy) ToString() string {
	switch videoSearchOrderBy {
	case VideoSearchOrderByRanking:
		return "ranking"
	case VideoSearchOrderByCreatedAt:
		return "created_at"
	default:
		return "created_at"
	}
}

type VideoSearchOrder string

const (
	VideoSearchOrderAsc  VideoSearchOrder = "asc"
	VideoSearchOrderDesc VideoSearchOrder = "desc"
)

func (videoSearchOrder VideoSearchOrder) ToString() string {
	switch videoSearchOrder {
	case VideoSearchOrderAsc:
		return "asc"
	case VideoSearchOrderDesc:
		return "desc"
	default:
		return "asc"
	}
}

type VideoSearch struct {
	Limit   int
	Offset  int
	Start   time.Time
	End     time.Time
	OrderBy VideoSearchOrderBy
	Order   VideoSearchOrder
}
