package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/api/schema"
	log "gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/log"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/input_port"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/interactor"
	"go.uber.org/zap"
)

type ThreadHandler struct {
	ThreadUC input_port.IThreadUseCase
}

func NewThreadHandler(ThreadUC input_port.IThreadUseCase) *ThreadHandler {
	return &ThreadHandler{ThreadUC: ThreadUC}
}

func (h *ThreadHandler) Search(c echo.Context) error {
	logger, _ := log.NewLogger()

	req := &schema.ThreadSearchQueryReq{}
	if err := echo.QueryParamsBinder(c).
		Int("limit", &req.Limit).
		Int("offset", &req.Offset).
		BindError(); err != nil {
		logger.Error("Failed to bind query", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res, err := h.ThreadUC.Search(
		input_port.ThreadSearch{
			Limit:  req.Limit,
			Offset: req.Offset,
		})
	if err != nil {
		logger.Info("Failed to find ", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, schema.ThreadsResFromSearchResult(res, len(res)))
}

func (h *ThreadHandler) CreateBulk(c echo.Context) error {
	logger, _ := log.NewLogger()

	req := &schema.ThreadCreateBulkReq{}
	if err := c.Bind(req); err != nil {
		logger.Error("Failed to bind request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	threads, err := req.ToEntity()
	if err != nil {
		logger.Error("Failed to convert request to entity", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err := h.ThreadUC.CreateBulk(threads); err != nil {
		logger.Info("Failed to create thread", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.NoContent(http.StatusOK)
}

func (h *ThreadHandler) FindByID(c echo.Context) error {
	logger, _ := log.NewLogger()

	var threadID string
	if err := echo.PathParamsBinder(c).MustString("threadId", &threadID).BindError(); err != nil {
		logger.Info("Failed to bind path param id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	thread, err := h.ThreadUC.FindByID(threadID)
	if err != nil {
		logger.Info("Failed to find thread", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, schema.ThreadResFromEntity(thread))
}

func (h *ThreadHandler) FindByIDs(c echo.Context) error {
	logger, _ := log.NewLogger()

	idsParam := c.QueryParam("ids")
	idList := strings.Split(idsParam, ",")

	if len(idList) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "no ids provided"})
	}

	threads, err := h.ThreadUC.FindByIDs(idList)
	if err != nil {
		logger.Info("Failed to find threads", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, schema.ThreadsResFromSearchResult(threads, len(threads)))
}

func (h *ThreadHandler) Like(c echo.Context) error {
	logger, _ := log.NewLogger()

	var threadID string
	if err := echo.PathParamsBinder(c).MustString("threadId", &threadID).BindError(); err != nil {
		logger.Info("Failed to bind path param id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.ThreadUC.Like(threadID); err != nil {
		logger.Info("Failed to like thread", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.NoContent(http.StatusOK)
}

func (h *ThreadHandler) Comment(c echo.Context) error {
	logger, _ := log.NewLogger()

	var threadID string
	if err := echo.PathParamsBinder(c).MustString("threadId", &threadID).BindError(); err != nil {
		logger.Info("Failed to bind path param id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	req := &schema.ThreadCommentReq{}
	if err := c.Bind(req); err != nil {
		logger.Error("Failed to bind request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.ThreadUC.Comment(threadID, req.Comment); err != nil {
		logger.Info("Failed to comment on thread", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.NoContent(http.StatusOK)
}
