package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/models"
	"go-admin-full/internal/utils"
)

// -------------------- 热点评论 --------------------

func (uc *UserOpsController) ListHotComments(c *gin.Context) {
	var items []models.WComment
	if err := uc.db.Where("status = 1 AND post_id > 0").Order("likes DESC, id DESC").Limit(200).Find(&items).Error; err != nil {
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

	item := models.WComment{
		PostID:   req.PostID,
		UserID:   req.UserID,
		ParentID: 0,
		Content:  strings.TrimSpace(req.Content),
		Likes:    0,
		Status:   1,
	}
	if req.ParentID != nil {
		item.ParentID = *req.ParentID
	}
	if req.Likes != nil {
		item.Likes = *req.Likes
	}
	if req.Status != nil {
		item.Status = *req.Status
	}

	if err := uc.db.Create(&item).Error; err != nil {
		utils.JSONError(c, 500, err.Error())
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

	if err := uc.db.Model(&models.WComment{}).Where("id = ? AND post_id > 0", id).Updates(updates).Error; err != nil {
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
	if err := uc.db.Where("post_id > 0").Delete(&models.WComment{}, id).Error; err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}
