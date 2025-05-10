package database

import (
	"gorm.io/gorm"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"
)

func NewGormTransaction(db *gorm.DB) output_port.Transaction {
	return &GormTransaction{
		db: db,
	}
}

type GormTransaction struct {
	db *gorm.DB
}

func (t *GormTransaction) StartTransaction(function func(tx interface{}) error) error {
	return t.db.Transaction(func(gormTx *gorm.DB) error {
		return function(gormTx)
	})
}
