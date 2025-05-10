package entity

import "gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entconst"

type FileInterface interface {
	GetFileInfo() File
	SetFileURL(url string)
}

type FileForDeleteInterface interface {
	GetFileInfo() File
}

type FileWithStatusInterface interface {
	GetFileWithStatusInfo() FileWithStatus
	SetFileURL(url string)
	SetFileStatus(status entconst.FileStatus)
}

type File struct {
	FileID      string
	FileKind    entconst.FileKind
	FileName    string
	FileURL     string
	ContentType string
	FileSize    int32
	CreatedUser User
	UpdatedUser User
}

type FileForUser struct {
	FileID      string
	FileKind    entconst.FileKind
	FileName    string
	FileURL     string
	ContentType string
	FileSize    int32
	HasViewed   bool
	CreatedUser User
	UpdatedUser User
}

func (f File) GetID() string {
	return f.FileID
}

func (f FileForUser) GetID() string {
	return f.FileID
}

type FileWithStatusForUser struct {
	FileID      string
	FileKind    entconst.FileKind
	FileName    string
	FileURL     string
	ContentType string
	FileSize    int32
	FileStatus  entconst.FileStatus
	HasViewed   bool
	CreatedUser User
	UpdatedUser User
}

type FileWithStatus struct {
	FileID      string
	FileKind    entconst.FileKind
	FileName    string
	FileURL     string
	ContentType string
	FileSize    int32
	FileStatus  entconst.FileStatus
	CreatedUser User
	UpdatedUser User
}

func (f File) GetFileInfo() File {
	return f
}
func (f FileWithStatus) GetFileInfo() File {
	return File{
		FileID:      f.FileID,
		FileKind:    f.FileKind,
		FileName:    f.FileName,
		FileURL:     f.FileURL,
		ContentType: f.ContentType,
		FileSize:    f.FileSize,
		CreatedUser: f.CreatedUser,
		UpdatedUser: f.UpdatedUser,
	}
}

func (f FileWithStatus) GetFileWithStatusInfo() FileWithStatus {
	return f
}

func (f FileWithStatusForUser) GetFileWithStatusInfo() FileWithStatus {
	return FileWithStatus{
		FileID:      f.FileID,
		FileKind:    f.FileKind,
		FileName:    f.FileName,
		FileURL:     f.FileURL,
		ContentType: f.ContentType,
		FileSize:    f.FileSize,
		FileStatus:  f.FileStatus,
		CreatedUser: f.CreatedUser,
		UpdatedUser: f.UpdatedUser,
	}
}

func (f *File) SetFileURL(url string) {
	f.FileURL = url
}
func (f *FileWithStatus) SetFileURL(url string) {
	f.FileURL = url
}
func (f *FileWithStatusForUser) SetFileURL(url string) {
	f.FileURL = url
}

func (f *FileWithStatus) SetFileStatus(status entconst.FileStatus) {
	f.FileStatus = status
}
func (f *FileWithStatusForUser) SetFileStatus(status entconst.FileStatus) {
	f.FileStatus = status
}
