package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/interactor"

	"go.uber.org/zap"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/api/middleware"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/api/schema"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/log"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/input_port"
)

type UserHandler struct {
	UserUC input_port.IUserUseCase
}

func NewUserHandler(userUC input_port.IUserUseCase) *UserHandler {
	return &UserHandler{UserUC: userUC}
}

func (h *UserHandler) Create(c echo.Context) error {
	logger, _ := log.NewLogger()

	req := &schema.CreateUserReq{}
	if err := c.Bind(req); err != nil {
		logger.Error("Failed to bind request", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	res, err := h.UserUC.Create(
		input_port.UserCreate{
			StudentID: req.StudentID,
			IdmUniv:   req.IdmUniv,
			IdmBus:    req.IdmBus,
		})
	if err != nil {
		logger.Info("Failed to create user", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, schema.UserResFromEntity(res))
}

func (h *UserHandler) FindMe(c echo.Context) error {
	logger, _ := log.NewLogger()

	ctx := c.Request().Context()
	user, err := middleware.GetUserFromContext(ctx) // トークンからIDを取得
	if err != nil {
		logger.Error("Failed to get id from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	res, err := h.UserUC.FindByID(user, user.ID)
	if err != nil {
		logger.Info("Failed to find me", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, schema.UserResFromEntity(res))
}

func (h *UserHandler) FindByID(c echo.Context) error {
	logger, _ := log.NewLogger()

	var id string
	if err := echo.PathParamsBinder(c).MustString("user-id", &id).BindError(); err != nil {
		logger.Error("Failed to bind path param id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	user, err := middleware.GetUserFromContext(ctx) // トークンからIDを取得
	if err != nil {
		logger.Error("Failed to get id from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	res, err := h.UserUC.FindByIDByAdmin(user, id)
	if err != nil {
		logger.Info("Failed to find me", zap.Error(err))
		switch {
		case errors.Is(err, interactor.ErrKind.NotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, schema.UserResFromEntity(res))
}

func (h *UserHandler) Update(c echo.Context) error {
	logger, _ := log.NewLogger()

	req := &schema.UpdateUserReq{}
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

	res, err := h.UserUC.Update(
		user,
		input_port.UserUpdate{
			ID:    user.ID,
			StudentID: req.StudentID,
			IdmUniv:   req.IdmUniv,
			IdmBus:    req.IdmBus,
			UserType:  req.UserType,
		})
	if err != nil {
		logger.Info("Failed to update user", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, schema.UserResFromEntity(res))
}

func (h *UserHandler) Search(c echo.Context) error {
	logger, _ := log.NewLogger()

	ctx := c.Request().Context()
	user, err := middleware.GetUserFromContext(ctx) // トークンからIDを取得
	if err != nil {
		logger.Error("Failed to get id from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	req := &schema.UserSearchQueryReq{}
	if err := c.Bind(req); err != nil {
		logger.Error("Failed to bind request", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	list, total, err := h.UserUC.Search(
		user,
		req.Query,
		req.UserType,
		req.Skip,
		req.Limit,
	)
	if err != nil {
		logger.Info("Failed by invalid request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, schema.UsersResFromSearchResult(list, total))
}

func (h *UserHandler) Delete(c echo.Context) error {
	logger, _ := log.NewLogger()

	ctx := c.Request().Context()
	user, err := middleware.GetUserFromContext(ctx) // トークンからIDを取得
	if err != nil {
		logger.Error("Failed to get id from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var id string
	if err := echo.PathParamsBinder(c).MustString("user-id", &id).BindError(); err != nil {
		logger.Error("Failed to bind path param id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res, err := h.UserUC.Delete(user, id)
	if err != nil {
		logger.Info("Failed to delete", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, schema.UserResFromEntity(res))
}
