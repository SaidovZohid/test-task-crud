package api

import (
	v1 "github.com/SaidovZohid/test-task-crud/api/v1"
	"github.com/SaidovZohid/test-task-crud/config"
	"github.com/SaidovZohid/test-task-crud/pkg/logger"
	"github.com/SaidovZohid/test-task-crud/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RoutetOptions struct {
	Cfg      *config.Config
	Storage  storage.StorageI
	InMemory storage.InMemoryStorageI
	Log      logger.Logger
}

// New @title       Swagger for blogs and news test task
// @version         2.0
// @description     This is a documentation for new's and blog's apis.
// @BasePath  		/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @Security ApiKeyAuth
func New(opt *RoutetOptions) *gin.Engine {
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "*")
	router.Use(cors.New(corsConfig))

	handlerV1 := v1.New(&v1.HandlerV1Options{
		Cfg:      opt.Cfg,
		Storage:  opt.Storage,
		InMemory: opt.InMemory,
		Log:      opt.Log,
	})

	apiV1 := router.Group("/v1")
	{
		apiV1.POST("/auth/signup", handlerV1.SignUp)
		apiV1.POST("/auth/verify", handlerV1.Verify)
		apiV1.POST("/auth/login", handlerV1.Login)

		apiV1.POST("/blogs", handlerV1.AuthMiddleWare, handlerV1.CreateBlog)
		apiV1.GET("/blogs/:id", handlerV1.GetBlogByID)
		apiV1.GET("/blogs", handlerV1.GetAllPosts)
		apiV1.DELETE("/blogs/:id", handlerV1.AuthMiddleWare, handlerV1.DeleteBlog)
		apiV1.PUT("/blogs/:id", handlerV1.AuthMiddleWare, handlerV1.UpdateBlog)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
