package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/seth16888/wxbusiness/internal/biz"
	"github.com/seth16888/wxbusiness/internal/model/r"
	"go.uber.org/zap"
)

type MPMemberHandler struct {
	Base
	uc  *biz.MPMemberUsecase
	log *zap.Logger
}

func NewMPMemberHandler(log *zap.Logger, uc *biz.MPMemberUsecase) *MPMemberHandler {
	return &MPMemberHandler{uc: uc, log: log}
}

// Query
func (h *MPMemberHandler) Query(ctx *gin.Context) {
	// 路径参数
	appId, err := h.GetPID(ctx)
	if err != nil {
		ctx.JSON(400, r.Error(400, err.Error()))
		return
	}

	c := ctx
	tags, err := h.uc.Query(c, appId)
	if err != nil {
		h.log.Error("query tag error", zap.Error(err))
		ctx.JSON(500, r.Error(500, "查询标签失败"))
		return
	}

	ctx.JSON(200, r.SuccessData(tags))
}

// GetMemberInfo
func (h *MPMemberHandler) GetMemberInfo(ctx *gin.Context) {
	// 路径参数
	appId, err := h.GetPID(ctx)
	if err != nil {
		ctx.JSON(400, r.Error(400, err.Error()))
		return
	}

	// query 参数
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(400, r.Error(400, "参数错误"))
		return
	}
	c := ctx
	member, err := h.uc.GetMemberInfo(c, appId, id)
	if err != nil {
		h.log.Error("get member info error", zap.Error(err))
		ctx.JSON(500, r.Error(500, "获取用户信息失败"))
		return
	}

	ctx.JSON(200, r.SuccessData(member))
}

// GetMemberTags
func (h *MPMemberHandler) GetMemberTags(ctx *gin.Context) {
}

// UpdateRemark
func (h *MPMemberHandler) UpdateRemark(ctx *gin.Context) {
	// 路径参数
	appId, err := h.GetPID(ctx)
	if err != nil {
		ctx.JSON(400, r.Error(400, err.Error()))
		return
	}
	type req struct {
		Remark string `json:"remark" binding:"required" msg:"备注不能为空"`
		Id     string `json:"id" binding:"required" msg:"id不能为空"`
		OpenId string `json:"openid" binding:"required" msg:"openid不能为空"`
	}
	var params req
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(400, r.Error(400, err.Error()))
		return
	}

	c := ctx
	if err := h.uc.UpdateRemark(c, appId, params.Id, params.OpenId,
		params.Remark); err != nil {
		h.log.Error("update remark error", zap.Error(err))
		ctx.JSON(500, r.Error(500, "更新备注失败"))
		return
	}

	ctx.JSON(200, r.Success())
}

// GetBlackList
func (h *MPMemberHandler) GetBlackList(ctx *gin.Context) {
	// 路径参数
	appId, err := h.GetPID(ctx)
	if err != nil {
		ctx.JSON(400, r.Error(400, err.Error()))
		return
	}

	c := ctx
	members, err := h.uc.GetBlackList(c, appId)
	if err != nil {
		h.log.Error("get black list error", zap.Error(err))
		ctx.JSON(500, r.Error(500, "获取黑名单失败"))
		return
	}

	ctx.JSON(200, r.SuccessData(members))
}

// BatchBlock 批量拉黑,一次最多20个
func (h *MPMemberHandler) BatchBlock(ctx *gin.Context) {
	// 路径参数
	appId, err := h.GetPID(ctx)
	if err != nil {
		ctx.JSON(400, r.Error(400, err.Error()))
		return
	}
	type req struct {
		OpenIds []string `json:"openids" binding:"required" msg:"openids不能为空"`
	}
	var params req
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(400, r.Error(400, err.Error()))
		return
	}
	// 最多20个
	if len(params.OpenIds) > 20 {
		ctx.JSON(400, r.Error(400, "一次最多20个"))
		return
	}

	c := ctx
	if err := h.uc.BatchBlock(c, appId, params.OpenIds); err != nil {
		h.log.Error("batch block error", zap.Error(err))
		ctx.JSON(500, r.Error(500, "批量拉黑失败"))
		return
	}

	ctx.JSON(200, r.Success())
}

// BatchUnblock 批量取消拉黑,一次最多20个
func (h *MPMemberHandler) BatchUnblock(ctx *gin.Context) {
	// 路径参数
	appId, err := h.GetPID(ctx)
	if err != nil {
		ctx.JSON(400, r.Error(400, err.Error()))
		return
	}
	type req struct {
		OpenIds []string `json:"openids" binding:"required" msg:"openids不能为空"`
	}
	var params req
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(400, r.Error(400, err.Error()))
		return
	}
	// 最多20个
	if len(params.OpenIds) > 20 {
		ctx.JSON(400, r.Error(400, "一次最多20个"))
		return
	}

	c := ctx
	if err := h.uc.BatchUnblock(c, appId, params.OpenIds); err != nil {
		h.log.Error("batch unblock error", zap.Error(err))
		ctx.JSON(500, r.Error(500, "批量取消拉黑失败"))
		return
	}
	ctx.JSON(200, r.Success())
}

// Pull 同步微信粉丝
func (h *MPMemberHandler) Pull(ctx *gin.Context) {
	// 路径参数
	_, err := h.GetPID(ctx)
	if err != nil {
		ctx.JSON(400, r.Error(400, err.Error()))
		return
	}

	ctx.JSON(200, r.Success())
}
