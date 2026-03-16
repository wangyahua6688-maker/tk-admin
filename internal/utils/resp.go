package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	commonresp "tk-common/utils/httpresp"
)

func JSONOK(c *gin.Context, data interface{}) {
	commonresp.GinOK(c, data)
}

func JSONError(c *gin.Context, code int, msg string) {
	commonresp.GinError(c, code, msg)
}

func OK(c *gin.Context, data interface{}) {
	commonresp.GinOK(c, data)
}

func Error(c *gin.Context, msg string) {
	commonresp.GinFailWithStatus(c, http.StatusBadRequest, http.StatusBadRequest, msg)
}

func Forbidden(c *gin.Context, msg string) {
	commonresp.GinFailWithStatus(c, http.StatusForbidden, http.StatusForbidden, msg)
}
