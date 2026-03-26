package rbac

import (
	"strconv"
	"strings"

	commonresp "github.com/wangyahua6688-maker/tk-common/utils/httpresp"
	"go-admin-full/internal/constants"
	rbacdao "go-admin-full/internal/dao/rbac"
	rbacsvc "go-admin-full/internal/services/rbac"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuditController 审计控制器（当前提供登录日志查询）。
type AuditController struct {
	loginLogSvc *rbacsvc.LoginLogService // 登录日志服务
}

// NewAuditController 创建AuditController实例。
func NewAuditController(db *gorm.DB) *AuditController {
	// 返回当前处理结果。
	return &AuditController{
		loginLogSvc: rbacsvc.NewLoginLogService(rbacdao.NewLoginLogDao(db)), // 注入登录日志服务
	}
}

// ListLoginLogs 分页查询登录日志。
func (ac *AuditController) ListLoginLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))           // 页码
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20")) // 每页数量
	username := strings.TrimSpace(c.Query("username"))             // 用户名筛选

	// 判断条件并进入对应分支逻辑。
	if page < 1 {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "page必须大于0")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if pageSize < 1 || pageSize > 100 {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "page_size范围必须在1-100")
		// 返回当前处理结果。
		return
	}

	// 查询登录日志
	logs, total, err := ac.loginLogSvc.ListLoginLogs(c.Request.Context(), page, pageSize, username)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}

	// 返回分页数据
	commonresp.GinOK(c, gin.H{
		// 处理当前语句逻辑。
		"list": logs,
		// 处理当前语句逻辑。
		"total": total,
		// 处理当前语句逻辑。
		"page": page,
		// 处理当前语句逻辑。
		"page_size": pageSize,
	})
}
