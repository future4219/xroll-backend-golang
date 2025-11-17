package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/api/middleware"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/api/schema"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entconst"
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

func (h *AuthHandler) Boot(c echo.Context) error {
	return nil // for compatibility
	logger, _ := log.NewLogger()

	ctx := c.Request().Context()
	user, _ := middleware.GetUserFromContext(ctx) // トークンからIDを取得

	user, token, err := h.UserUC.Boot(user)
	if err != nil {
		logger.Info("Failed to boot", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, &schema.BootRes{
		AccessToken: token,
		TokenType:   schema.TokenType,
		UserID:      user.ID,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	logger, _ := log.NewLogger()

	req := &schema.LoginReq{}
	if err := c.Bind(req); err != nil {
		logger.Error("Failed to bind request", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	user, token, err := h.UserUC.Login(req.Email, req.Password)
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
		ID: user.ID,
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

func (h *AuthHandler) CreateByMe(c echo.Context) error {
	logger, _ := log.NewLogger()

	var req schema.CreateByMeReq
	if err := c.Bind(&req); err != nil {
		logger.Error("Failed to bind request body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := h.UserUC.CreateByMe(input_port.CreateByMe{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		logger.Info("Failed to create user", zap.Error(err))
		var validationError *entconst.ValidationError
		switch {
		case errors.Is(err, interactor.ErrKind.Conflict):
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		case errors.As(err, &validationError):
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *AuthHandler) VerifyEmail(c echo.Context) error {
	logger, _ := log.NewLogger()

	var req schema.VerifyEmailReq
	if err := c.Bind(&req); err != nil {
		logger.Error("Failed to bind request body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	user, err := middleware.GetUserFromContext(ctx) // トークンからIDを取得
	if err != nil {
		logger.Error("Failed to get id from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	user, token, err := h.UserUC.VerifyEmail(user, input_port.VerifyEmail{
		Email:              req.Email,
		AuthenticationCode: req.AuthenticationCode,
	})
	if err != nil {
		logger.Info("Failed to verify email", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		case errors.Is(err, interactor.ErrKind.BadRequest):
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, schema.VerifyEmailResFromEntity(user, token))
}
