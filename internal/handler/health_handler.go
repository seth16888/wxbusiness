package handler

import "github.com/gin-gonic/gin"

type HealthHandler struct {
}

func (h *HealthHandler) Health(c *gin.Context) {
	c.String(200, "ok")
}

func (h *HealthHandler) Ping(c *gin.Context) {
	c.String(200, "pong")
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}
