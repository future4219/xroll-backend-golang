package model

import (
	"time"
)

type User struct {
	ID             string    `gorm:"primaryKey;type:varchar(255)"`
	Name           string    `gorm:"type:varchar(100)"`
	Age            int       `gorm:"default:0"`
	UserType       string    `gorm:"type:varchar(20);index"`        // "guest" | "member" | "admin"
	Email          *string   `gorm:"type:varchar(255);uniqueIndex"` // ゲストはNULL
	HashedPassword *string   `gorm:"type:varchar(255)"`             // ゲストはNULL
	GofileToken    *string    `gorm:"type:varchar(255);unique"`
	EmailVerified  bool      `gorm:"default:false"`
	IsDeleted      bool      `gorm:"default:false"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

type RefreshToken struct {
	ID         string     `gorm:"primaryKey;type:varchar(255)"`
	UserID     string     `gorm:"index;type:varchar(255)"`
	TokenHash  string     `gorm:"size:64;uniqueIndex"`
	SessionID  string     `gorm:"size:64;index"`
	UserAgent  string     `gorm:"size:255"`
	IPAddress  string     `gorm:"size:45"`
	CreatedAt  time.Time  `gorm:"autoCreateTime"`
	ExpiresAt  time.Time  `gorm:"index"`
	RevokedAt  *time.Time `gorm:"index"`
	ReplacedBy *int64
}

func (RefreshToken) TableName() string { return "refresh_tokens" }

type Video struct {
	ID            string         `gorm:"primaryKey;type:varchar(255)"`
	Ranking       int            `gorm:"not null"`
	VideoURL      string         `gorm:"type:varchar(768);unique"` // ← 追加要件対応済み
	ThumbnailURL  string         `gorm:"type:text;not null"`
	TweetURL      *string        `gorm:"type:text"` // NULL許容
	DownloadCount int            `gorm:"default:0"`
	LikeCount     int            `gorm:"default:0"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	Comments      []VideoComment `gorm:"foreignKey:VideoID"`
}

type VideoComment struct {
	ID        string    `gorm:"primaryKey;type:varchar(255)"`
	VideoID   string    `gorm:"type:varchar(255);not null"` // 外部キー制約つけたければ↓追加
	Comment   string    `gorm:"type:text;not null"`
	LikeCount int       `gorm:"default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type GofileVideo struct {
	ID                  string               `gorm:"primaryKey;type:varchar(255)"`
	Name                string               `gorm:"type:varchar(255);not null"`
	GofileID            string               `gorm:"type:varchar(255);not null"` // GofileのID
	GofileDirectURL     string               `gorm:"type:text;not null"`         // Gofileの直接ダウンロードURL
	VideoURL            string               `gorm:"type:text;not null"`
	ThumbnailURL        string               `gorm:"type:text;not null"`                                                                                              // 動画のサムネイルURL
	LikeCount           int                  `gorm:"default:0"`                                                                                                       // いいね数
	IsShared            bool                 `gorm:"default:false"`                                                                                                   // 動画が共有されているかどうか
	GofileTags          []GofileTag          `gorm:"many2many:gofile_video_tags;foreignKey:ID;joinForeignKey:GofileVideoID;References:ID;joinReferences:GofileTagID"` // GofileTagとの多対多リレーション
	GofileVideoComments []GofileVideoComment `gorm:"foreignKey:GofileVideoID"`                                                                                        // GofileVideoCommentとのリレーション
	UserID              string               `gorm:"type:varchar(255);not null"`                                                                                      // ユーザーID
	User                User                 `gorm:"foreignKey:UserID;references:ID"`                                                                                 // ユーザーとのリレーション
	CreatedAt           time.Time            `gorm:"autoCreateTime"`
	UpdatedAt           time.Time            `gorm:"autoUpdateTime"` // 更新日時
}

type GofileVideoComment struct {
	ID            string    `gorm:"primaryKey;type:varchar(255)"`
	GofileVideoID string    `gorm:"type:varchar(255);not null"` // 外部キー制約つけたければ↓追加
	Comment       string    `gorm:"type:text;not null"`
	LikeCount     int       `gorm:"default:0"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"` // 更新日時
}

type GofileTag struct {
	ID           string        `gorm:"primaryKey;type:varchar(255)"`
	Name         string        `gorm:"type:varchar(100);not null"`                                                                                      // タグ名
	GofileVideos []GofileVideo `gorm:"many2many:gofile_video_tags;foreignKey:ID;joinForeignKey:GofileTagID;References:ID;joinReferences:GofileVideoID"` // GofileVideoとの多対多リレーション
	CreatedAt    time.Time     `gorm:"autoCreateTime"`
	UpdatedAt    time.Time     `gorm:"autoUpdateTime"` // 更新日時
}
