package repository

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/database/model"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/interactor"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"
)

type RegisterVerificationRepository struct {
	db *gorm.DB
}

func NewRegisterVerificationRepository(db *gorm.DB) output_port.RegisterVerificationRepository {
	return RegisterVerificationRepository{db: db}
}

func (r RegisterVerificationRepository) FindByEmail(email string) (_ entity.RegisterVerification, err error) {
	defer output_port.WrapDatabaseError(&err)
	var registerVerification model.RegisterVerification
	err = r.db.Model(&model.RegisterVerification{}).Where("email = ?", email).First(&registerVerification).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.RegisterVerification{}, interactor.ErrKind.NotFound
	}
	if err != nil {
		return entity.RegisterVerification{}, err
	}
	return registerVerification.Entity(), nil
}

func (r RegisterVerificationRepository) UpsertInTx(tx interface{}, entity entity.RegisterVerification) (err error) {
	defer output_port.WrapDatabaseError(&err)

	if txAsserted, ok := tx.(*gorm.DB); !ok {
		return output_port.ErrInvalidTransaction
	} else {
		return r.upsert(entity, txAsserted)
	}
}

func (r RegisterVerificationRepository) upsert(entity entity.RegisterVerification, db *gorm.DB) (err error) {
	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "Email"}},
		DoUpdates: clause.AssignmentColumns([]string{"hashed_authentication_code", "expires_at", "hashed_password"}),
	}).Create(&model.RegisterVerification{
		Email:                    entity.Email,
		ExpiresAt:                entity.ExpiresAt,
		HashedPassword:           entity.HashedPassword,
		HashedAuthenticationCode: entity.HashedAuthenticationCode,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r RegisterVerificationRepository) DeleteByEmailInTx(tx interface{}, email string) (err error) {
	defer output_port.WrapDatabaseError(&err)
	txAsserted, ok := tx.(*gorm.DB)
	if !ok {
		return output_port.ErrInvalidTransaction
	}
	return txAsserted.Where("email = ?", email).Delete(&model.RegisterVerification{}).Error
}
