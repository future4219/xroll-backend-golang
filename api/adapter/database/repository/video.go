package repository

import (
	"fmt"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/database/model"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VideoRepository struct {
	db   *gorm.DB
	ulid output_port.ULID
}

func NewVideoRepository(
	db *gorm.DB,
	ulid output_port.ULID,
) output_port.VideoRepository {
	return &VideoRepository{db: db, ulid: ulid}
}

func (r *VideoRepository) Search(search output_port.VideoSearch) (videos []entity.Video, err error) {
	defer output_port.WrapDatabaseError(&err)

	var videosModel []model.Video
	if err = r.db.Model(&model.Video{}).
		Preload("Comments").
		Limit(search.Limit).
		Offset(search.Offset).
		Where("created_at BETWEEN ? AND ?", search.Start, search.End).
		Order(fmt.Sprintf("%s %s", search.OrderBy.ToString(), search.Order.ToString())).
		Find(&videosModel).Error; err != nil {
		return nil, err
	}

	videos = make([]entity.Video, len(videosModel))
	for i, videoModel := range videosModel {
		videos[i] = videoModel.Entity()
	}
	return videos, nil
}

func (r *VideoRepository) Create(video entity.Video) (err error) {
	defer output_port.WrapDatabaseError(&err)

	m := &model.Video{
		ID:            video.ID,
		Ranking:       video.Ranking,
		VideoURL:      video.VideoURL,
		ThumbnailURL:  video.ThumbnailURL,
		TweetURL:      video.TweetURL,
		DownloadCount: video.DownloadCount,
		LikeCount:     video.LikeCount,
		CreatedAt:     video.CreatedAt,
	}
	if err = r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "video_url"}}, // UNIQUE指定カラム
		DoUpdates: clause.AssignmentColumns([]string{
			"ranking", "thumbnail_url", "tweet_url", "download_count", "like_count", "created_at",
		}),
	}).Create(m).Error; err != nil {
		return err
	}

	return nil
}

func (r *VideoRepository) CreateBulk(videos []entity.Video) (err error) {
	defer output_port.WrapDatabaseError(&err)

	// IDを生成する
	for i := range videos {
		videos[i].ID = r.ulid.GenerateID()
	}

	m := make([]*model.Video, len(videos))
	for i := range videos {
		m[i] = &model.Video{
			ID:            videos[i].ID,
			Ranking:       videos[i].Ranking,
			VideoURL:      videos[i].VideoURL,
			ThumbnailURL:  videos[i].ThumbnailURL,
			TweetURL:      videos[i].TweetURL,
			DownloadCount: videos[i].DownloadCount,
			LikeCount:     videos[i].LikeCount,
			CreatedAt:     videos[i].CreatedAt,
		}
	}
	if err = r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "video_url"}}, // ← これ必須！
		DoUpdates: clause.AssignmentColumns([]string{
			"ranking", "thumbnail_url", "tweet_url", "download_count", "like_count", "created_at",
		}),
	}).CreateInBatches(m, 100).Error; err != nil {
		return err
	}

	return nil
}

func (r *VideoRepository) FindByID(id string) (video entity.Video, err error) {
	defer output_port.WrapDatabaseError(&err)

	var videoModel model.Video
	if err = r.db.Model(&model.Video{}).
		Where("id = ?", id).
		First(&videoModel).Error; err != nil {
		return entity.Video{}, err
	}
	video = videoModel.Entity()
	return video, nil
}

func (r *VideoRepository) FindByIDs(ids []string) (videos []entity.Video, err error) {
	defer output_port.WrapDatabaseError(&err)

	var videosModel []model.Video
	if err = r.db.Model(&model.Video{}).
		Where("id IN ?", ids).
		Find(&videosModel).Error; err != nil {
		return nil, err
	}
	videos = make([]entity.Video, len(videosModel))
	for i, videoModel := range videosModel {
		videos[i] = videoModel.Entity()
	}

	return videos, nil
}

func (r *VideoRepository) Update(video entity.Video) (err error) {
	defer output_port.WrapDatabaseError(&err)
	fmt.Printf("video: %+v\n", video)
	m := &model.Video{
		ID:            video.ID,
		Ranking:       video.Ranking,
		VideoURL:      video.VideoURL,
		ThumbnailURL:  video.ThumbnailURL,
		TweetURL:      video.TweetURL,
		DownloadCount: video.DownloadCount,
		LikeCount:     video.LikeCount,
	}
	if err = r.db.Model(&model.Video{}).
		Where("id = ?", video.ID).
		Updates(m).Error; err != nil {
		return err
	}

	return nil
}

func (r *VideoRepository) CreateComment(videoID string, comment entity.Comment) (err error) {
	defer output_port.WrapDatabaseError(&err)

	m := &model.VideoComment{
		ID:      comment.ID,
		VideoID: videoID,
		Comment: comment.Comment,
	}
	if err = r.db.Model(&model.VideoComment{}).
		Create(m).Error; err != nil {
		return err
	}

	return nil
}
