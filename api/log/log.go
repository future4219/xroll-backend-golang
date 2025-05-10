package log

import (
	"go.uber.org/zap"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/config"
)

var (
	logger *zap.Logger
)

func NewLogger() (*zap.Logger, error) {
	if logger != nil {
		return logger, nil
	}

	var err error

	if config.IsDevelopment() {
		logger, err = zap.NewDevelopment()
	} else if config.IsTest() {
		logger = zap.NewNop()
	} else {
		logger, err = zap.NewProduction()
	}

	return logger, err
}
