package interactor

import (
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entconst"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/input_port"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"
)

type FileUseCase struct {
	ulid       output_port.ULID
	fileDriver output_port.FileDriver
}

func NewFileUseCase(ulid output_port.ULID, FileDriver output_port.FileDriver) input_port.IFileUseCase {
	return &FileUseCase{
		ulid:       ulid,
		fileDriver: FileDriver,
	}
}

func (u *FileUseCase) IssuePreSignedURLForPut(user entity.User) (url string, key string, err error) {
	key = u.ulid.GenerateID()
	url, err = u.fileDriver.CreatePreSignedURLForPut(key)
	return
}

func (u *FileUseCase) IssuePreSignedURLForPutVideo(user entity.User, fileName string) (url string, key string, err error) {
	key = u.ulid.GenerateID()
	url, err = u.fileDriver.CreatePreSignedURLForPut("video/original/" + key + "/" + fileName)
	return
}

func (u *FileUseCase) IssuePreSignedURLForGetVideo(user entity.User, fileName, fileID string) (url string, status entconst.FileStatus, err error) {
	url, status, err = u.fileDriver.CreateVideoPreSignedURLForGet(fileID, fileName)
	return
}
