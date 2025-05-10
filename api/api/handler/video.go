package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/api/schema"
	log "gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/log"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/input_port"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/interactor"
	"go.uber.org/zap"
)

type VideoHandler struct {
	VideoUC input_port.IVideoUseCase
}

func NewVideoHandler(VideoUC input_port.IVideoUseCase) *VideoHandler {
	return &VideoHandler{VideoUC: VideoUC}
}

func (h *VideoHandler) Search(c echo.Context) error {
	logger, _ := log.NewLogger()

	//クエリパラメータを取得

	req := &schema.VideoSearchQueryReq{}
	if err := echo.QueryParamsBinder(c).
		Int("limit", &req.Limit).
		Int("offset", &req.Offset).
		BindError(); err != nil {
		logger.Error("Failed to bind query", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	fmt.Printf("🐛 Search called with limit: %d, offset: %d\n", req.Limit, req.Offset) // ← これ表示されるか？
	res, err := h.VideoUC.Search(
		input_port.VideoSearch{
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
	return c.JSON(http.StatusOK, schema.VideosResFromSearchResult(res, len(res)))
}

func (h *VideoHandler) CreateBulk(c echo.Context) error {
	logger, _ := log.NewLogger()

	req := &schema.VideoCreateBulkReq{}
	if err := c.Bind(req); err != nil {
		logger.Error("Failed to bind request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	videos, err := req.ToEntity()
	if err != nil {
		logger.Error("Failed to convert request to entity", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err := h.VideoUC.CreateBulk(videos); err != nil {
		logger.Info("Failed to create video", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.NoContent(http.StatusOK)
}

func (h *VideoHandler) FindByID(c echo.Context) error {
	logger, _ := log.NewLogger()

	var videoID string
	if err := echo.PathParamsBinder(c).MustString("videoId", &videoID).BindError(); err != nil {
		logger.Info("Failed to bind path param id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	video, err := h.VideoUC.FindByID(videoID)
	if err != nil {
		logger.Info("Failed to find video", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, schema.VideoResFromEntity(video))
}

func (h *VideoHandler) FindByIDs(c echo.Context) error {
	logger, _ := log.NewLogger()
	
	idsParam := c.QueryParam("ids")
	idList := strings.Split(idsParam, ",")
	fmt.Printf("🐛 FindByIDs called with ids: %s\n", idsParam) // ← これ表示されるか？
	if len(idList) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "no ids provided"})
	}

	videos, err := h.VideoUC.FindByIDs(idList)
	if err != nil {
		logger.Info("Failed to find videos", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, schema.VideosResFromSearchResult(videos, len(videos)))
}

func (h *VideoHandler) Like(c echo.Context) error {
	logger, _ := log.NewLogger()

	var videoID string
	if err := echo.PathParamsBinder(c).MustString("videoId", &videoID).BindError(); err != nil {
		logger.Info("Failed to bind path param id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.VideoUC.Like(videoID); err != nil {
		logger.Info("Failed to like video", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.NoContent(http.StatusOK)
}

func (h *VideoHandler) Comment(c echo.Context) error {
	logger, _ := log.NewLogger()

	var videoID string
	if err := echo.PathParamsBinder(c).MustString("videoId", &videoID).BindError(); err != nil {
		logger.Info("Failed to bind path param id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	req := &schema.VideoCommentReq{}
	if err := c.Bind(req); err != nil {
		logger.Error("Failed to bind request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.VideoUC.Comment(videoID, req.Comment); err != nil {
		logger.Info("Failed to comment on video", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.NoContent(http.StatusOK)
}