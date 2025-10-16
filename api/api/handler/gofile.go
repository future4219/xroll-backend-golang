package handler

import (
	"errors"
	"io"
	"net/http"
	"os"
	"syscall"

	"github.com/labstack/echo/v4"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/api/middleware"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/api/schema"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entconst"
	log "gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/log"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/input_port"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/interactor"
	"go.uber.org/zap"
)

type GofileHandler struct {
	GofileUC input_port.IGofileUseCase
}

func NewGofileHandler(gofileUC input_port.IGofileUseCase) *GofileHandler {
	return &GofileHandler{GofileUC: gofileUC}
}

func (g *GofileHandler) Create(c echo.Context) error {
	logger, _ := log.NewLogger()

	req := &schema.GofileCreateReq{}
	if err := c.Bind(req); err != nil {
		logger.Error("Failed to bind request", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	ctx := c.Request().Context()
	user, err := middleware.GetUserFromContext(ctx) // トークンからIDを取得
	if err != nil {
		logger.Error("Failed to get id from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	res, err := g.GofileUC.Create(
		user,
		input_port.GofileCreate{
			Name:        req.Name,
			GofileID:    req.GofileID,
			TagIDs:      req.TagIDs,
			UserID:      req.UserID,
			GofileToken: req.GofileToken,
		})
	if err != nil {
		logger.Info("Failed to create gofile", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, schema.GofileCreateResFromEntity(res))
}

func (g *GofileHandler) Update(c echo.Context) error {
	logger, _ := log.NewLogger()
	var id string
	if err := echo.PathParamsBinder(c).MustString("id", &id).BindError(); err != nil {
		logger.Info("Failed to bind path param id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	req := &schema.GofileUpdateReq{}
	if err := c.Bind(req); err != nil {
		logger.Error("Failed to bind request", zap.Error(err))

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	ctx := c.Request().Context()
	user, err := middleware.GetUserFromContext(ctx) // トークンからIDを取得
	if err != nil {
		logger.Error("Failed to get id from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	res, err := g.GofileUC.Update(
		user,
		input_port.GofileUpdate{
			ID:          id,
			Name:        req.Name,
			Description: req.Description,
			TagIDs:      req.TagIDs,
			IsShare:     req.IsShared,
		})
	if err != nil {
		logger.Info("Failed to update gofile", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())

		case errors.Is(err, interactor.ErrKind.Unauthorized):
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, schema.GofileVideoResFromEntity(res))
}

func (g *GofileHandler) FindByID(c echo.Context) error {
	logger, _ := log.NewLogger()

	var id string
	if err := echo.PathParamsBinder(c).MustString("id", &id).BindError(); err != nil {
		logger.Info("Failed to bind path param id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	user, err := middleware.GetUserFromContext(ctx) // トークンからIDを取得
	if err != nil {
		logger.Error("Failed to get id from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	res, hasLike, err := g.GofileUC.FindByID(user, id)
	if err != nil {
		logger.Error("Failed to find by id", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, schema.GofileVideoResWithLike{
		GofileVideoRes: schema.GofileVideoResFromEntity(res),
		HasLike:        hasLike,
	})
}

func (g *GofileHandler) FindByUserID(c echo.Context) error {
	logger, _ := log.NewLogger()

	ctx := c.Request().Context()
	user, err := middleware.GetUserFromContext(ctx) // トークンからIDを取得
	if err != nil {
		logger.Error("Failed to get id from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	res, err := g.GofileUC.FindByUserID(user)
	if err != nil {
		logger.Error("Failed to find by user id", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, schema.GofileVideoListFromEntity(res))
}

func (g *GofileHandler) FindByUserIDShared(c echo.Context) error {
	logger, _ := log.NewLogger()

	var targetUserId string
	if err := echo.PathParamsBinder(c).MustString("userId", &targetUserId).BindError(); err != nil {
		logger.Info("Failed to bind path param id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	user, err := middleware.GetUserFromContext(ctx) // トークンからIDを取得
	if err != nil {
		logger.Error("Failed to get id from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	res, err := g.GofileUC.FindByUserIDShared(user, targetUserId)
	if err != nil {
		logger.Error("Failed to find non-shared by user id", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, schema.GofileVideoListFromEntity(res))
}

func (g *GofileHandler) UpdateIsShareVideo(c echo.Context) error {
	logger, _ := log.NewLogger()

	ctx := c.Request().Context()
	user, err := middleware.GetUserFromContext(ctx) // トークンからIDを取得
	if err != nil {
		logger.Error("Failed to get id from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var req schema.GofileUpdateIsShareReq
	if err := c.Bind(&req); err != nil {
		logger.Error("Failed to bind request", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = g.GofileUC.UpdateIsShareVideo(user, req.VideoID, req.IsShared)
	if err != nil {
		logger.Error("Failed to update is share", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.NoContent(http.StatusOK)
}

func (g *GofileHandler) Delete(c echo.Context) error {
	logger, _ := log.NewLogger()

	ctx := c.Request().Context()
	user, err := middleware.GetUserFromContext(ctx) // トークンからIDを取得
	if err != nil {
		logger.Error("Failed to get id from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var id string
	if err := echo.PathParamsBinder(c).MustString("id", &id).BindError(); err != nil {
		logger.Info("Failed to bind path param id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = g.GofileUC.Delete(user, id)
	if err != nil {
		logger.Error("Failed to delete gofile video", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.NoContent(http.StatusOK)
}
func (g *GofileHandler) LikeVideo(c echo.Context) error {
	logger, _ := log.NewLogger()
	ctx := c.Request().Context()
	user, err := middleware.GetUserFromContext(ctx) // トークンからIDを取得
	if err != nil {
		logger.Error("Failed to get id from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var id string
	if err := echo.PathParamsBinder(c).MustString("id", &id).BindError(); err != nil {
		logger.Info("Failed to bind path param id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = g.GofileUC.LikeVideo(user, id)
	if err != nil {
		logger.Error("Failed to like video", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.NoContent(http.StatusOK)
}

func (g *GofileHandler) UnlikeVideo(c echo.Context) error {
	logger, _ := log.NewLogger()

	ctx := c.Request().Context()
	user, err := middleware.GetUserFromContext(ctx) // トークンからIDを取得
	if err != nil {
		logger.Error("Failed to get id from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var id string
	if err := echo.PathParamsBinder(c).MustString("id", &id).BindError(); err != nil {
		logger.Info("Failed to bind path param id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = g.GofileUC.UnlikeVideo(user, id)
	if err != nil {
		logger.Error("Failed to unlike video", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.NoContent(http.StatusOK)
}

func (g *GofileHandler) FindLikedVideos(c echo.Context) error {
	logger, _ := log.NewLogger()

	ctx := c.Request().Context()
	user, err := middleware.GetUserFromContext(ctx) // トークンからIDを取得
	if err != nil {
		logger.Error("Failed to get id from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	videos, err := g.GofileUC.FindLikedVideos(user)
	if err != nil {
		logger.Error("Failed to find liked videos", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, schema.GofileVideoListFromEntity(videos))
}

func (g *GofileHandler) Search(c echo.Context) error {
	logger, _ := log.NewLogger()

	ctx := c.Request().Context()
	user, err := middleware.GetUserFromContext(ctx) // トークンからIDを取得
	if err != nil {
		logger.Error("Failed to get id from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	req, err := schema.BindGofileVideoSearchReq(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	order, err := entconst.NewOrder(req.Order)
	if err != nil {
		logger.Info("Failed to new sort order", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	orderBy, err := entconst.NewGofileOrderBy(req.OrderBy)
	if err != nil {
		logger.Info("Failed to new gofile order by", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	query := input_port.GofileSearchQuery{
		Q:       req.Q,
		Skip:    req.Skip,
		Limit:   req.Limit,
		OrderBy: orderBy,
		Order:   order,
	}

	videos, err := g.GofileUC.Search(user, query)
	if err != nil {
		logger.Error("Failed to search videos", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, schema.GofileVideoListFromEntity(videos))
}

func (g *GofileHandler) CreateComment(c echo.Context) error {
	logger, _ := log.NewLogger()
	ctx := c.Request().Context()
	user, err := middleware.GetUserFromContext(ctx) // トークンからIDを取得
	if err != nil {
		logger.Error("Failed to get id from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	var videoId string
	if err := echo.PathParamsBinder(c).MustString("video-id", &videoId).BindError(); err != nil {
		logger.Info("Failed to bind path param id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var req schema.GofileCreateCommentReq
	if err := c.Bind(&req); err != nil {
		logger.Error("Failed to bind request", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	res, err := g.GofileUC.CreateComment(user, input_port.GofileVideoCommentCreate{
		VideoID: videoId,
		Comment: req.Comment,
	})
	if err != nil {
		logger.Error("Failed to create comment", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusCreated, schema.GofileVideoCommentResFromEntity(res))
}

func (g *GofileHandler) CreateFromTwimgURL(c echo.Context) error {
	logger, _ := log.NewLogger()
	ctx := c.Request().Context()
	user, err := middleware.GetUserFromContext(ctx) // トークンからIDを取得
	if err != nil {
		logger.Error("Failed to get id from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	var req schema.GofileCreateFromTwimgURLReq
	if err := c.Bind(&req); err != nil {
		logger.Error("Failed to bind request", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	res, err := g.GofileUC.CreateFromTwimgURL(user, req.TwimgURL)
	if err != nil {
		logger.Error("Failed to create from twimg url", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusCreated, schema.GofileCreateResFromEntity(res))
}

func (g *GofileHandler) ProxyGofileVideo(c echo.Context) error {
	logger, _ := log.NewLogger()

	ctx := c.Request().Context()
	user, err := middleware.GetUserFromContext(ctx) // トークンからIDを取得
	if err != nil {
		logger.Error("Failed to get id from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	gofileVideoId := c.QueryParam("id")
	if gofileVideoId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "IDが必要です (?id=...)")
	}

	gofile, _, err := g.GofileUC.FindByID(user, gofileVideoId)
	if err != nil {
		logger.Error("Failed to find by id", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	decodedURL := gofile.GofileDirectURL

	if gofile.IsShared == false && gofile.UserID != user.ID {
		return echo.NewHTTPError(http.StatusForbidden, "this video is not shared")
	}

	token := os.Getenv("GOFILE_API_KEY")
	if token == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "GOFILE_API_KEY が未設定です")
	}

	// クライアントメソッド踏襲（HEAD対応）
	upMethod := c.Request().Method
	if upMethod != http.MethodGet && upMethod != http.MethodHead {
		upMethod = http.MethodGet
	}

	reqUp, err := http.NewRequestWithContext(c.Request().Context(), upMethod, decodedURL, nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, "上流リクエスト生成に失敗しました: "+err.Error())
	}

	// Range / If-Range 転送
	if rng := c.Request().Header.Get("Range"); rng != "" {
		reqUp.Header.Set("Range", rng)
	}
	if ifr := c.Request().Header.Get("If-Range"); ifr != "" {
		reqUp.Header.Set("If-Range", ifr)
	}

	// まず最初のリクエストに付与
	reqUp.Header.Set("Authorization", "Bearer "+token)
	reqUp.Header.Set("User-Agent", "Mozilla/5.0")
	reqUp.Header.Set("Accept", "*/*")

	// ★ ここがキモ：リダイレクト時にもヘッダを引き継ぐ
	client := &http.Client{
		Timeout: 0,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// リダイレクト先にも同じヘッダを再設定
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("User-Agent", "Mozilla/5.0")
			req.Header.Set("Accept", "*/*")
			// 必要なら Range も引き継ぎ
			if rng := c.Request().Header.Get("Range"); rng != "" {
				req.Header.Set("Range", rng)
			}
			if ifr := c.Request().Header.Get("If-Range"); ifr != "" {
				req.Header.Set("If-Range", ifr)
			}
			return nil
		},
	}

	up, err := client.Do(reqUp)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, "Gofile取得に失敗しました: "+err.Error())
	}
	defer up.Body.Close()

	// デバッグしたい時は以下で中身確認（一時的に）
	// if up.StatusCode >= 400 {
	// 	b, _ := io.ReadAll(io.LimitReader(up.Body, 2048))
	// 	c.Logger().Warnf("upstream %d from %s: %s", up.StatusCode, up.Request.URL, string(b))
	// }

	resp := c.Response()
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Expose-Headers", "Content-Range,Content-Length,Accept-Ranges")

	hopByHop := map[string]struct{}{
		"Connection": {}, "Keep-Alive": {}, "Proxy-Authenticate": {}, "Proxy-Authorization": {},
		"TE": {}, "Trailers": {}, "Transfer-Encoding": {}, "Upgrade": {},
	}
	for k, vv := range up.Header {
		if _, skip := hopByHop[k]; skip {
			continue
		}
		for _, v := range vv {
			resp.Header().Add(k, v)
		}
	}

	if resp.Header().Get("Content-Type") == "" {
		resp.Header().Set("Content-Type", "video/mp4")
	}
	if resp.Header().Get("Content-Disposition") == "" {
		resp.Header().Set("Content-Disposition", "inline")
	}

	resp.WriteHeader(up.StatusCode)
	if upMethod == http.MethodHead {
		return nil
	}
	_, copyErr := io.Copy(resp, up.Body)
	if copyErr != nil && !errors.Is(copyErr, syscall.EPIPE) && !errors.Is(copyErr, syscall.ECONNRESET) {
		c.Logger().Warnf("stream copy error: %v", copyErr)
	}
	return nil
}
