package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/dao"
	"go-admin-full/internal/services"
	"go-admin-full/internal/utils"
	"gorm.io/gorm"
)

// AuditController 审计控制器（当前提供登录日志查询）。
type AuditController struct {
	loginLogSvc *services.LoginLogService
}

func NewAuditController(db *gorm.DB) *AuditController {
	return &AuditController{
		loginLogSvc: services.NewLoginLogService(dao.NewLoginLogDao(db)),
	}
}

// ListLoginLogs 分页查询登录日志。
func (ac *AuditController) ListLoginLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	username := strings.TrimSpace(c.Query("username"))

	if page < 1 {
		utils.JSONError(c, http.StatusBadRequest, "page必须大于0")
		return
	}
	if pageSize < 1 || pageSize > 100 {
		utils.JSONError(c, http.StatusBadRequest, "page_size范围必须在1-100")
		return
	}

	logs, total, err := ac.loginLogSvc.ListLoginLogs(c.Request.Context(), page, pageSize, username)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONOK(c, gin.H{
		"list":      logs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
