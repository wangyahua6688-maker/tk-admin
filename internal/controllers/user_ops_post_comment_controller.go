package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/models"
	"go-admin-full/internal/utils"
)

// -------------------- 帖子评论管理（按帖子维度） --------------------

func (uc *UserOpsController) ListPostComments(c *gin.Context) {
	postID, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid post id")
		return
	}

	var rows []models.WComment
	if err := uc.db.Where("post_id = ?", postID).Order("parent_id ASC, id ASC").Limit(500).Find(&rows).Error; err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}

	userIDs := make([]uint, 0)
	seen := make(map[uint]struct{})
	for _, row := range rows {
		if _, ok := seen[row.UserID]; ok {
			continue
		}
		seen[row.UserID] = struct{}{}
		userIDs = append(userIDs, row.UserID)
	}

	userMap := map[uint]models.WUser{}
	if len(userIDs) > 0 {
		var users []models.WUser
		_ = uc.db.Select("id,username,nickname,user_type").Where("id IN ?", userIDs).Find(&users).Error
		for _, u := range users {
			userMap[u.ID] = u
		}
	}

	items := make([]map[string]interface{}, 0, len(rows))
	for _, row := range rows {
		u := userMap[row.UserID]
		items = append(items, map[string]interface{}{
			"id":         row.ID,
			"post_id":    row.PostID,
			"user_id":    row.UserID,
			"parent_id":  row.ParentID,
			"content":    row.Content,
			"likes":      row.Likes,
			"status":     row.Status,
			"created_at": row.CreatedAt,
			"username":   u.Username,
			"nickname":   u.Nickname,
			"user_type":  u.UserType,
		})
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

	if !uc.isUserTypes(req.UserID, "robot", "official") {
		utils.JSONError(c, http.StatusBadRequest, "user must be robot or official account")
		return
	}

	var post models.WPostArticle
	if err := uc.db.Select("id").Where("id = ?", postID).First(&post).Error; err != nil {
		utils.JSONError(c, http.StatusBadRequest, "post not found")
		return
	}

	item := models.WComment{
		PostID:   postID,
		UserID:   req.UserID,
		ParentID: req.ParentID,
		Content:  strings.TrimSpace(req.Content),
		Likes:    0,
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
