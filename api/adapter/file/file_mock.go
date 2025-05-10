package file

import (
	"fmt"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entconst"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"
)

type Mock struct{}

func NewFileDriverMock() output_port.FileDriver {
	return &Mock{}
}

func (f Mock) CopyFile(srcID, dstID string) error {
	fmt.Println(srcID + " to " + dstID)
	return nil
}

func (f Mock) CreatePreSignedURLForGet(filepath string) (string, error) {
	fmt.Println(filepath + "presignedForGet")
	return "https://dummyimage.com/600x400/000/fff&text=test", nil
}

func (f Mock) CreateVideoPreSignedURLForGet(key, fileName string) (string, entconst.FileStatus, error) {
	fmt.Println(key + "presignedForGet")
	return "https://dummyimage.com/600x400/000/fff&text=test", entconst.FileStatusSuccess, nil
}

func (f Mock) CreatePreSignedURLForPut(filepath string) (string, error) {
	fmt.Println(filepath + "presignedForPut")
	return filepath + "presignedForPut", nil
}

func (f Mock) DeleteFileWithPath(filepath string) error {
	fmt.Println(filepath + "DeleteFileWithPath")
	return nil
}

func (f Mock) DeleteDirectoryWithPath(filepath string) error {
	fmt.Println(filepath + "DeleteFileWithPath")
	return nil
}

func (f Mock) DeleteVideoByKey(key string) error {
	fmt.Println("video/conv/" + key + "DeleteVideoWithPath")
	fmt.Println("video/original/" + key + "DeleteVideoWithPath")
	return nil
}

func (f Mock) UploadCsv(filepath string, data []byte) error {
	fmt.Println(filepath + "UploadCsv" + string(data))
	return nil
}
