package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/models"
	"go-admin-full/internal/utils"
)

// -------------------- Lottery Category --------------------

func (bc *BizConfigController) ListLotteryCategories(c *gin.Context) {
	var items []models.WLotteryCategory
	query := bc.db.Order("sort ASC, id ASC")

	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("category_key LIKE ? OR name LIKE ? OR search_keywords LIKE ?", like, like, like)
	}

	if err := query.Find(&items).Error; err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"items": items})
}

func (bc *BizConfigController) CreateLotteryCategory(c *gin.Context) {
	var req struct {
		CategoryKey    string `json:"category_key"`
		Name           string `json:"name"`
		SearchKeywords string `json:"search_keywords"`
		ShowOnHome     *int8  `json:"show_on_home"`
		Status         *int8  `json:"status"`
		Sort           *int   `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		return
	}

	item := models.WLotteryCategory{
		CategoryKey:    strings.TrimSpace(req.CategoryKey),
		Name:           strings.TrimSpace(req.Name),
		SearchKeywords: strings.TrimSpace(req.SearchKeywords),
		ShowOnHome:     1,
		Status:         1,
		Sort:           0,
	}
	if item.CategoryKey == "" || item.Name == "" {
		utils.JSONError(c, http.StatusBadRequest, "category_key/name required")
		return
	}
	if req.ShowOnHome != nil {
		item.ShowOnHome = *req.ShowOnHome
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

func (bc *BizConfigController) UpdateLotteryCategory(c *gin.Context) {
	id, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req struct {
		CategoryKey    *string `json:"category_key"`
		Name           *string `json:"name"`
		SearchKeywords *string `json:"search_keywords"`
		ShowOnHome     *int8   `json:"show_on_home"`
		Status         *int8   `json:"status"`
		Sort           *int    `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		return
	}

	updates := map[string]interface{}{}
	if req.CategoryKey != nil {
		updates["category_key"] = strings.TrimSpace(*req.CategoryKey)
	}
	if req.Name != nil {
		updates["name"] = strings.TrimSpace(*req.Name)
	}
	if req.SearchKeywords != nil {
		updates["search_keywords"] = strings.TrimSpace(*req.SearchKeywords)
	}
	if req.ShowOnHome != nil {
		updates["show_on_home"] = *req.ShowOnHome
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}

	if len(updates) == 0 {
		utils.JSONError(c, http.StatusBadRequest, "empty updates")
		return
	}

	if keyRaw, ok := updates["category_key"]; ok && strings.TrimSpace(keyRaw.(string)) == "" {
		utils.JSONError(c, http.StatusBadRequest, "category_key cannot be empty")
		return
	}
	if nameRaw, ok := updates["name"]; ok && strings.TrimSpace(nameRaw.(string)) == "" {
		utils.JSONError(c, http.StatusBadRequest, "name cannot be empty")
		return
	}

	if err := bc.db.Model(&models.WLotteryCategory{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}

func (bc *BizConfigController) DeleteLotteryCategory(c *gin.Context) {
	id, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := bc.db.Delete(&models.WLotteryCategory{}, id).Error; err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}
