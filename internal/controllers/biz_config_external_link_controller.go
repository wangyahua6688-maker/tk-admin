package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/models"
	"go-admin-full/internal/utils"
)

// -------------------- 外链 --------------------

func (bc *BizConfigController) ListExternalLinks(c *gin.Context) {
	var items []models.WExternalLink
	if err := bc.db.Order("position ASC, sort ASC, id DESC").Limit(200).Find(&items).Error; err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"items": items})
}

func (bc *BizConfigController) CreateExternalLink(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		URL      string `json:"url"`
		Position string `json:"position"`
		IconURL  string `json:"icon_url"`
		GroupKey string `json:"group_key"`
		Status   *int8  `json:"status"`
		Sort     *int   `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		return
	}
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.URL) == "" || strings.TrimSpace(req.Position) == "" {
		utils.JSONError(c, http.StatusBadRequest, "name/url/position required")
		return
	}
	item := models.WExternalLink{
		Name:     strings.TrimSpace(req.Name),
		URL:      strings.TrimSpace(req.URL),
		Position: strings.TrimSpace(req.Position),
		IconURL:  strings.TrimSpace(req.IconURL),
		GroupKey: strings.TrimSpace(req.GroupKey),
		Status:   1,
		Sort:     0,
	}
	if req.Status != nil {
		item.Status = *req.Status
	}
	if req.Sort != nil {
		item.Sort = *req.Sort
	}
	if err := bc.db.Create(&item).Error; err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, item)
}

func (bc *BizConfigController) UpdateExternalLink(c *gin.Context) {
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
	if len(req) == 0 {
		utils.JSONError(c, http.StatusBadRequest, "empty updates")
		return
	}
	if err := bc.db.Model(&models.WExternalLink{}).Where("id = ?", id).Updates(req).Error; err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}

func (bc *BizConfigController) DeleteExternalLink(c *gin.Context) {
	id, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := bc.db.Delete(&models.WExternalLink{}, id).Error; err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}
