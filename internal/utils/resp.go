package utils

import (
	"github.com/gin-gonic/gin"
	commonresp "tk-common/utils/httpresp"
)

// JSONOK 处理JSONOK相关逻辑。
func JSONOK(c *gin.Context, data interface{}) {
	// 调用commonresp.GinOK完成当前处理。
	commonresp.GinOK(c, data)
}

// JSONError 处理JSONError相关逻辑。
func JSONError(c *gin.Context, code int, msg string) {
	// 调用commonresp.GinError完成当前处理。
	commonresp.GinError(c, code, msg)
}
