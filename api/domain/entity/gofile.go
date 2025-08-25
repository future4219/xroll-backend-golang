package entity

import "time"

type GofileVideo struct {
	ID                  string
	Name                string               // 動画の名前
	GofileID            string               // GofileのID
	GofileDirectURL     string               // Gofileの直接ダウンロードURL
	VideoURL            string               // 動画のURL
	ThumbnailURL        string               // 動画のサムネイルURL
	LikeCount           int                  // いいね数
	IsShared            bool                 // 動画が共有されているかどうか
	GofileTags          []GofileTag          // GofileTagとの多対多リレーション
	GofileVideoComments []GofileVideoComment // GofileVideoCommentとのリレーション
	UserID              string               // ユーザーID
	User                User                 // ユーザーとのリレーション
	CreatedAt           time.Time            // 作成日時
	UpdatedAt           time.Time            // 更新日時
}

type GofileVideoComment struct {
	ID            string    // コメントのID
	GofileVideoID string    // GofileVideoのID
	Comment       string    // コメント内容
	LikeCount     int       // いいね数
	CreatedAt     time.Time // 作成日時
	UpdatedAt     time.Time // 更新日時
}

type GofileTag struct {
	ID           string        // タグのID
	Name         string        // タグの名前
	GofileVideos []GofileVideo // GofileVideoとの多対多リレーション
}
