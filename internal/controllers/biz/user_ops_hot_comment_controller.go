package biz

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/utils"
)

// -------------------- 热点评论 --------------------

func (uc *UserOpsController) ListHotComments(c *gin.Context) {
	items, err := uc.svc.ListHotComments(c.Request.Context(), 200)
	if err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"items": items})
}

func (uc *UserOpsController) CreateHotComment(c *gin.Context) {
	var req struct {
		PostID   uint   `json:"post_id"`
		UserID   uint   `json:"user_id"`
		ParentID *uint  `json:"parent_id"`
		Content  string `json:"content"`
		Likes    *int64 `json:"likes"`
		Status   *int8  `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		return
	}
	if req.PostID == 0 || req.UserID == 0 || strings.TrimSpace(req.Content) == "" {
		utils.JSONError(c, http.StatusBadRequest, "post_id/user_id/content required")
		return
	}

	item, err := uc.svc.CreateHotComment(c.Request.Context(), req.PostID, req.UserID, req.ParentID, strings.TrimSpace(req.Content), req.Likes, req.Status)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.JSONOK(c, item)
}

func (uc *UserOpsController) UpdateHotComment(c *gin.Context) {
	id, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		Content *string `json:"content"`
		Likes   *int64  `json:"likes"`
		Status  *int8   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		return
	}

	updates := map[string]interface{}{}
	if req.Content != nil {
		updates["content"] = strings.TrimSpace(*req.Content)
	}
	if req.Likes != nil {
		updates["likes"] = *req.Likes
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if len(updates) == 0 {
		utils.JSONError(c, http.StatusBadRequest, "empty updates")
		return
	}

	if err := uc.svc.UpdateHotComment(c.Request.Context(), id, updates); err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}

func (uc *UserOpsController) DeleteHotComment(c *gin.Context) {
	id, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := uc.svc.DeleteHotComment(c.Request.Context(), id); err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}
