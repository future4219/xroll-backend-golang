package repository

import (
	"fmt"
	"strings"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/database/model"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entconst"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/interactor"
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
		IsShared:        gofile.IsShared,
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

func (r *GofileRepository) Update(gofile entity.GofileVideo) (err error) {
	defer output_port.WrapDatabaseError(&err)

	var m model.GofileVideo

	gofileTagsModel := make([]model.GofileTag, 0, len(gofile.GofileTags))
	for _, tag := range gofile.GofileTags {
		gofileTagsModel = append(gofileTagsModel, model.GofileTag{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}
	m = model.GofileVideo{
		ID:              gofile.ID,
		Name:            gofile.Name,
		GofileID:        gofile.GofileID,
		GofileDirectURL: gofile.GofileDirectURL,
		VideoURL:        gofile.VideoURL,
		ThumbnailURL:    gofile.ThumbnailURL,
		Description:     gofile.Description,
		PlayCount:       gofile.PlayCount,
		LikeCount:       gofile.LikeCount,
		IsShared:        gofile.IsShared,
		UserID:          gofile.UserID,
		GofileTags:      gofileTagsModel,
		CreatedAt:       gofile.CreatedAt,
		UpdatedAt:       gofile.UpdatedAt,
	}
	fmt.Printf("gofile: %+v\n", m)
	if err := r.db.Save(&m).Error; err != nil {
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
		Preload("GofileVideoComments").
		Preload("User").
		Where("is_deleted = ?", false).
		First(&m, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.GofileVideo{}, fmt.Errorf("%w: gofile video", interactor.ErrKind.NotFound)
		}
		return entity.GofileVideo{}, err
	}

	tags := make([]entity.GofileTag, 0, len(m.GofileTags))
	for _, tg := range m.GofileTags {
		tags = append(tags, entity.GofileTag{ID: tg.ID, Name: tg.Name})
	}

	gofileVideoComments := make([]entity.GofileVideoComment, 0, len(m.GofileVideoComments))
	for _, comment := range m.GofileVideoComments {
		gofileVideoComments = append(gofileVideoComments, entity.GofileVideoComment{
			ID:        comment.ID,
			Comment:   comment.Comment,
			LikeCount: comment.LikeCount,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		})
	}

	likeCount := int64(0)
	if err := r.db.Model(&model.GofileVideoLike{}).Where("gofile_video_id = ?", m.ID).Count(&likeCount).Error; err != nil {
		return entity.GofileVideo{}, err
	}
	m.LikeCount = int(likeCount)

	return entity.GofileVideo{
		ID:                  m.ID,
		Name:                m.Name,
		GofileID:            m.GofileID,
		GofileDirectURL:     m.GofileDirectURL,
		VideoURL:            m.VideoURL,
		ThumbnailURL:        m.ThumbnailURL,
		Description:         m.Description,
		PlayCount:           m.PlayCount,
		LikeCount:           m.LikeCount,
		IsShared:            m.IsShared,
		GofileTags:          tags,
		GofileVideoComments: gofileVideoComments,
		UserID:              m.UserID,
		User:                m.User.Entity(),
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
	}, nil
}

func (r *GofileRepository) FindByUserID(userID string) ([]entity.GofileVideo, error) {
	var res []model.GofileVideo
	if err := r.db.
		Preload("User").
		Preload("GofileTags").
		Preload("GofileVideoComments").
		Where("user_id = ? AND is_deleted = ?", userID, false).
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
			Description:         video.Description,
			PlayCount:           video.PlayCount,
			LikeCount:           video.LikeCount,
			IsShared:            video.IsShared,
			UserID:              video.UserID,
			User:                video.User.Entity(),
			GofileTags:          tags,
			GofileVideoComments: gofileVideoComments,
			CreatedAt:           video.CreatedAt,
			UpdatedAt:           video.UpdatedAt,
		}
	}

	return videos, nil
}

func (r *GofileRepository) FindByUserIDShared(userId string) ([]entity.GofileVideo, error) {
	var res []model.GofileVideo
	fmt.Println("--------------------------")
	if err := r.db.
		Preload("User").
		Preload("GofileTags").
		Preload("GofileVideoComments").
		Where("user_id = ? AND is_shared = ? AND is_deleted = ?", userId, true, false).
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
			Description:         video.Description,
			PlayCount:           video.PlayCount,
			LikeCount:           video.LikeCount,
			IsShared:            video.IsShared,
			UserID:              video.UserID,
			User:                video.User.Entity(),
			GofileTags:          tags,
			GofileVideoComments: gofileVideoComments,
			CreatedAt:           video.CreatedAt,
			UpdatedAt:           video.UpdatedAt,
		}
	}

	return videos, nil
}

func (r *GofileRepository) Delete(id string) error {
	return r.db.Model(&model.GofileVideo{}).
		Where("id = ?", id).
		Update("is_deleted", true).Error
}

func (r *GofileRepository) HasLike(userID, videoID string) (bool, error) {
	var cnt int64
	if err := r.db.
		Model(&model.GofileVideoLike{}).
		Where("user_id = ? AND gofile_video_id = ?", userID, videoID).
		Count(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (r *GofileRepository) CreateLike(l entity.GofileVideoLike) error {
	m := model.GofileVideoLike{
		ID:            l.ID,
		GofileVideoID: l.GofileVideoID,
		UserID:        l.UserID,
		CreatedAt:     l.CreatedAt,
	}
	return r.db.Create(&m).Error
}

func (r *GofileRepository) DeleteLike(userID, videoID string) (int64, error) {
	tx := r.db.
		Where("user_id = ? AND gofile_video_id = ?", userID, videoID).
		Delete(&model.GofileVideoLike{})
	return tx.RowsAffected, tx.Error
}

func (r *GofileRepository) FindLikedVideos(userID string) ([]entity.GofileVideo, error) {
	var rows []model.GofileVideo

	q := r.db.
		Model(&model.GofileVideo{}).
		Joins("JOIN gofile_video_likes l ON l.gofile_video_id = gofile_videos.id AND l.user_id = ?", userID).
		Where("gofile_videos.is_deleted = ?", false).
		Where("(gofile_videos.is_shared = ? OR gofile_videos.user_id = ?)", true, userID).
		Select("gofile_videos.*"). // 余計な列を拾わない
		Preload("User").
		Preload("GofileTags").
		Preload("GofileVideoComments").
		Order("l.created_at DESC") // ← DISTINCTを外したのでOK

	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}

	videos := make([]entity.GofileVideo, len(rows))
	for i, v := range rows {
		tags := make([]entity.GofileTag, 0, len(v.GofileTags))
		for _, tg := range v.GofileTags {
			tags = append(tags, entity.GofileTag{ID: tg.ID, Name: tg.Name})
		}
		comments := make([]entity.GofileVideoComment, 0, len(v.GofileVideoComments))
		for _, c := range v.GofileVideoComments {
			comments = append(comments, entity.GofileVideoComment{
				ID:        c.ID,
				Comment:   c.Comment,
				LikeCount: c.LikeCount,
				CreatedAt: c.CreatedAt,
				UpdatedAt: c.UpdatedAt,
			})
		}
		videos[i] = entity.GofileVideo{
			ID:                  v.ID,
			Name:                v.Name,
			GofileID:            v.GofileID,
			GofileDirectURL:     v.GofileDirectURL,
			VideoURL:            v.VideoURL,
			ThumbnailURL:        v.ThumbnailURL,
			Description:         v.Description,
			PlayCount:           v.PlayCount,
			LikeCount:           v.LikeCount,
			IsShared:            v.IsShared,
			UserID:              v.UserID,
			User:                v.User.Entity(),
			GofileTags:          tags,
			GofileVideoComments: comments,
			CreatedAt:           v.CreatedAt,
			UpdatedAt:           v.UpdatedAt,
		}
	}
	return videos, nil
}

// Search: 公開(Shared) & 未削除を対象に、名前/説明/タグ名で横断検索。
// 並び順は entconst の enum に従う。Skip/Limit ページング対応。
func (r *GofileRepository) Search(qry output_port.GofileSearchQuery) ([]entity.GofileVideo, error) {
	// ---- base query ----
	dbq := r.db.
		Model(&model.GofileVideo{}).
		Joins("LEFT JOIN gofile_video_tags gvt ON gvt.gofile_video_id = gofile_videos.id").
		Joins("LEFT JOIN gofile_tags gt ON gt.id = gvt.gofile_tag_id").
		Where("gofile_videos.is_deleted = ?", false).
		Where("gofile_videos.is_shared = ?", true).
		Select("DISTINCT gofile_videos.*").
		Preload("User").
		Preload("GofileTags").
		Preload("GofileVideoComments")

	// ---- keyword (case-insensitive LIKE) ----
	if q := strings.TrimSpace(qry.Q); q != "" {
		like := "%" + strings.ToLower(q) + "%"
		dbq = dbq.Where(`
			LOWER(gofile_videos.name) LIKE ? OR
			LOWER(gofile_videos.description) LIKE ? OR
			LOWER(gt.name) LIKE ?
		`, like, like, like)
	}

	// ---- ordering (enum-safe) ----
	col := "gofile_videos.updated_at" // default
	switch qry.OrderBy {
	case entconst.GofileOrderByCreatedAt:
		col = "gofile_videos.created_at"
	case entconst.GofileOrderByUpdatedAt:
		col = "gofile_videos.updated_at"
	case entconst.GofileOrderByLikeCount:
		col = "gofile_videos.like_count"
	case entconst.GofileOrderByPlayCount:
		col = "gofile_videos.play_count"
	}

	dir := "DESC"
	switch qry.Order {
	case entconst.ASC:
		dir = "ASC"
	case entconst.DESC:
		dir = "DESC"
	}
	dbq = dbq.Order(col + " " + dir)

	// ---- paging ----
	if qry.Skip > 0 {
		dbq = dbq.Offset(qry.Skip)
	}
	if qry.Limit > 0 {
		dbq = dbq.Limit(qry.Limit)
	} else {
		dbq = dbq.Limit(60) // 安全な上限（必要なら調整）
	}

	// ---- execute ----
	var rows []model.GofileVideo
	if err := dbq.Find(&rows).Error; err != nil {
		return nil, err
	}

	// ---- model -> entity ----
	videos := make([]entity.GofileVideo, len(rows))
	for i, v := range rows {
		tags := make([]entity.GofileTag, 0, len(v.GofileTags))
		for _, tg := range v.GofileTags {
			tags = append(tags, entity.GofileTag{ID: tg.ID, Name: tg.Name})
		}
		comments := make([]entity.GofileVideoComment, 0, len(v.GofileVideoComments))
		for _, c := range v.GofileVideoComments {
			comments = append(comments, entity.GofileVideoComment{
				ID:        c.ID,
				Comment:   c.Comment,
				LikeCount: c.LikeCount,
				CreatedAt: c.CreatedAt,
				UpdatedAt: c.UpdatedAt,
			})
		}
		videos[i] = entity.GofileVideo{
			ID:                  v.ID,
			Name:                v.Name,
			GofileID:            v.GofileID,
			GofileDirectURL:     v.GofileDirectURL,
			VideoURL:            v.VideoURL,
			ThumbnailURL:        v.ThumbnailURL,
			Description:         v.Description,
			PlayCount:           v.PlayCount,
			LikeCount:           v.LikeCount,
			IsShared:            v.IsShared,
			UserID:              v.UserID,
			User:                v.User.Entity(),
			GofileTags:          tags,
			GofileVideoComments: comments,
			CreatedAt:           v.CreatedAt,
			UpdatedAt:           v.UpdatedAt,
		}
	}

	return videos, nil
}
