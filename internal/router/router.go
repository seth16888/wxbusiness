package router

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/seth16888/wxbusiness/internal/handler"
	"github.com/seth16888/wxbusiness/internal/middleware"
	"go.uber.org/zap"
)

func InitRouter(log *zap.Logger) http.Handler {
  gin.DisableBindValidation()
  gin.SetMode(gin.ReleaseMode)

  r := gin.New()

  r.Use(middleware.LoggingMiddleware(log))
	r.Use(gin.Recovery())
	r.Use(cors.Default())
	r.Use(requestid.New())

	registerRoutes(r)

  return r
}

func registerRoutes(r *gin.Engine) {
	r.GET("/ping", handler.NewHealthHandler().Ping)
	r.GET("/health", handler.NewHealthHandler().Health)
}
