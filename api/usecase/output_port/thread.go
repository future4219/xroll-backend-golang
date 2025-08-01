package output_port

import (
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
)

type ThreadRepository interface {
	Search(search ThreadSearch) ([]entity.Thread, error)
	Create(thread entity.Thread) error
	FindByID(id string) (entity.Thread, error)
	FindByIDs(ids []string) ([]entity.Thread, error)
	Update(thread entity.Thread) error
	CreateComment(threadID string, comment entity.ThreadComment) error
}

type ThreadSearchOrderBy string

const (
	ThreadSearchOrderByRanking   ThreadSearchOrderBy = "ranking"
	ThreadSearchOrderByCreatedAt ThreadSearchOrderBy = "created_at"
)

func (threadSearchOrderBy ThreadSearchOrderBy) ToString() string {
	switch threadSearchOrderBy {
	case ThreadSearchOrderByRanking:
		return "ranking"
	case ThreadSearchOrderByCreatedAt:
		return "created_at"
	default:
		return "created_at"
	}
}

type ThreadSearchOrder string

const (
	ThreadSearchOrderAsc  ThreadSearchOrder = "asc"
	ThreadSearchOrderDesc ThreadSearchOrder = "desc"
)

func (threadSearchOrder ThreadSearchOrder) ToString() string {
	switch threadSearchOrder {
	case ThreadSearchOrderAsc:
		return "asc"
	case ThreadSearchOrderDesc:
		return "desc"
	default:
		return "asc"
	}
}

type ThreadSearch struct {
	Limit   int
	Offset  int
	OrderBy ThreadSearchOrderBy
	Order   ThreadSearchOrder
}
