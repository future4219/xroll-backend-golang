package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/api/schema"
	log "gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/log"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/input_port"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/interactor"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/authentication"
)

type AuthHandler struct {
	UserUC input_port.IUserUseCase
}

func NewAuthHandler(userUC input_port.IUserUseCase) *AuthHandler {
	return &AuthHandler{UserUC: userUC}
}

// Login POST /auth/access-token
func (h *AuthHandler) Login(c echo.Context) error {
	logger, _ := log.NewLogger()

	req := &schema.LoginReq{}
	if err := c.Bind(req); err != nil {
		logger.Error("Failed to bind request", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	user, token, err := h.UserUC.Login(req.LoginID, req.Password)
	if err != nil {
		logger.Info("Failed to login", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		case errors.Is(err, authentication.ErrWrongPassword):
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	loginUser := &schema.LoginResUser{
		ID:  user.ID,
	}

	return c.JSON(http.StatusOK, &schema.LoginRes{
		AccessToken: token,
		TokenType:   schema.TokenType,
		User:        *loginUser,
	})
}

// ResetPassword POST /auth/reset-password
func (h *AuthHandler) ResetPassword(c echo.Context) error {
	logger, _ := log.NewLogger()

	req := &schema.ResetPasswordReq{}
	if err := c.Bind(req); err != nil {
		logger.Error("Failed to bind request", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err := h.UserUC.SendResetPasswordMail(req.Email)
	if err != nil {
		logger.Info("Failed to reset password", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.NoContent(http.StatusOK)
}
