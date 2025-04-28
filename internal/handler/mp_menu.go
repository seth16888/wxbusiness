package handler

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/seth16888/wxbusiness/internal/biz"
	"github.com/seth16888/wxbusiness/internal/model/r"
	"github.com/seth16888/wxcommon/mp"
)

type MPMenuHandler struct {
  Base
  log *zap.Logger
  menuUc *biz.MPMenuUsecase
}

func NewMPMenuHandler(log *zap.Logger, menuUc *biz.MPMenuUsecase) *MPMenuHandler {
  return &MPMenuHandler{
    log: log,
    menuUc: menuUc,
  }
}

func (h *MPMenuHandler) Create(c *gin.Context) {
  pId,err:= h.GetPID(c)
  if err != nil {
    c.JSON(400, r.Error(400, err.Error()))
    return
  }

  var params mp.CreateMenuReq
  if err := c.ShouldBindJSON(&params); err != nil {
    c.JSON(400, r.Error(400, err.Error()))
    return
  }

  ctx, cancel := context.WithTimeout(c, time.Second*10)
  defer cancel()
	rt := h.menuUc.Create(ctx, pId, &params)

	c.JSON(rt.StatusCode(), rt)
}

// GetMenuInfo
func (h *MPMenuHandler) GetMenuInfo(c *gin.Context) {
  pId,err:= h.GetPID(c)
  if err != nil {
    c.JSON(400, r.Error(400, err.Error()))
    return
  }

  ctx, cancel := context.WithTimeout(c, time.Second*10)
  defer cancel()
	rt := h.menuUc.GetMenuInfo(ctx, pId)

	c.JSON(rt.StatusCode(), rt)
}

// Delete
func (h *MPMenuHandler) Delete(c *gin.Context) {
  pId,err:= h.GetPID(c)
  if err != nil {
    c.JSON(400, r.Error(400, err.Error()))
    return
  }

  ctx, cancel := context.WithTimeout(c, time.Second*10)
  defer cancel()
	rt := h.menuUc.Delete(ctx, pId)

	c.JSON(rt.StatusCode(), rt)
}

// CreateConditional
func (h *MPMenuHandler) CreateConditional(c *gin.Context) {

}

// DeleteConditional
func (h *MPMenuHandler) DeleteConditional(c *gin.Context) {

}

// Pull
func (h *MPMenuHandler) Pull(c *gin.Context) {
  pId,err:= h.GetPID(c)
  if err != nil {
    c.JSON(400, r.Error(400, err.Error()))
    return
  }

  ctx, cancel := context.WithTimeout(c, time.Second*10)
  defer cancel()
  rt := h.menuUc.Pull(ctx, pId)

  c.JSON(rt.StatusCode(), rt)
}
