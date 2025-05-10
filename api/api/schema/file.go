package schema

import (
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
)

type FileRes struct {
	FileID      string `json:"fileId"`
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	FileSize    int32  `json:"fileSize"`
	FileURL     string `json:"fileUrl"`
}

type FileResForUser struct {
	FileID      string `json:"fileId"`
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	FileSize    int32  `json:"fileSize"`
	HasViewed   bool   `json:"hasViewed"`
	FileURL     string `json:"fileUrl"`
}

type FileResWithStatusForUser struct {
	FileID      string `json:"fileId"`
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	FileSize    int32  `json:"fileSize"`
	FileStatus  string `json:"fileStatus"`
	HasViewed   bool   `json:"hasViewed"`
	FileURL     string `json:"fileUrl"`
}

type FileResWithStatus struct {
	FileID      string `json:"fileId"`
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	FileSize    int32  `json:"fileSize"`
	FileStatus  string `json:"fileStatus"`
	FileURL     string `json:"fileUrl"`
}

type FileCreateReq struct {
	FileID      string `json:"fileId"`
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	FileSize    int32  `json:"fileSize"`
}

func FileResFromEntity(entity entity.File) FileRes {
	return FileRes{
		FileID:      entity.FileID,
		FileName:    entity.FileName,
		ContentType: entity.ContentType,
		FileSize:    entity.FileSize,
		FileURL:     entity.FileURL,
	}
}

func FileResFromEntityForUser(entity entity.FileForUser) FileResForUser {
	return FileResForUser{
		FileID:      entity.FileID,
		FileName:    entity.FileName,
		ContentType: entity.ContentType,
		FileSize:    entity.FileSize,
		HasViewed:   entity.HasViewed,
		FileURL:     entity.FileURL,
	}
}

func FileResFromEntityWithStatusForUser(entity entity.FileWithStatusForUser) FileResWithStatusForUser {
	return FileResWithStatusForUser{
		FileID:      entity.FileID,
		FileName:    entity.FileName,
		ContentType: entity.ContentType,
		FileSize:    entity.FileSize,
		FileStatus:  entity.FileStatus.String(),
		HasViewed:   entity.HasViewed,
		FileURL:     entity.FileURL,
	}
}

func FileResFromEntityWithStatus(entity entity.FileWithStatus) FileResWithStatus {
	return FileResWithStatus{
		FileID:      entity.FileID,
		FileName:    entity.FileName,
		ContentType: entity.ContentType,
		FileSize:    entity.FileSize,
		FileStatus:  entity.FileStatus.String(),
		FileURL:     entity.FileURL,
	}
}

func FileListResFromEntity(entity []entity.File) []FileRes {
	res := make([]FileRes, len(entity))
	for i, v := range entity {
		res[i] = FileResFromEntity(v)
	}
	return res
}

func FileListResFromEntityForUser(entity []entity.FileForUser) []FileResForUser {
	res := make([]FileResForUser, len(entity))
	for i, v := range entity {
		res[i] = FileResFromEntityForUser(v)
	}
	return res
}

func FileListResFromEntityWithStatusForUser(entity []entity.FileWithStatusForUser) []FileResWithStatusForUser {
	res := make([]FileResWithStatusForUser, len(entity))
	for i, v := range entity {
		res[i] = FileResFromEntityWithStatusForUser(v)
	}
	return res
}

func FileListResFromEntityWithStatus(entity []entity.FileWithStatus) []FileResWithStatus {
	res := make([]FileResWithStatus, len(entity))
	for i, v := range entity {
		res[i] = FileResFromEntityWithStatus(v)
	}
	return res
}

type IssuePreSignedURLForVideoReq struct {
	FileName string `json:"fileName"`
}

type IssuePreSignedURLForPutRes struct {
	PreSignedUrl string `json:"presignedUrl"`
	Key          string `json:"key"`
}

type IssuePreSignedURLForGetRes struct {
	PreSignedUrl string `json:"presignedUrl"`
	Status       string `json:"status"`
}
