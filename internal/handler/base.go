package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Base 基础控制器
type Base struct{}

// GetUserId 获取登录用户的ID
func (b *Base) GetUserId(c *gin.Context) (string, error) {
	userId, exists := c.Get("UID")
	if !exists {
		return "", fmt.Errorf("user not logged in")
	}
  id, ok := userId.(string)
  if !ok {
    return "", fmt.Errorf("user id is not string")
  }
  return id, nil
}

// GetPID 获取应用ID
//
// c.Param("id") 路径参数，平台应用ID
func (b *Base) GetPID(c *gin.Context) (string, error) {
  appID := c.Param("id")
  if len(appID) == 0 {
    return "", fmt.Errorf("invalid appid")
  }
  return appID, nil
}
