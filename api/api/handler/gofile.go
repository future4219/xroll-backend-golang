package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"syscall"

	"github.com/labstack/echo/v4"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/api/schema"
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

	res, err := g.GofileUC.Create(
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

func (g *GofileHandler) FindByUserID(c echo.Context) error {
	logger, _ := log.NewLogger()

	var userID string
	if err := echo.PathParamsBinder(c).MustString("userId", &userID).BindError(); err != nil {
		logger.Info("Failed to bind path param id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res, err := g.GofileUC.FindByUserID(userID)
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

func (g *GofileHandler) ProxyGofileVideo(c echo.Context) error {
	rawURL := c.QueryParam("url")
	if rawURL == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "URLパラメータが必要です (?url=...)")
	}

	decodedURL := rawURL
	if strings.Contains(rawURL, "%") {
		u, err := url.QueryUnescape(rawURL)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "URLのデコードに失敗しました")
		}
		decodedURL = u
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
