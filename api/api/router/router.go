package router

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4/middleware"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/api/handler"
	apiMiddleware "gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/api/middleware"
	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/input_port"
)

func NewServer(
	userUC input_port.IUserUseCase,
	fileUC input_port.IFileUseCase,
	isLogging bool,
) *echo.Echo {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch},
	}))

	if isLogging {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())

	authHandler := handler.NewAuthHandler(userUC)
	userHandler := handler.NewUserHandler(userUC)
	fileHandler := handler.NewFileHandler(fileUC)

	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	api := e.Group("/api")
	api.POST("/auth/access-token", authHandler.Login)
	api.POST("/auth/reset-password", authHandler.ResetPassword)

	// auth
	// 認可の例
	auth := api.Group("", apiMiddleware.NewAuthMiddleware(userUC).Authenticate)
	// user
	user := auth.Group("/users")
	user.GET("", userHandler.Search)
	user.GET("/me", userHandler.FindMe)
	user.GET("/:user-id", userHandler.FindByID)
	user.PATCH("/:user-id", userHandler.Update)
	user.DELETE("/:user-id", userHandler.Delete)

	// file
	file := auth.Group("/files")
	file.POST("/upload", fileHandler.IssuePreSignedURLForPut)
	file.POST("/upload/video", fileHandler.IssuePreSignedURLForPutVideo)
	file.GET("/video/:fileId/:fileName", fileHandler.IssuePresignedURLForGetVideo)

	return e
}
