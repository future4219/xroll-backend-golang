package interactor

import "errors"

var ErrKind struct {
	NotFound   error
	Conflict   error
	BadRequest error
}

func init() {
	ErrKind.NotFound = errors.New("not found error")
	ErrKind.Conflict = errors.New("conflict error")
	ErrKind.BadRequest = errors.New("bad request error")
}
