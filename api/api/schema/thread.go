package schema

import (
	"time"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
)

type ThreadRes struct {
	ID            string       `json:"id"`
	Ranking       int          `json:"ranking"`
	ThreadURL      string       `json:"thread_url"`
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

type ThreadsRes struct {
	Threads []ThreadRes `json:"threads"`
	Total  int        `json:"total"`
}

type ThreadSearchQueryReq struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

type ThreadCreateBulkReq struct {
	Threads []ThreadCreateReq `json:"threads"`
}

type ThreadCreateReq struct {
	ID            string  `json:"id"`
	Ranking       int     `json:"ranking"`
	ThreadURL      string  `json:"thread_url"`
	ThumbnailURL  string  `json:"thumbnail_url"`
	DownloadCount int     `json:"download_count"`
	TweetURL      *string `json:"tweet_url"`
}

type ThreadCommentReq struct {
	Comment string `json:"comment"`
}

func (vcq *ThreadCreateBulkReq) ToEntity() ([]entity.Thread, error) {
	threads := make([]entity.Thread, len(vcq.Threads))
	for i, thread := range vcq.Threads {
		threads[i] = entity.Thread{
			ID:            thread.ID,
			Ranking:       thread.Ranking,
			ThreadURL:      thread.ThreadURL,
			ThumbnailURL:  thread.ThumbnailURL,
			DownloadCount: thread.DownloadCount,
			TweetURL:      thread.TweetURL,
		}
	}
	return threads, nil
}

func ThreadCreateReqFromEntity(thread entity.Thread) ThreadCreateReq {
	return ThreadCreateReq{
		ID:            thread.ID,
		Ranking:       thread.Ranking,
		ThreadURL:      thread.ThreadURL,
		ThumbnailURL:  thread.ThumbnailURL,
		DownloadCount: thread.DownloadCount,
	}
}

func ThreadCreateReqsFromEntity(threads []entity.Thread) []ThreadCreateReq {
	res := make([]ThreadCreateReq, len(threads))
	for i, thread := range threads {
		res[i] = ThreadCreateReqFromEntity(thread)
	}
	return res
}

func ThreadResFromEntity(thread entity.Thread) ThreadRes {
	return ThreadRes{
		ID:           thread.ID,
		Ranking:      thread.Ranking,
		ThreadURL:     thread.ThreadURL,
		ThumbnailURL: thread.ThumbnailURL,
		TweetURL: func() string {
			if thread.TweetURL != nil {
				return *thread.TweetURL
			}
			return ""
		}(),
		DownloadCount: thread.DownloadCount,
		LikeCount:     thread.LikeCount,
		Comments:      CommentsResFromEntity(thread.Comments),
		CreatedAt:     thread.CreatedAt,
	}
}

func ThreadsResFromSearchResult(list []entity.Thread, total int) ThreadsRes {
	return ThreadsRes{
		Threads: ThreadsResFromEntity(list),
		Total:  total,
	}
}

func ThreadsResFromEntity(threads []entity.Thread) []ThreadRes {
	res := make([]ThreadRes, len(threads))
	for i, thread := range threads {
		res[i] = ThreadResFromEntity(thread)
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
