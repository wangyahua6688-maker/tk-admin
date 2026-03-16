package biz

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/models"
	"go-admin-full/internal/utils"
)

// -------------------- Broadcast --------------------

func (bc *BizConfigController) ListBroadcasts(c *gin.Context) {
	items, err := bc.svc.ListBroadcasts(c.Request.Context(), 200)
	if err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"items": items})
}

func (bc *BizConfigController) CreateBroadcast(c *gin.Context) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Status  *int8  `json:"status"`
		Sort    *int   `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		return
	}
	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.Content) == "" {
		utils.JSONError(c, http.StatusBadRequest, "title/content required")
		return
	}
	item := models.WBroadcast{Title: strings.TrimSpace(req.Title), Content: strings.TrimSpace(req.Content), Status: 1, Sort: 0}
	if req.Status != nil {
		item.Status = *req.Status
	}
	if req.Sort != nil {
		item.Sort = *req.Sort
	}
	if err := bc.svc.CreateBroadcast(c.Request.Context(), &item); err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, item)
}

func (bc *BizConfigController) UpdateBroadcast(c *gin.Context) {
	id, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		Title   *string `json:"title"`
		Content *string `json:"content"`
		Status  *int8   `json:"status"`
		Sort    *int    `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		return
	}
	updates := map[string]interface{}{}
	if req.Title != nil {
		updates["title"] = strings.TrimSpace(*req.Title)
	}
	if req.Content != nil {
		updates["content"] = strings.TrimSpace(*req.Content)
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
	if err := bc.svc.UpdateBroadcast(c.Request.Context(), id, updates); err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}

func (bc *BizConfigController) DeleteBroadcast(c *gin.Context) {
	id, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := bc.svc.DeleteBroadcast(c.Request.Context(), id); err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}
