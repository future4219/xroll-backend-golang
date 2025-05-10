package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/api/middleware"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/api/schema"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
	log "gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/log"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/input_port"
)

type FileHandler struct {
	FileUC input_port.IFileUseCase
}

func NewFileHandler(FileUC input_port.IFileUseCase) *FileHandler {
	return &FileHandler{FileUC: FileUC}
}

func (h *FileHandler) IssuePreSignedURLForPut(c echo.Context) error {
	logger, _ := log.NewLogger()

	ctx := c.Request().Context()
	user, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		logger.Error("Failed to get user from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	url, key, err := h.FileUC.IssuePreSignedURLForPut(entity.User{
		ID:   user.ID,
		UserType: user.UserType,
	})
	if err != nil {
		logger.Info("Failed to issue pre-signed url for put", zap.Error(err))
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, &schema.IssuePreSignedURLForPutRes{
		PreSignedUrl: url,
		Key:          key,
	})
}

func (h *FileHandler) IssuePreSignedURLForPutVideo(c echo.Context) error {
	logger, _ := log.NewLogger()

	ctx := c.Request().Context()
	user, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		logger.Error("Failed to get user from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	req := &schema.IssuePreSignedURLForVideoReq{}
	if err := c.Bind(req); err != nil {
		logger.Error("Failed to bind request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	url, key, err := h.FileUC.IssuePreSignedURLForPutVideo(entity.User{
		ID:   user.ID,
		UserType: user.UserType,
	}, req.FileName)
	if err != nil {
		logger.Info("Failed to issue pre-signed url for put", zap.Error(err))
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, &schema.IssuePreSignedURLForPutRes{
		PreSignedUrl: url,
		Key:          key,
	})
}

func (h *FileHandler) IssuePresignedURLForGetVideo(c echo.Context) error {
	logger, _ := log.NewLogger()

	ctx := c.Request().Context()
	user, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		logger.Error("Failed to get user from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var fileID string
	var fileName string
	if err := echo.PathParamsBinder(c).MustString("fileId", &fileID).MustString("fileName", &fileName).BindError(); err != nil {
		logger.Error("Failed to bind path param id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	url, status, err := h.FileUC.IssuePreSignedURLForGetVideo(entity.User{
		ID:   user.ID,
		UserType: user.UserType,
	}, fileName, fileID)
	if err != nil {
		logger.Info("Failed to issue pre-signed url for get", zap.Error(err))
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, &schema.IssuePreSignedURLForGetRes{
		PreSignedUrl: url,
		Status:       status.String(),
	})
}
