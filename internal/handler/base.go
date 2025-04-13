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
