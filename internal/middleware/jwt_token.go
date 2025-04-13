package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/seth16888/wxbusiness/pkg/jwt"
)

func NewJwtTokenMW(jwt *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// authoriz := c.Request.Header.Get("Authorization")
		// if authoriz == "" {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, r.Error[any](401, "Unauthorized"))
		// 	return
		// }
		// // Bearer
		// const prefix = "Bearer "
		// if !strings.HasPrefix(authoriz, prefix) {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, r.Error[any](401, "Authorization must be Bearer"))
		// 	return
		// }

		// token := strings.TrimPrefix(authoriz, prefix)

		// claims, err := jwt.ParseToken(token)
		// if err != nil {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, r.Error[any](401, err.Error()))
		// 	return
		// }

		// uid := claims.UserID
		// if uid <= 0 {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, r.Error[any](401, "uid is empty"))
		// 	return
		// }
		c.Set("UID", "10001")

		// 当前用户信息
		// user := &model.CurrentUser{
		// 	UserID:      uid,
		// 	Username:    claims.Subject,
		// 	DeptID:      claims.DepartmentID,
		// 	DataScope:   claims.DataScope,
		// 	Authorities: claims.Authorities,
		// }
		// c.Set("CURRENT_USER", user)

		c.Next()
	}
}
