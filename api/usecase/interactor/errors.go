package interactor

import "errors"

var ErrKind struct {
	NotFound     error
	Conflict     error
	BadRequest   error
	Unauthorized error
}

func init() {
	ErrKind.NotFound = errors.New("not found error")
	ErrKind.Conflict = errors.New("conflict error")
	ErrKind.BadRequest = errors.New("bad request error")
	ErrKind.Unauthorized = errors.New("unauthorized error")
}
