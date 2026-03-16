package biz

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/models"
	"go-admin-full/internal/utils"
)

// -------------------- 官方发帖 --------------------

func (bc *BizConfigController) ListOfficialPosts(c *gin.Context) {
	items, err := bc.svc.ListOfficialPosts(c.Request.Context(), 200)
	if err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"items": items})
}

func (bc *BizConfigController) CreateOfficialPost(c *gin.Context) {
	var req struct {
		UserID     uint   `json:"user_id"`
		Title      string `json:"title"`
		CoverImage string `json:"cover_image"`
		Content    string `json:"content"`
		Status     *int8  `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		return
	}
	if strings.TrimSpace(req.Title) == "" {
		utils.JSONError(c, http.StatusBadRequest, "title required")
		return
	}
	if req.UserID == 0 {
		utils.JSONError(c, http.StatusBadRequest, "user_id required")
		return
	}
	if !bc.svc.IsUserTypes(c.Request.Context(), req.UserID, "official") {
		utils.JSONError(c, http.StatusBadRequest, "user_id must be official account")
		return
	}

	item := models.WPostArticle{
		UserID:     req.UserID,
		Title:      strings.TrimSpace(req.Title),
		CoverImage: strings.TrimSpace(req.CoverImage),
		Content:    req.Content,
		IsOfficial: 1,
		Status:     1,
	}
	if req.Status != nil {
		item.Status = *req.Status
	}
	if err := bc.svc.CreateOfficialPost(c.Request.Context(), &item); err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, item)
}

func (bc *BizConfigController) UpdateOfficialPost(c *gin.Context) {
	id, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		return
	}
	req["is_official"] = int8(1)
	if rawUserID, ok := req["user_id"]; ok {
		userID := toUint(rawUserID)
		if userID == 0 || !bc.svc.IsUserTypes(c.Request.Context(), userID, "official") {
			utils.JSONError(c, http.StatusBadRequest, "user_id must be official account")
			return
		}
		req["user_id"] = userID
	}
	if err := bc.svc.UpdateOfficialPost(c.Request.Context(), id, req); err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}

func (bc *BizConfigController) DeleteOfficialPost(c *gin.Context) {
	id, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := bc.svc.DeleteOfficialPost(c.Request.Context(), id); err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}

func toUint(raw interface{}) uint {
	switch v := raw.(type) {
	case float64:
		if v > 0 {
			return uint(v)
		}
	case int:
		if v > 0 {
			return uint(v)
		}
	case int64:
		if v > 0 {
			return uint(v)
		}
	case string:
		v = strings.TrimSpace(v)
		if v == "" {
			return 0
		}
		var n uint
		_, _ = fmt.Sscanf(v, "%d", &n)
		return n
	}
	return 0
}
