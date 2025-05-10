package repository

// import (
// 	"errors"
// 	"fmt"

// 	"gorm.io/gorm"

// 	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/database/model"
// 	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
// 	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/interactor"
// 	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"
// )

// type FileRepository struct {
// 	db *gorm.DB
// }

// func NewFileRepository(db *gorm.DB) output_port.FileRepository {
// 	return &FileRepository{db: db}
// }

// func (r FileRepository) DeleteBulk(tx interface{}, fileIDs []string) (err error) {
// 	defer output_port.WrapDatabaseError(&err)
// 	if txAsserted, ok := tx.(*gorm.DB); !ok {
// 		return output_port.ErrInvalidTransaction
// 	} else {
// 		return txAsserted.Model(&model.File{}).
// 			Where("file_id IN ?", fileIDs).
// 			Update("is_deleted", true).Error
// 	}
// }

// func (r FileRepository) FindByID(fileID string) (file entity.File, err error) {
// 	defer output_port.WrapDatabaseError(&err)
// 	var fileModel model.File
// 	if err := r.db.Model(&model.File{}).
// 		Where("file_id = ?", fileID).
// 		First(&fileModel).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
// 		return entity.File{}, fmt.Errorf("%w: file", interactor.ErrKind.NotFound)
// 	} else if err != nil {
// 		return entity.File{}, err
// 	}
// 	return fileModel.Entity(), nil
// }