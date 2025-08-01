package repository

import (
	"fmt"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/database/model"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"
	"gorm.io/gorm"
)

type ThreadRepository struct {
	db   *gorm.DB
	ulid output_port.ULID
}

func NewThreadRepository(
	db *gorm.DB,
	ulid output_port.ULID,
) output_port.ThreadRepository {
	return &ThreadRepository{db: db, ulid: ulid}
}

func (r *ThreadRepository) Search(search output_port.ThreadSearch) (threads []entity.Thread, err error) {
	defer output_port.WrapDatabaseError(&err)

	var threadsModel []model.Thread
	if err = r.db.Model(&model.Thread{}).
		Preload("Comments").
		Limit(search.Limit).
		Offset(search.Offset).
		Order(fmt.Sprintf("%s %s", search.OrderBy.ToString(), search.Order.ToString())).
		Find(&threadsModel).Error; err != nil {
		return nil, err
	}

	threads = make([]entity.Thread, len(threadsModel))
	for i, threadModel := range threadsModel {
		threads[i] = threadModel.Entity()
	}
	return threads, nil
}

func (r *ThreadRepository) Create(thread entity.Thread) (err error) {
	defer output_port.WrapDatabaseError(&err)

	m := &model.Thread{
		ID:        thread.ID,
		Title:     thread.Title,
		LikeCount: thread.LikeCount,
		CreatedAt: thread.CreatedAt,
	}

	return r.db.Model(&model.Thread{}).Create(m).Error
}

func (r *ThreadRepository) FindByID(id string) (thread entity.Thread, err error) {
	defer output_port.WrapDatabaseError(&err)

	var threadModel model.Thread
	if err = r.db.Model(&model.Thread{}).
		Preload("Comments").
		Where("id = ?", id).
		First(&threadModel).Error; err != nil {
		return entity.Thread{}, err
	}
	thread = threadModel.Entity()
	return thread, nil
}

func (r *ThreadRepository) FindByIDs(ids []string) (threads []entity.Thread, err error) {
	defer output_port.WrapDatabaseError(&err)

	var threadsModel []model.Thread
	if err = r.db.Model(&model.Thread{}).
		Preload("Comments").
		Where("id IN ?", ids).
		Find(&threadsModel).Error; err != nil {
		return nil, err
	}
	threads = make([]entity.Thread, len(threadsModel))
	for i, threadModel := range threadsModel {
		threads[i] = threadModel.Entity()
	}

	return threads, nil
}

func (r *ThreadRepository) Update(thread entity.Thread) (err error) {
	defer output_port.WrapDatabaseError(&err)
	m := &model.Thread{
		ID:        thread.ID,
		Title:     thread.Title,
		LikeCount: thread.LikeCount,
	}
	if err = r.db.Model(&model.Thread{}).
		Where("id = ?", thread.ID).
		Updates(m).Error; err != nil {
		return err
	}

	return nil
}

func (r *ThreadRepository) CreateComment(threadID string, comment entity.ThreadComment) (err error) {
	defer output_port.WrapDatabaseError(&err)

	m := &model.ThreadComment{
		ID:         comment.ID,
		ThreadID:   threadID,
		ThreaderID: comment.ThreaderID,
		Comment:    comment.Comment,
	}
	if err = r.db.Model(&model.ThreadComment{}).
		Create(m).Error; err != nil {
		return err
	}

	return nil
}
