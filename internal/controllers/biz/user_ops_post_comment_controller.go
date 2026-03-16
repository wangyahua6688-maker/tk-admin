package biz

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/utils"
)

// -------------------- 帖子评论管理（按帖子维度） --------------------

func (uc *UserOpsController) ListPostComments(c *gin.Context) {
	postID, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid post id")
		return
	}
	items, err := uc.svc.ListPostComments(c.Request.Context(), postID)
	if err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"items": items})
}

func (uc *UserOpsController) CreatePostComment(c *gin.Context) {
	postID, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid post id")
		return
	}

	var req struct {
		UserID   uint   `json:"user_id"`
		ParentID uint   `json:"parent_id"`
		Content  string `json:"content"`
		Status   *int8  `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		return
	}
	if req.UserID == 0 || strings.TrimSpace(req.Content) == "" {
		utils.JSONError(c, http.StatusBadRequest, "user_id/content required")
		return
	}
	item, err := uc.svc.CreatePostComment(c.Request.Context(), postID, req.UserID, req.ParentID, strings.TrimSpace(req.Content), req.Status)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.JSONOK(c, item)
}
