package output_port

import (
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entconst"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
)

type FileDriver interface {
	CopyFile(srcID, dstID string) error
	CreatePreSignedURLForGet(filepath string) (string, error)
	CreateVideoPreSignedURLForGet(key, fileName string) (string, entconst.FileStatus, error) // fileNameは拡張子付き
	CreatePreSignedURLForPut(filepath string) (string, error)
	DeleteFileWithPath(filepath string) error
	DeleteDirectoryWithPath(filepath string) error
	DeleteVideoByKey(key string) error
	UploadCsv(url string, data []byte) error
}

type FileRepository interface {
	DeleteBulk(tx interface{}, fileIDs []string) error
	FindByID(fileID string) (entity.File, error)
}
