package controllers

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
	var items []models.WPostArticle
	if err := bc.db.Where("is_official = 1").Order("id DESC").Limit(200).Find(&items).Error; err != nil {
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
	if !bc.isUserType(req.UserID, "official") {
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
	if err := bc.db.Create(&item).Error; err != nil {
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
		if userID == 0 || !bc.isUserType(userID, "official") {
			utils.JSONError(c, http.StatusBadRequest, "user_id must be official account")
			return
		}
		req["user_id"] = userID
	}
	if err := bc.db.Model(&models.WPostArticle{}).Where("id = ?", id).Updates(req).Error; err != nil {
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
	if err := bc.db.Delete(&models.WPostArticle{}, id).Error; err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}

func (bc *BizConfigController) isUserType(userID uint, expected string) bool {
	var user models.WUser
	if err := bc.db.Select("id,user_type").Where("id = ? AND status = 1", userID).First(&user).Error; err != nil {
		return false
	}
	return strings.TrimSpace(user.UserType) == expected
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
