package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/seth16888/wxbusiness/internal/biz"
	"github.com/seth16888/wxbusiness/internal/model/r"
)

// AppMiddleware 获取app信息
func NewAppMiddleware(uc *biz.AppUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
    id:= c.Param("id")
    if id == "" {
      c.AbortWithStatusJSON(400, r.Error(400, "app id is required"))
      return
    }

		app, err := uc.GetById(c, id)
		if err != nil {
			c.AbortWithStatusJSON(404, r.Error(404, "app not found"))
			return
		}
		c.Set("APP", app)
    c.Set("MP_ID", app.MpId)
		c.Next()
	}
}
