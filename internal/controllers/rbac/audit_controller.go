package rbac

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	rbacdao "go-admin-full/internal/dao/rbac"
	rbacsvc "go-admin-full/internal/services/rbac"
	"go-admin-full/internal/utils"
	"gorm.io/gorm"
)

// AuditController 审计控制器（当前提供登录日志查询）。
type AuditController struct {
	loginLogSvc *rbacsvc.LoginLogService // 登录日志服务
}

func NewAuditController(db *gorm.DB) *AuditController {
	return &AuditController{
		loginLogSvc: rbacsvc.NewLoginLogService(rbacdao.NewLoginLogDao(db)), // 注入登录日志服务
	}
}

// ListLoginLogs 分页查询登录日志。
func (ac *AuditController) ListLoginLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))           // 页码
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20")) // 每页数量
	username := strings.TrimSpace(c.Query("username"))             // 用户名筛选

	if page < 1 {
		utils.JSONError(c, http.StatusBadRequest, "page必须大于0")
		return
	}
	if pageSize < 1 || pageSize > 100 {
		utils.JSONError(c, http.StatusBadRequest, "page_size范围必须在1-100")
		return
	}

	// 查询登录日志
	logs, total, err := ac.loginLogSvc.ListLoginLogs(c.Request.Context(), page, pageSize, username)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 返回分页数据
	utils.JSONOK(c, gin.H{
		"list":      logs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
