package entconst

type FileStatus string

const (
	FileStatusSuccess    FileStatus = "success"
	FileStatusInProgress FileStatus = "inProgress"
	FileStatusFailed     FileStatus = "failed"
)

func (e FileStatus) String() string {
	return string(e)
}

type FileKind string

const (
	FileKindEditor    FileKind = "editor"
	FileKindVideo     FileKind = "video"
	FileKindDocument  FileKind = "document"
	FileKindThumbnail FileKind = "thumbnail"
)

func (e FileKind) String() string {
	return string(e)
}
