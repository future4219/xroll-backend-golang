package output_port

import (
	"context"
	"io"
)

type GofileAPIDriver interface {
	GetContent(gofileID string, gofileToken string) (GofileGetContentRes, error)
	IssueDirectLink(contentID, gofileToken string) (GofileDirectLink, error)
	Upload(ctx context.Context, filename, folderID string, r io.Reader) (GofileUploadData, error)
}

type GofileGetContentRes struct {
	Status   string        `json:"status"`
	Data     GofileContent `json:"data"`
	Metadata interface{}   `json:"metadata"`
}

type GofileContent struct {
	CanAccess      bool                        `json:"canAccess"`
	IsOwner        bool                        `json:"isOwner"`
	ID             string                      `json:"id"`
	ParentFolder   string                      `json:"parentFolder"`
	Type           string                      `json:"type"`
	Name           string                      `json:"name"`
	CreateTime     int64                       `json:"createTime"`
	ModTime        int64                       `json:"modTime"`
	LastAccess     int64                       `json:"lastAccess"`
	Size           int64                       `json:"size"`
	DownloadCount  int64                       `json:"downloadCount"`
	MD5            string                      `json:"md5"`
	Mimetype       string                      `json:"mimetype"`
	Servers        []string                    `json:"servers"`
	ServerSelected string                      `json:"serverSelected"`
	DirectLinks    map[string]GofileDirectLink `json:"directLinks"`
	Link           string                      `json:"link"`
	Thumbnail      string                      `json:"thumbnail"`
}

type GofileDirectLink struct {
	ExpireTime       int64    `json:"expireTime"`
	SourceIpsAllowed []string `json:"sourceIpsAllowed"`
	DomainsAllowed   []string `json:"domainsAllowed"`
	Auth             []string `json:"auth"`
	IsReqLink        bool     `json:"isReqLink"`
	DirectLink       string   `json:"directLink"`
}

type UploadResult struct {
	Status string           `json:"status"`
	Data   GofileUploadData `json:"data"`
}

// Data 部分（2025/10 現行スキーマ相当）
type GofileUploadData struct {
	ID               string   `json:"id"` // contentId
	ParentFolder     string   `json:"parentFolder"`
	ParentFolderCode string   `json:"parentFolderCode"`
	DownloadPage     string   `json:"downloadPage"`
	GuestToken       string   `json:"guestToken"`
	MD5              string   `json:"md5"`
	MimeType         string   `json:"mimetype"`
	Name             string   `json:"name"`
	Servers          []string `json:"servers"`
	Size             int64    `json:"size"`
	Type             string   `json:"type"`
	CreateTime       int64    `json:"createTime"`
	ModTime          int64    `json:"modTime"`
	// 互換
	FileID     string `json:"fileId,omitempty"`
	DirectLink string `json:"directLink,omitempty"`
	Code       string `json:"code,omitempty"`
}
