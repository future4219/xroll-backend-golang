package proxy

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

	"go.uber.org/zap"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/config"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/log"
)

// NewPostCodeJPProxyHandler
// postcode-jp.comのAPIをプロキシするハンドラの例
func NewPostCodeJPProxyHandler() PostCodeJPProxy {
	return &PostCodeJPProxyHandler{}
}

type PostCodeJPProxy interface {
	SearchAboutPostCode(c echo.Context) error
}

type PostCodeJPProxyHandler struct{}

const postCodeJPEndpoint string = "https://apis.postcode-jp.com/api/v5/"

func (h *PostCodeJPProxyHandler) SearchAboutPostCode(c echo.Context) error {
	logger, _ := log.NewLogger()

	var postCode string
	if err := echo.PathParamsBinder(c).MustString("post-code", &postCode).BindError(); err != nil {
		logger.Error("Failed to bind path param post-code", zap.Error(err))
		return echo.ErrBadRequest
	}

	request, err := http.NewRequest("GET", postCodeJPEndpoint+"postcodes/"+postCode, nil)
	if err != nil {
		logger.Error("Failed to create internal request", zap.Error(err))
		return echo.ErrInternalServerError
	}

	request.Header.Add("apikey", config.PostCodeJPToken())
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.Error("Failed to request", zap.Error(err))
		return echo.ErrInternalServerError
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("Failed to close request body", zap.Error(err))
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read request body", zap.Error(err))
		return echo.ErrInternalServerError
	}
	return c.JSONBlob(http.StatusOK, body)
}
