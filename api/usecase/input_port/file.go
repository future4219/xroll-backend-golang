package input_port

import (
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entconst"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
)

type FileCreate struct {
	FileID      string
	FileName    string
	ContentType string
	FileSize    int32
}

type FileUpdate struct {
	FileID        string
	FileName      string
	ContentType   string
	FileSize      int32
	CreatedUserID string
	UpdatedUserID string
}

type IFileUseCase interface {
	IssuePreSignedURLForPut(user entity.User) (url string, key string, err error)
	IssuePreSignedURLForPutVideo(user entity.User, fileName string) (url string, key string, err error)
	IssuePreSignedURLForGetVideo(user entity.User, fileName, id string) (url string, status entconst.FileStatus, err error)
}
