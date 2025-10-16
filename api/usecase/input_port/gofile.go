package input_port

import (
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entconst"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
)

type IGofileUseCase interface {
	Create(user entity.User, gofile GofileCreate) (entity.GofileVideo, error)
	Update(user entity.User, update GofileUpdate) (entity.GofileVideo, error)
	FindByUserID(user entity.User) ([]entity.GofileVideo, error)
	FindByID(user entity.User, id string) (entity.GofileVideo, bool, error)
	FindByUserIDShared(user entity.User, targetUserID string) ([]entity.GofileVideo, error)
	UpdateIsShareVideo(user entity.User, id string, isShare bool) error
	Delete(user entity.User, id string) error
	LikeVideo(user entity.User, videoID string) error
	UnlikeVideo(user entity.User, videoID string) error
	FindLikedVideos(user entity.User) ([]entity.GofileVideo, error)
	Search(user entity.User, query GofileSearchQuery) ([]entity.GofileVideo, error)
	CreateComment(user entity.User, input GofileVideoCommentCreate) (entity.GofileVideoComment, error)
	CreateFromTwimgURL(user entity.User, srcUrl string) (entity.GofileVideo, error)
}

type GofileCreate struct {
	Name     string
	GofileID string
	TagIDs   []string
	// User情報
	UserID      *string
	GofileToken *string
}

type GofileUpdate struct {
	ID          string
	Name        string
	Description string
	TagIDs      []string
	IsShare     bool
}

type GofileSearchQuery struct {
	Q       string
	Skip    int
	Limit   int
	OrderBy entconst.GofileOrderBy
	Order   entconst.Order
}

type GofileVideoCommentCreate struct {
	VideoID string
	Comment string
}
