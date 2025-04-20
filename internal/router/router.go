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
		userCtr := handler.NewUserHandler(deps.Log, deps.Validator, deps.UserUsecase)
		v1.GET("/users/:userId/apps", userCtr.ListMPApps)
		// v1/auth
		auth := v1.Group("/auth")
		{
			authCtr := handler.NewAuthHandler(deps.Log, deps.Validator, deps.CoAuthClient)
			auth.GET("/captcha", authCtr.GetCaptcha)
			auth.POST("/login", authCtr.Login)
		}
		// v1/portal/:id
		portGrp := v1.Group("/portal")
		{
			portalCtr := handler.NewPortalHandler(deps.Log, deps.Validator, deps.PortalUsecase)
			portGrp.GET("/:id", portalCtr.Verify)
			portGrp.POST("/:id", portalCtr.Portal)
		}

		// v1/apps
		appGrp := v1.Group("/apps")
		{
			appCtr := handler.NewAppHandler(deps.Log, deps.Validator, deps.AppUsecase)
			appGrp.POST("", appCtr.Create)

			// v1/apps/:id
			appGrp := appGrp.Group("/:id", middleware.NewAppMiddleware(deps.AppUsecase))
			{
				appGrp.GET("", appCtr.GetById)
				// v1/apps/:id/menu
				menuGrp := appGrp.Group("/menu")
				{
					menuCtr := handler.NewMPMenuHandler(deps.Log, deps.MenuUsecase)
					menuGrp.POST("", menuCtr.Create)
				}
				// v1/apps/:id/tags
				tagGrp := appGrp.Group("/tags")
				{
					tagCtr := handler.NewMemberTagHandler(deps.Log, deps.MemberTagUsecase)
					tagGrp.GET("", tagCtr.Query)
					tagGrp.POST("", tagCtr.Create)
					tagGrp.PUT("/:tagId", tagCtr.Update)
					tagGrp.DELETE("/:tagId", tagCtr.Delete)
					tagGrp.POST("/pull", tagCtr.Pull)
				}
				// v1/apps/:id/members
				memberGrp := appGrp.Group("/members")
				{
					memberCtr := handler.NewMPMemberHandler(deps.Log, deps.MPMemberUsecase, deps.Validator)
					memberGrp.GET("", memberCtr.Query)
					memberGrp.GET("/info", memberCtr.GetMemberInfo)
					memberGrp.GET("/tags", memberCtr.GetMemberTags)
          memberGrp.POST("/tags", memberCtr.BatchTagging)
          memberGrp.DELETE("/tags", memberCtr.BatchUnTagging)
					memberGrp.POST("/info/remark", memberCtr.UpdateRemark)
					memberGrp.POST("/blacklist/list", memberCtr.GetBlackList)
					memberGrp.POST("/blacklist/block", memberCtr.BatchBlock)
					memberGrp.POST("/blacklist/unblock", memberCtr.BatchUnblock)
          memberGrp.POST("/blacklist/pull", memberCtr.PullBlackList)
					memberGrp.POST("/pull", memberCtr.Pull)
				}
				// v1/apps/:id/materials
				materialGrp := appGrp.Group("/materials")
				{
					maCtr := handler.NewMaterialHandler(deps.Log, deps.MaterialUsecase)
					materialGrp.POST("/temporary", maCtr.UploadTemporaryMedia)
					materialGrp.POST("/news_image", maCtr.UploadNewsImage)
					materialGrp.POST("/limit", maCtr.UploadMedia)
					materialGrp.POST("/list", maCtr.GetMaterialList)
				}
        qrcodeGrp := appGrp.Group("/qrcode")
        {
          qrcodeCtr := handler.NewQRCodeHandler(deps.Log, deps.MpQRCodeUsecase, deps.Validator)
          qrcodeGrp.POST("/temporary", qrcodeCtr.CreateTemporary)
          qrcodeGrp.POST("/limit", qrcodeCtr.CreateLimit)
          qrcodeGrp.GET("/url", qrcodeCtr.GetURL)
        }
			}
		}
	}
}
