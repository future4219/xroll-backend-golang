package proxy

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// PostCodeJPTokenが空の時用のmockハンドラー
func NewPostCodeJPProxyMockHandler() PostCodeJPProxy {
	return &PostCodeJPProxyHandlerMock{}
}

type PostCodeJPProxyHandlerMock struct{}

func (h *PostCodeJPProxyHandlerMock) SearchAboutPostCode(c echo.Context) error {
	return c.JSONBlob(http.StatusOK, nil)
}
