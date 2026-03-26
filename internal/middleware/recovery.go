package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	commonlogx "github.com/wangyahua6688-maker/tk-common/utils/logx"
)

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	// 返回当前处理结果。
	return func(c *gin.Context) {
		// 注册延迟执行逻辑。
		defer func() {
			// 判断条件并进入对应分支逻辑。
			if err := recover(); err != nil {
				// 获取堆栈信息
				stack := debug.Stack()

				// 记录错误日志
				logger := commonlogx.GetLogger()
				// 调用logger.Error完成当前处理。
				logger.Error("Panic recovered: %v\n%s", err, stack)

				// 返回500错误
				c.JSON(http.StatusInternalServerError, gin.H{
					// 处理当前语句逻辑。
					"code": 500,
					// 处理当前语句逻辑。
					"msg": "服务器内部错误",
				})
				// 调用c.Abort完成当前处理。
				c.Abort()
			}
		}()

		// 调用c.Next完成当前处理。
		c.Next()
	}
}
