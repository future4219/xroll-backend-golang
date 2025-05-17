package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	log "gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/log"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/input_port"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/interactor"
	"go.uber.org/zap"
)

type TwitterHandler struct {
	TwitterUC input_port.ITwitterUseCase
}

func NewTwitterHandler(TwitterUC input_port.ITwitterUseCase) *TwitterHandler {
	return &TwitterHandler{TwitterUC: TwitterUC}
}

func (h *TwitterHandler) GetVideoByURL(c echo.Context) error {
	logger, _ := log.NewLogger()

	url := c.QueryParam("url")
	if url == "" {
		logger.Error("URL is empty")
		return echo.NewHTTPError(http.StatusBadRequest, "URL is empty")
	}

	video, err := h.TwitterUC.GetVideoByURL(url)
	if err != nil {
		logger.Info("Failed to get twitter video by url ", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, video)
}
