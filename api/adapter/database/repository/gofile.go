package repository

import (
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/database/model"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"
	"gorm.io/gorm"
)

type GofileRepository struct {
	db   *gorm.DB
	ulid output_port.ULID
}

func NewGofileRepository(
	db *gorm.DB,
	ulid output_port.ULID,
) output_port.GofileRepository {
	return &GofileRepository{db: db, ulid: ulid}
}

func (r *GofileRepository) Create(gofile entity.GofileVideo) error {
	var m model.GofileVideo

	gofileTagsModel := make([]model.GofileTag, 0, len(gofile.GofileTags))
	for _, tag := range gofile.GofileTags {
		gofileTagsModel = append(gofileTagsModel, model.GofileTag{
			ID: tag.ID,
		})
	}

	m = model.GofileVideo{
		ID:              gofile.ID,
		Name:            gofile.Name,
		GofileID:        gofile.GofileID,
		GofileDirectURL: gofile.GofileDirectURL,
		VideoURL:        gofile.VideoURL,
		ThumbnailURL:    gofile.ThumbnailURL,
		UserID:          gofile.UserID,
		GofileTags:      gofileTagsModel,
	}

	if err := r.db.Create(&m).Error; err != nil {
		return err
	}
	// 既存タグと関連づけ（join テーブル）
	if len(gofileTagsModel) > 0 {
		if err := r.db.Model(&m).Association("GofileTags").Replace(gofileTagsModel); err != nil {
			return err
		}
	}
	return nil
}

func (r *GofileRepository) FindByID(id string) (entity.GofileVideo, error) {
	var m model.GofileVideo
	if err := r.db.
		Preload("GofileTags").
		First(&m, "id = ?", id).Error; err != nil {
		return entity.GofileVideo{}, err
	}

	tags := make([]entity.GofileTag, 0, len(m.GofileTags))
	for _, tg := range m.GofileTags {
		tags = append(tags, entity.GofileTag{ID: tg.ID, Name: tg.Name})
	}

	return entity.GofileVideo{
		ID:              m.ID,
		Name:            m.Name,
		GofileID:        m.GofileID,
		GofileDirectURL: m.GofileDirectURL,
		VideoURL:        m.VideoURL,
		ThumbnailURL:    m.ThumbnailURL,
		UserID:          m.UserID,
		GofileTags:      tags,
	}, nil
}

func (r *GofileRepository) FindByUserID(userID string) ([]entity.GofileVideo, error) {
	var res []model.GofileVideo
	if err := r.db.
		Preload("GofileTags").
		Preload("GofileVideoComments").
		Where("user_id = ?", userID).
		Find(&res).Error; err != nil {
		return nil, err
	}

	videos := make([]entity.GofileVideo, len(res))
	for i, video := range res {
		
		tags := make([]entity.GofileTag, 0, len(video.GofileTags))
		for _, tg := range video.GofileTags {
			tags = append(tags, entity.GofileTag{ID: tg.ID, Name: tg.Name})
		}

		gofileVideoComments := make([]entity.GofileVideoComment, 0, len(video.GofileVideoComments))
		for _, comment := range video.GofileVideoComments {
			gofileVideoComments = append(gofileVideoComments, entity.GofileVideoComment{
				ID:        comment.ID,
				Comment:   comment.Comment,
				LikeCount: comment.LikeCount,
				CreatedAt: comment.CreatedAt,
				UpdatedAt: comment.UpdatedAt,
			})
		}

		videos[i] = entity.GofileVideo{
			ID:                  video.ID,
			Name:                video.Name,
			GofileID:            video.GofileID,
			GofileDirectURL:     video.GofileDirectURL,
			VideoURL:            video.VideoURL,
			ThumbnailURL:        video.ThumbnailURL,
			UserID:              video.UserID,
			GofileTags:          tags,
			GofileVideoComments: gofileVideoComments,
			CreatedAt:           video.CreatedAt,
			UpdatedAt:           video.UpdatedAt,
		}
	}

	return videos, nil
}
