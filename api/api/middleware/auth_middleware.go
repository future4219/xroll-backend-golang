package middleware

import (
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/api/schema"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/log"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/input_port"

	"go.uber.org/zap"

	"context"
	"errors"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entconst"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity"
)

var (
	ErrNoAuthorizationHeader   = errors.New("no authorization header passed")
	ErrNoStripeSignatureHeader = errors.New("no Stripe-Signature header passed")
)

type AuthMiddleware struct {
	userUC input_port.IUserUseCase
}

func NewAuthMiddleware(userUC input_port.IUserUseCase) *AuthMiddleware {
	return &AuthMiddleware{userUC}
}

// Authenticate
// tokenを取得して、認証するmiddlewareの例
func (m *AuthMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	logger, _ := log.NewLogger()

	return func(c echo.Context) error {
		// Get JWT Token From Header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, schema.TokenType+" ") {
			logger.Info("Failed to authenticate", zap.Error(ErrNoAuthorizationHeader))
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		token := strings.TrimPrefix(authHeader, schema.TokenType+" ")

		// Authenticate
		userID, err := m.userUC.Authenticate(token)
		if err != nil {
			logger.Info("Failed to authenticate", zap.Error(err))
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		// set user detail to context
		user, err := m.userUC.FindByID(entity.User{
			UserType: entconst.SystemAdmin.String(),
		}, userID)
		if err != nil {
			logger.Error("Failed to find me", zap.Error(err))
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		c = SetToContext(c, user)

		return next(c)
	}
}

// tokenを取得するが認証はしないmiddleware
func (m *AuthMiddleware) NotAuthenticateButToSetUserToContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get JWT Token From Header
		authHeader := c.Request().Header.Get("Authorization")
		token := strings.TrimPrefix(authHeader, schema.TokenType+" ")
		userID, _ := m.userUC.Authenticate(token)
		user, err := m.userUC.FindByID(entity.User{
			UserType: entconst.SystemAdmin.String(),
		}, userID)
		if err == nil {
			c = SetToContext(c, user)
		}
		return next(c)
	}
}

func (m *AuthMiddleware) AuthenticateForUpdatePassword(next echo.HandlerFunc) echo.HandlerFunc {
	logger, _ := log.NewLogger()

	return func(c echo.Context) error {
		// Get JWT Token From Header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, schema.TokenType+" ") {
			logger.Info("Failed to authenticate", zap.Error(ErrNoAuthorizationHeader))
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		token := strings.TrimPrefix(authHeader, schema.TokenType+" ")

		// Authenticate
		userID, err := m.userUC.AuthenticateForUpdatePassword(token)
		if err != nil {
			logger.Info("Failed to authenticate", zap.Error(err))
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		// set user detail to context
		user, err := m.userUC.FindByID(entity.User{
			UserType: entconst.NonMemberUser.String(),
		}, userID) //FindByIDのdecoratorの認証を通すために書いています。

		if err != nil {
			logger.Error("Failed to find me", zap.Error(err))
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		c = SetToContext(c, user)

		return next(c)
	}
}

func (m *AuthMiddleware) AuthenticateForUpdateEmail(next echo.HandlerFunc) echo.HandlerFunc {
	logger, _ := log.NewLogger()

	return func(c echo.Context) error {
		// Get JWT Token From Header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, schema.TokenType+" ") {
			logger.Info("Failed to authenticate", zap.Error(ErrNoAuthorizationHeader))
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		token := strings.TrimPrefix(authHeader, schema.TokenType+" ")

		// Authenticate
		userID, err := m.userUC.AuthenticateForUpdateEmail(token)
		if err != nil {
			logger.Info("Failed to authenticate", zap.Error(err))
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		// set user detail to context
		user, err := m.userUC.FindByID(entity.User{
			UserType: entconst.NonMemberUser.String(),
		}, userID)
		if err != nil {
			logger.Error("Failed to find me", zap.Error(err))
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		c = SetToContext(c, user)

		return next(c)
	}
}

func SetToContext(c echo.Context, user entity.User) echo.Context {
	ctx := c.Request().Context()
	ctx = SetUserToContext(ctx, user)
	c.SetRequest(c.Request().WithContext(ctx))
	return c
}

type ContextKey string

var (
	userKey ContextKey = "userKey"
)

func SetUserToContext(ctx context.Context, user entity.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func GetUserFromContext(ctx context.Context) (entity.User, error) {
	v := ctx.Value(userKey)
	user, ok := v.(entity.User)
	if !ok {
		return entity.User{}, errors.New("no user found in context")
	}
	return user, nil
}

func GetStripeWebhookInfo(c echo.Context) (header string, body []byte, err error) {
	req := c.Request()

	header = req.Header.Get("Stripe-Signature")
	if header == "" {
		return "", nil, ErrNoStripeSignatureHeader
	}

	const MaxBodyBytes = int64(65536)

	reqBody := http.MaxBytesReader(c.Response().Writer, req.Body, MaxBodyBytes)

	body, err = io.ReadAll(reqBody)
	if err != nil {
		return "", nil, err
	}

	return header, body, nil
}
