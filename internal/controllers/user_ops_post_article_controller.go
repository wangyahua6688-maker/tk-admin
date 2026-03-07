package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/models"
	"go-admin-full/internal/utils"
)

// -------------------- 发帖管理 --------------------

func (uc *UserOpsController) ListPostArticles(c *gin.Context) {
	var items []models.WPostArticle
	if err := uc.db.Where("is_official = 0").Order("id DESC").Limit(200).Find(&items).Error; err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"items": items})
}

func (uc *UserOpsController) CreatePostArticle(c *gin.Context) {
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
	if !uc.isUserType(req.UserID, "robot") {
		utils.JSONError(c, http.StatusBadRequest, "user_id must be robot account")
		return
	}

	item := models.WPostArticle{
		UserID:     req.UserID,
		Title:      strings.TrimSpace(req.Title),
		CoverImage: strings.TrimSpace(req.CoverImage),
		Content:    req.Content,
		IsOfficial: 0,
		Status:     1,
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

func (uc *UserOpsController) UpdatePostArticle(c *gin.Context) {
	id, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		UserID     *uint   `json:"user_id"`
		Title      *string `json:"title"`
		CoverImage *string `json:"cover_image"`
		Content    *string `json:"content"`
		Status     *int8   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		return
	}
	updates := map[string]interface{}{}
	if req.UserID != nil {
		if *req.UserID == 0 || !uc.isUserType(*req.UserID, "robot") {
			utils.JSONError(c, http.StatusBadRequest, "user_id must be robot account")
			return
		}
		updates["user_id"] = *req.UserID
	}
	if req.Title != nil {
		updates["title"] = strings.TrimSpace(*req.Title)
	}
	if req.CoverImage != nil {
		updates["cover_image"] = strings.TrimSpace(*req.CoverImage)
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	updates["is_official"] = int8(0)
	if len(updates) == 0 {
		utils.JSONError(c, http.StatusBadRequest, "empty updates")
		return
	}

	if err := uc.db.Model(&models.WPostArticle{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}

func (uc *UserOpsController) DeletePostArticle(c *gin.Context) {
	id, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := uc.db.Delete(&models.WPostArticle{}, id).Error; err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}

func (uc *UserOpsController) isUserType(userID uint, expected string) bool {
	return uc.isUserTypes(userID, expected)
}

func (uc *UserOpsController) isUserTypes(userID uint, expectedTypes ...string) bool {
	var u models.WUser
	if err := uc.db.Select("id,user_type").Where("id = ? AND status = 1", userID).First(&u).Error; err != nil {
		return false
	}
	current := strings.TrimSpace(u.UserType)
	for _, t := range expectedTypes {
		if current == t {
			return true
		}
	}
	return false
}
