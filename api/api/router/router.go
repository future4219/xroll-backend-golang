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
	videoUC input_port.IVideoUseCase,
	twitterUC input_port.ITwitterUseCase,
	gofileUC input_port.IGofileUseCase,
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
	videoHandler := handler.NewVideoHandler(videoUC)
	twitterHandler := handler.NewTwitterHandler(twitterUC)
	gofileHandler := handler.NewGofileHandler(gofileUC)

	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	api := e.Group("/api")
	api.POST("/auth/access-token", authHandler.Login)
	api.POST("/auth/reset-password", authHandler.ResetPassword)
	api.POST("/auth/create-by-me", authHandler.CreateByMe)

	// auth
	// 認可の例
	auth := api.Group("", apiMiddleware.NewAuthMiddleware(userUC).Authenticate)
	authIfPossible := api.Group("", apiMiddleware.NewAuthMiddleware(userUC).AuthenticateIfPossible)
	authCookieOrHeader := api.Group("", apiMiddleware.NewAuthMiddleware(userUC).AuthenticateCookieOrHeader)
	notAuth := api.Group("")

	auth.POST("/auth/verify-email", authHandler.VerifyEmail)

	// auth
	authIfPossible.GET("/auth/boot", authHandler.Boot)

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

	// video
	video := notAuth.Group("/videos")
	video.GET("/search", videoHandler.Search)
	video.POST("/create-bulk", videoHandler.CreateBulk)
	video.GET("/multiple", videoHandler.FindByIDs)
	video.GET("/:videoId", videoHandler.FindByID)
	video.POST("/like/:videoId", videoHandler.Like)
	video.POST("/comment/:videoId", videoHandler.Comment)

	// twitter
	twitter := notAuth.Group("/twitter")
	twitter.GET("/get-video-url", twitterHandler.GetVideoByURL)

	//gofile
	gofile := auth.Group("/gofile")
	gofile.POST("/create", gofileHandler.Create)
	gofile.PATCH("/update/:id", gofileHandler.Update)
	gofile.GET("/video/:id", gofileHandler.FindByID)
	gofile.GET("/:userId", gofileHandler.FindByUserID)
	gofile.GET("/search", gofileHandler.Search)
	gofile.GET("/:userId/shared", gofileHandler.FindByUserIDShared)
	gofile.PATCH("/update-is-shared", gofileHandler.UpdateIsShareVideo)
	gofile.DELETE("/delete/:id", gofileHandler.Delete)
	gofile.POST("/like/:id", gofileHandler.LikeVideo)
	gofile.POST("/unlike/:id", gofileHandler.UnlikeVideo)
	gofile.GET("/liked-videos", gofileHandler.FindLikedVideos)

	// authCookieOrHeader
	// Cookie or Header どちらでも認証可能
	gofileAuthCookieOrHeader := authCookieOrHeader.Group("/gofile")
	gofileAuthCookieOrHeader.GET("/proxy", gofileHandler.ProxyGofileVideo)

	return e
}
