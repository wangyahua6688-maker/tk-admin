package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/models"
	"go-admin-full/internal/utils"
)

// -------------------- 客户端用户 --------------------

func (uc *UserOpsController) ListClientUsers(c *gin.Context) {
	var items []models.WUser
	query := uc.db.Order("id DESC").Limit(300)
	if userType := strings.TrimSpace(c.Query("user_type")); userType != "" {
		query = query.Where("user_type = ?", userType)
	}

	if err := query.Find(&items).Error; err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"items": items})
}

func (uc *UserOpsController) CreateClientUser(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
		UserType string `json:"user_type"`
		Status   *int8  `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		return
	}
	if strings.TrimSpace(req.Username) == "" {
		utils.JSONError(c, http.StatusBadRequest, "username required")
		return
	}

	userType := normalizeUserType(req.UserType)

	item := models.WUser{
		Username: strings.TrimSpace(req.Username),
		Nickname: strings.TrimSpace(req.Nickname),
		Avatar:   strings.TrimSpace(req.Avatar),
		UserType: userType,
		Status:   1,
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

func (uc *UserOpsController) UpdateClientUser(c *gin.Context) {
	id, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		Nickname *string `json:"nickname"`
		Avatar   *string `json:"avatar"`
		UserType *string `json:"user_type"`
		Status   *int8   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		return
	}

	updates := map[string]interface{}{}
	if req.Nickname != nil {
		updates["nickname"] = strings.TrimSpace(*req.Nickname)
	}
	if req.Avatar != nil {
		updates["avatar"] = strings.TrimSpace(*req.Avatar)
	}
	if req.UserType != nil {
		updates["user_type"] = normalizeUserType(*req.UserType)
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if len(updates) == 0 {
		utils.JSONError(c, http.StatusBadRequest, "empty updates")
		return
	}

	if err := uc.db.Model(&models.WUser{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}

func (uc *UserOpsController) DeleteClientUser(c *gin.Context) {
	id, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := uc.db.Delete(&models.WUser{}, id).Error; err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}

func normalizeUserType(v string) string {
	switch strings.TrimSpace(v) {
	case "official", "robot", "natural":
		return strings.TrimSpace(v)
	default:
		return "natural"
	}
}
