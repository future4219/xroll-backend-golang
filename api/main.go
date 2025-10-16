package main

import (
	"fmt"
	"os"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/cache"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/database/repository"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/file"
	gofileAPI "gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/gofile"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/twitter"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/clock"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/email"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/config"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/log"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/aws"

	"go.uber.org/zap"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/authentication"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/database"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/ulid"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/api/router"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/interactor"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"
)

func main() {
	logger, err := log.NewLogger()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to create logger : %v\n", err)
		os.Exit(1)
	}

	db, err := database.NewMySQLDB(logger, true)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return
	}
	transaction := database.NewGormTransaction(db)

	if err := database.Migrate(db); err != nil {
		logger.Error("Failed to migrate database", zap.Error(err))
		return
	}

	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			logger.Error("Failed to get sql.DB", zap.Error(err))
		}
		err = sqlDB.Close()
		if err != nil {
			logger.Error("Failed to close database connection", zap.Error(err))
		}
	}()

	var awsCli *aws.Cli
	if config.IsAWSConfigFilled() {
		awsCli = aws.NewCli()
	}

	var mailDriver output_port.Email
	if config.EmailFrom() == "" {
		mailDriver = email.NewEmailDriverMock()
	} else {
		mailDriver = email.NewEmailDriver(awsCli)
	}

	var userAuth output_port.UserAuth
	if config.SigKey() == "" {
		userAuth = authentication.NewUserAuthMock()
	} else {
		userAuth = authentication.NewUserAuth()
	}

	cacheDriver := cache.GetInstance()
	fileDriver := file.NewFileDriver(awsCli, cacheDriver)
	clockDriver := clock.New()
	authCode := authentication.NewAuthenticationCode()

	registerVerification := repository.NewRegisterVerificationRepository(db)

	ulidDriver := ulid.NewULID()
	userRepo := repository.NewUserRepository(db, ulidDriver)
	userUC := interactor.NewAuthorizationUserUseCase(interactor.NewUserUseCase(
		clockDriver, mailDriver, ulidDriver, transaction, userAuth, userRepo, authCode, registerVerification))

	fileUC := interactor.NewAuthorizationFileUseCase(interactor.NewFileUseCase(ulidDriver, fileDriver))
	videoRepo := repository.NewVideoRepository(db, ulidDriver)
	videoUC := interactor.NewVideoUseCase(ulidDriver, videoRepo, clockDriver)

	twitter := twitter.NewTwitter()
	twitterUC := interactor.NewTwitterUseCase(twitter, videoRepo, ulidDriver)

	gofileRepo := repository.NewGofileRepository(db, ulidDriver)
	gofileAPIDriver := gofileAPI.NewGofileAPI()
	gofileUC := interactor.NewAuthorizationGofileUseCase(interactor.NewGofileUseCase(
		ulidDriver,
		gofileRepo,
		gofileAPIDriver,
		userRepo,
		clockDriver,
		twitter,
	))

	s := router.NewServer(
		userUC,
		fileUC,
		videoUC,
		twitterUC,
		gofileUC,
		true,
	)

	if err := s.Start(":8000"); err != nil {
		logger.Error("Failed to start server", zap.Error(err))
		os.Exit(1)
	}

}
