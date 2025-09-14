package database

import (
	"fmt"
	"math"
	"time"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/database/model"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/config"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func NewMySQLDB(logger *zap.Logger, isLogging bool) (*gorm.DB, error) {
	dsn := config.DSN()

	gormConfig := &gorm.Config{}
	if !isLogging {
		gormConfig.Logger = gormLogger.Default.LogMode(gormLogger.Silent)
	}
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)

	if err != nil {
		return nil, fmt.Errorf("failed to open MySQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetConnMaxIdleTime(100)
	sqlDB.SetMaxOpenConns(100)

	// Check connection
	const retryMax = 10
	for i := 0; i < retryMax; i++ {
		err = sqlDB.Ping()
		if err == nil {
			break
		}
		if i == retryMax-1 {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
		duration := time.Millisecond * time.Duration(math.Pow(1.5, float64(i))*1000)
		logger.Warn("failed to connect to database retrying", zap.Error(err), zap.Duration("sleepSeconds", duration))
		time.Sleep(duration)
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	if err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&model.User{},
		&model.RegisterVerification{},
		&model.Video{},
		&model.VideoComment{},
		&model.GofileVideo{},
		&model.GofileVideoComment{},
		&model.GofileTag{},
		&model.GofileVideoLike{},
	); err != nil {
		return err
	}
	return nil
}

func DropDB(db *gorm.DB) {
	_ = db.Migrator().DropTable(
		&model.GofileVideoLike{},
		&model.User{},
		&model.RegisterVerification{},
		&model.Video{},
		&model.VideoComment{},
		&model.GofileVideo{},
		&model.GofileVideoComment{},
		&model.GofileTag{},
	)
}
