package router

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/seth16888/wxbusiness/internal/di"
	"github.com/seth16888/wxbusiness/internal/handler"
	"github.com/seth16888/wxbusiness/internal/middleware"
)

func InitRouter(deps *di.Container) http.Handler {
	gin.DisableBindValidation()
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(middleware.LoggingMiddleware(deps.Log))
	r.Use(gin.Recovery())
	r.Use(cors.Default())
	r.Use(requestid.New())
  r.Use(middleware.NewJwtTokenMW(deps.JWT))

	registerRoutes(r, deps)

	return r
}

func registerRoutes(r *gin.Engine, deps *di.Container) {
	r.GET("/ping", deps.HealthHandler.Ping)
	r.GET("/health", deps.HealthHandler.Health)

	v1 := r.Group("/v1")
	{
    mp := v1.Group("/mp")
    {
      portalCtr:= handler.NewPortalHandler(deps.Log, deps.Validator, deps.PortalUsecase)
      mp.GET("/:id/portal", portalCtr.Verify)
      mp.POST("/:id/portal", portalCtr.Portal)
    }

    platform := v1.Group("/platform")
    {
      platformAppHdl := handler.NewPlatformAppHandler(deps.Log, deps.Validator, deps.PlatformAppUsecase)
      platform.POST("/apps", platformAppHdl.Create)
    }
	}
}
