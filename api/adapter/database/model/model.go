package model

import (
	"time"
)

type User struct {
	ID       string `gorm:"primaryKey;type:varchar(255)"`
	Name     string `gorm:"type:varchar(100)"`
	Age      int
	UserType string `gorm:"type:varchar(20)"`
}

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

type Thread struct {
	ID        string          `gorm:"primaryKey;type:varchar(255)"` // スレッドID
	Title     string          `gorm:"type:varchar(255);not null"`   // スレッドタイトル
	LikeCount int             `gorm:"default:0"`                    // いいね数
	Comments  []ThreadComment `gorm:"foreignKey:ThreadID"`          // スレッドに紐づくコメント
	CreatedAt time.Time       `gorm:"autoCreateTime"`               // 作成日時
}

type ThreadComment struct {
	ID         string    `gorm:"primaryKey;type:varchar(255)"` // コメントID
	ThreadID   string    `gorm:"type:varchar(255);not null;foreignKey:ThreadID"`
	ThreaderID string    `gorm:"type:varchar(8);not null"` // コメントしたユーザーのID
	Comment    string    `gorm:"type:text;not null"`       // コメント内容
	LikeCount  int       `gorm:"default:0"`                // いいね数
	CreatedAt  time.Time `gorm:"autoCreateTime"`           // 作成日時
}
