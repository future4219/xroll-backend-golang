package entconst

import "errors"

var (
	ErrInvalidGofileOrderBy = errors.New("invalid gofile order by")
)

type GofileOrderBy string

const (
	GofileOrderByCreatedAt GofileOrderBy = "created_at"
	GofileOrderByUpdatedAt GofileOrderBy = "updated_at"
	GofileOrderByLikeCount GofileOrderBy = "like_count"
	GofileOrderByPlayCount GofileOrderBy = "play_count"
)

func (j GofileOrderBy) String() string {
	return string(j)
}

func validateGofileOrderBy(orderBy string) error {
	switch GofileOrderBy(orderBy) {
	case GofileOrderByCreatedAt, GofileOrderByLikeCount, GofileOrderByPlayCount, GofileOrderByUpdatedAt:
		return nil
	default:
		return ErrInvalidGofileOrderBy
	}
}

func NewGofileOrderBy(orderBy string) (GofileOrderBy, error) {
	if orderBy == "" {
		return GofileOrderByUpdatedAt, nil
	}
	if err := validateGofileOrderBy(orderBy); err != nil {
		return "", err
	}
	return GofileOrderBy(orderBy), nil
}
