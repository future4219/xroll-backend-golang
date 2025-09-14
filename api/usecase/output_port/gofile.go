package output_port

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entconst"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
)

type GofileRepository interface {
	Create(gofile entity.GofileVideo) error
	FindByID(id string) (entity.GofileVideo, error)
	FindByUserID(userID string) ([]entity.GofileVideo, error)
	FindByUserIDShared(userID string) ([]entity.GofileVideo, error)
	Update(gofile entity.GofileVideo) error
	Delete(id string) error
	HasLike(userID, videoID string) (bool, error)
	CreateLike(l entity.GofileVideoLike) error
	DeleteLike(userID, videoID string) (int64, error)
	FindLikedVideos(userID string) ([]entity.GofileVideo, error)
	Search(query GofileSearchQuery) ([]entity.GofileVideo, error)
}

type GofileSearchQuery struct {
	Q       string
	Skip    int
	Limit   int
	OrderBy entconst.GofileOrderBy
	Order   entconst.Order
}

// IsUniqueViolation 判定：一意制約違反なら true
func IsUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	// MySQL (go-sql-driver/mysql)
	var myErr *mysql.MySQLError
	if errors.As(err, &myErr) {
		// 1062 = ER_DUP_ENTRY
		return myErr.Number == 1062
	}

	return false
}
