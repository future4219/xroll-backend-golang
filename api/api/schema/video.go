package schema

import (
	"time"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
)

type VideoRes struct {
	ID            string       `json:"id"`
	Ranking       int          `json:"ranking"`
	VideoURL      string       `json:"video_url"`
	ThumbnailURL  string       `json:"thumbnail_url"`
	TweetURL      string       `json:"tweet_url"`
	DownloadCount int          `json:"download_count"`
	LikeCount     int          `json:"like_count"`
	Comments      []CommentRes `json:"comments"`
	CreatedAt     time.Time    `json:"created_at"`
}

type CommentRes struct {
	ID        string    `json:"id"`
	Comment   string    `json:"comment"`
	LikeCount int       `json:"like_count"`
	CreatedAt time.Time `json:"created_at"`
}

type VideosRes struct {
	Videos []VideoRes `json:"videos"`
	Total  int        `json:"total"`
}

type VideoSearchQueryReq struct {
	Limit      int  `query:"limit"`
	Offset     int  `query:"offset"`
	IsRealtime bool `query:"is_realtime"`
}

type VideoCreateBulkReq struct {
	Videos []VideoCreateReq `json:"videos"`
}

type VideoCreateReq struct {
	ID            string  `json:"id"`
	Ranking       int     `json:"ranking"`
	VideoURL      string  `json:"video_url"`
	ThumbnailURL  string  `json:"thumbnail_url"`
	DownloadCount int     `json:"download_count"`
	TweetURL      *string `json:"tweet_url"`
}

type VideoCommentReq struct {
	Comment string `json:"comment"`
}

func (vcq *VideoCreateBulkReq) ToEntity() ([]entity.Video, error) {
	videos := make([]entity.Video, len(vcq.Videos))
	for i, video := range vcq.Videos {
		videos[i] = entity.Video{
			ID:            video.ID,
			Ranking:       video.Ranking,
			VideoURL:      video.VideoURL,
			ThumbnailURL:  video.ThumbnailURL,
			DownloadCount: video.DownloadCount,
			TweetURL:      video.TweetURL,
		}
	}
	return videos, nil
}

func VideoCreateReqFromEntity(video entity.Video) VideoCreateReq {
	return VideoCreateReq{
		ID:            video.ID,
		Ranking:       video.Ranking,
		VideoURL:      video.VideoURL,
		ThumbnailURL:  video.ThumbnailURL,
		DownloadCount: video.DownloadCount,
	}
}

func VideoCreateReqsFromEntity(videos []entity.Video) []VideoCreateReq {
	res := make([]VideoCreateReq, len(videos))
	for i, video := range videos {
		res[i] = VideoCreateReqFromEntity(video)
	}
	return res
}

func VideoResFromEntity(video entity.Video) VideoRes {
	return VideoRes{
		ID:           video.ID,
		Ranking:      video.Ranking,
		VideoURL:     video.VideoURL,
		ThumbnailURL: video.ThumbnailURL,
		TweetURL: func() string {
			if video.TweetURL != nil {
				return *video.TweetURL
			}
			return ""
		}(),
		DownloadCount: video.DownloadCount,
		LikeCount:     video.LikeCount,
		Comments:      CommentsResFromEntity(video.Comments),
		CreatedAt:     video.CreatedAt,
	}
}

func VideosResFromSearchResult(list []entity.Video, total int) VideosRes {
	return VideosRes{
		Videos: VideosResFromEntity(list),
		Total:  total,
	}
}

func VideosResFromEntity(videos []entity.Video) []VideoRes {
	res := make([]VideoRes, len(videos))
	for i, video := range videos {
		res[i] = VideoResFromEntity(video)
	}
	return res
}

func CommentsResFromEntity(comments []entity.Comment) []CommentRes {
	res := make([]CommentRes, len(comments))
	for i, comment := range comments {
		res[i] = CommentRes{
			ID:        comment.ID,
			Comment:   comment.Comment,
			LikeCount: comment.LikeCount,
			CreatedAt: comment.CreatedAt,
		}
	}
	return res
}
