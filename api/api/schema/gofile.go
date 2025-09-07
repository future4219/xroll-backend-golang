package schema

import "gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"

type GofileVideoRes struct {
	ID                  string                  `json:"id"`                    // 動画のID
	Name                string                  `json:"name"`                  // 動画の名前
	GofileID            string                  `json:"gofile_id"`             // GofileのID
	GofileDirectURL     string                  `json:"gofile_direct_url"`     // GofileのダイレクトURL
	VideoURL            string                  `json:"video_url"`             // 動画のURL
	ThumbnailURL        string                  `json:"thumbnail_url"`         // サムネイルのURL
	LikeCount           int                     `json:"like_count"`            // いいねの数
	IsShared            bool                    `json:"is_shared"`             // 動画が共有されているかどうか
	GofileTags          []GofileTagRes          `json:"gofile_tags"`           // タグの情報
	GofileVideoComments []GofileVideoCommentRes `json:"gofile_video_comments"` // 動画に対するコメント
	UserID              *string                 `json:"user_id"`               // ユーザーID
	User                UserRes                 `json:"user"`                  // ユーザー情報
	CreatedAt           string                  `json:"created_at"`            // 作成日時
	UpdatedAt           string                  `json:"updated_at"`            // 更新日時
}

type GofileVideoListRes struct {
	Videos []GofileVideoRes `json:"videos"` // 動画のリスト
	Count  int              `json:"count"`  // 動画の総数
}

type GofileCreateReq struct {
	Name        string   `json:"name" validate:"required"`      // 動画の名前
	GofileID    string   `json:"gofile_id" validate:"required"` // GofileのID
	TagIDs      []string `json:"tag_ids"`                       // タグのIDリスト
	UserID      *string  `json:"user_id"`                       // ユーザーID
	GofileToken *string  `json:"gofile_token"`                  // Gofileのトークン
}

type GofileCreateRes struct {
	ID              string         `json:"id"`                // 動画のID
	Name            string         `json:"name"`              // 動画の名前
	GofileID        string         `json:"gofile_id"`         // GofileのID
	GofileDirectURL string         `json:"gofile_direct_url"` // GofileのダイレクトURL
	VideoURL        string         `json:"video_url"`         // 動画のURL
	ThumbnailURL    string         `json:"thumbnail_url"`     // サムネイルのURL
	UserID          *string        `json:"user_id"`           // ユーザーID
	GofileTags      []GofileTagRes `json:"gofile_tags"`       // タグの情報
}

type GofileTagRes struct {
	ID   string `json:"id"`   // タグのID
	Name string `json:"name"` // タグの名前
}

type GofileVideoCommentRes struct {
	ID        string `json:"id"`         // コメントのID
	Comment   string `json:"comment"`    // コメントの内容
	LikeCount int    `json:"like_count"` // いいねの数
	CreatedAt string `json:"created_at"` // 作成日時
	UpdatedAt string `json:"updated_at"` // 更新日時
}

type GofileUpdateIsShareReq struct {
	VideoID  string `json:"video_id" validate:"required"`  // 動画のID
	IsShared bool   `json:"is_shared" validate:"required"` // 動画の共有状態
}

func GofileCreateResFromEntity(e entity.GofileVideo) GofileCreateRes {
	tags := make([]GofileTagRes, 0, len(e.GofileTags))
	for _, t := range e.GofileTags {
		tags = append(tags, GofileTagResFromEntity(t))
	}
	return GofileCreateRes{
		ID:              e.ID,
		Name:            e.Name,
		GofileID:        e.GofileID,
		GofileDirectURL: e.GofileDirectURL,
		VideoURL:        e.VideoURL,
		ThumbnailURL:    e.ThumbnailURL,
		UserID:          &e.UserID,
		GofileTags:      tags,
	}
}

func GofileTagResFromEntity(tag entity.GofileTag) GofileTagRes {
	return GofileTagRes{
		ID:   tag.ID,
		Name: tag.Name,
	}
}

func GofileVideoListFromEntity(videos []entity.GofileVideo) GofileVideoListRes {
	res := make([]GofileVideoRes, 0, len(videos))
	for _, v := range videos {
		res = append(res, GofileVideoResFromEntity(v))
	}

	return GofileVideoListRes{
		Videos: res,
		Count:  len(videos),
	}
}

func GofileVideoResFromEntity(e entity.GofileVideo) GofileVideoRes {
	tags := make([]GofileTagRes, 0, len(e.GofileTags))
	for _, t := range e.GofileTags {
		tags = append(tags, GofileTagResFromEntity(t))
	}
	gofileComments := make([]GofileVideoCommentRes, 0, len(e.GofileVideoComments))
	for _, c := range e.GofileVideoComments {
		gofileComments = append(gofileComments, GofileVideoCommentRes{
			ID:        c.ID,
			Comment:   c.Comment,
			LikeCount: c.LikeCount,
			CreatedAt: c.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: c.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return GofileVideoRes{
		ID:                  e.ID,
		Name:                e.Name,
		// セキュリティ上、GofileIDとGofileDirectURLは返さない
		// GofileID:            e.GofileID,
		// GofileDirectURL:     e.GofileDirectURL,
		// VideoURL:            e.VideoURL,
		ThumbnailURL:        e.ThumbnailURL,
		LikeCount:           e.LikeCount,
		IsShared:            e.IsShared,
		GofileTags:          tags,
		GofileVideoComments: gofileComments,
		UserID:              &e.UserID,
		User:                UserResFromEntity(e.User),
		CreatedAt:           e.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:           e.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
