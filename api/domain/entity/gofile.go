package entity

type GofileVideo struct {
	ID string
	Name                string // 動画の名前
	GofileID            string // GofileのID
	GofileDirectURL     string // Gofileの直接ダウンロードURL
	VideoURL            string // 動画のURL
	LikeCount           int    // いいね数
	GofileTags          []GofileTag          // GofileTagとの多対多リレーション
	GofileVideoComments []GofileVideoComment // GofileVideoCommentとのリレーション
	UserID              string    // ユーザーID
	User                User      // ユーザーとのリレーション
}	

type GofileVideoComment struct {
	ID            string // コメントのID
	GofileVideoID string // GofileVideoのID
	Comment       string // コメント内容
	LikeCount     int    // いいね数
	CreatedAt     string // 作成日時
	UpdatedAt     string // 更新日時
}

type GofileTag struct {	
	ID           string // タグのID
	Name         string // タグの名前
	GofileVideos []GofileVideo // GofileVideoとの多対多リレーション
}
