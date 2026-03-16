package biz

import (
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/models"
	"go-admin-full/internal/utils"
)

// -------------------- Banner --------------------

func (bc *BizConfigController) ListBanners(c *gin.Context) {
	bannerType := strings.TrimSpace(c.Query("type"))
	items, err := bc.svc.ListBanners(c.Request.Context(), bannerType, 300)
	if err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"items": items})
}

func (bc *BizConfigController) CreateBanner(c *gin.Context) {
	var req struct {
		Title       string   `json:"title"`
		ImageURL    string   `json:"image_url"`
		LinkURL     string   `json:"link_url"`
		Type        string   `json:"type"`
		Position    string   `json:"position"`
		Positions   []string `json:"positions"`
		JumpType    string   `json:"jump_type"`
		JumpPostID  uint     `json:"jump_post_id"`
		JumpURL     string   `json:"jump_url"`
		ContentHTML string   `json:"content_html"`
		Status      *int8    `json:"status"`
		Sort        *int     `json:"sort"`
		StartAt     *string  `json:"start_at"`
		EndAt       *string  `json:"end_at"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		return
	}
	positions := normalizePositions(req.Positions, req.Position)
	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.ImageURL) == "" || strings.TrimSpace(req.Type) == "" || len(positions) == 0 {
		utils.JSONError(c, http.StatusBadRequest, "title/image_url/type/positions required")
		return
	}

	jumpType := normalizeJumpType(req.JumpType)
	if jumpType == "post" && req.JumpPostID == 0 {
		utils.JSONError(c, http.StatusBadRequest, "jump_post_id required when jump_type=post")
		return
	}
	if jumpType == "external" && strings.TrimSpace(req.JumpURL) == "" {
		utils.JSONError(c, http.StatusBadRequest, "jump_url required when jump_type=external")
		return
	}
	if jumpType == "post" {
		ok, err := bc.svc.IsPostExists(c.Request.Context(), req.JumpPostID)
		if err != nil {
			utils.JSONError(c, http.StatusInternalServerError, err.Error())
			return
		}
		if !ok {
			utils.JSONError(c, http.StatusBadRequest, "jump_post_id not found")
			return
		}
	}

	item := models.WBanner{
		Title:       strings.TrimSpace(req.Title),
		ImageURL:    strings.TrimSpace(req.ImageURL),
		LinkURL:     strings.TrimSpace(req.LinkURL),
		Type:        strings.TrimSpace(req.Type),
		Position:    positions[0],
		Positions:   strings.Join(positions, ","),
		JumpType:    jumpType,
		JumpPostID:  req.JumpPostID,
		JumpURL:     strings.TrimSpace(req.JumpURL),
		ContentHTML: req.ContentHTML,
		Status:      1,
		Sort:        0,
		StartAt:     parseRFC3339Ptr(req.StartAt),
		EndAt:       parseRFC3339Ptr(req.EndAt),
	}
	if jumpType == "external" {
		item.LinkURL = strings.TrimSpace(req.JumpURL)
	}
	if req.Status != nil {
		item.Status = *req.Status
	}
	if req.Sort != nil {
		item.Sort = *req.Sort
	}
	if err := bc.svc.CreateBanner(c.Request.Context(), &item); err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, item)
}

func (bc *BizConfigController) UpdateBanner(c *gin.Context) {
	id, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		Title       *string  `json:"title"`
		ImageURL    *string  `json:"image_url"`
		LinkURL     *string  `json:"link_url"`
		Type        *string  `json:"type"`
		Position    *string  `json:"position"`
		Positions   []string `json:"positions"`
		JumpType    *string  `json:"jump_type"`
		JumpPostID  *uint    `json:"jump_post_id"`
		JumpURL     *string  `json:"jump_url"`
		ContentHTML *string  `json:"content_html"`
		Status      *int8    `json:"status"`
		Sort        *int     `json:"sort"`
		StartAt     *string  `json:"start_at"`
		EndAt       *string  `json:"end_at"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		return
	}

	updates := map[string]interface{}{}
	if req.Title != nil {
		updates["title"] = strings.TrimSpace(*req.Title)
	}
	if req.ImageURL != nil {
		updates["image_url"] = strings.TrimSpace(*req.ImageURL)
	}
	if req.LinkURL != nil {
		updates["link_url"] = strings.TrimSpace(*req.LinkURL)
	}
	if req.Type != nil {
		updates["type"] = strings.TrimSpace(*req.Type)
	}
	if req.Position != nil {
		updates["position"] = strings.TrimSpace(*req.Position)
	}
	if req.Positions != nil {
		positions := normalizePositions(req.Positions, "")
		if len(positions) > 0 {
			updates["positions"] = strings.Join(positions, ",")
			updates["position"] = positions[0]
		}
	}
	if req.JumpType != nil {
		updates["jump_type"] = normalizeJumpType(*req.JumpType)
	}
	if req.JumpPostID != nil {
		updates["jump_post_id"] = *req.JumpPostID
	}
	if req.JumpURL != nil {
		updates["jump_url"] = strings.TrimSpace(*req.JumpURL)
	}
	if req.ContentHTML != nil {
		updates["content_html"] = *req.ContentHTML
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}
	if req.StartAt != nil {
		updates["start_at"] = parseRFC3339Ptr(req.StartAt)
	}
	if req.EndAt != nil {
		updates["end_at"] = parseRFC3339Ptr(req.EndAt)
	}
	if link, ok := updates["jump_url"]; ok {
		if jt, hasJT := updates["jump_type"]; hasJT && jt == "external" {
			updates["link_url"] = link
		}
	}
	if len(updates) == 0 {
		utils.JSONError(c, http.StatusBadRequest, "empty updates")
		return
	}

	if err := bc.svc.UpdateBanner(c.Request.Context(), id, updates); err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}

func (bc *BizConfigController) DeleteBanner(c *gin.Context) {
	id, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := bc.svc.DeleteBanner(c.Request.Context(), id); err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}

func normalizePositions(positions []string, fallback string) []string {
	out := make([]string, 0)
	seen := map[string]struct{}{}
	appendPos := func(v string) {
		pos := strings.TrimSpace(v)
		if pos == "" {
			return
		}
		if _, ok := seen[pos]; ok {
			return
		}
		seen[pos] = struct{}{}
		out = append(out, pos)
	}
	for _, p := range positions {
		appendPos(p)
	}
	if len(out) == 0 {
		appendPos(fallback)
	}
	sort.Strings(out)
	return out
}

func normalizeJumpType(v string) string {
	switch strings.TrimSpace(v) {
	case "post", "external", "custom", "none":
		return strings.TrimSpace(v)
	default:
		return "none"
	}
}
