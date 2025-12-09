package utils

import "github.com/gin-gonic/gin"

func JSONOK(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{"code": 0, "msg": "ok", "data": data})
}

func JSONError(c *gin.Context, code int, msg string) {
	c.JSON(200, gin.H{"code": code, "msg": msg})
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}

func Error(c *gin.Context, msg string) {
	c.JSON(400, gin.H{
		"code": 400,
		"msg":  msg,
	})
}

func Forbidden(c *gin.Context, msg string) {
	c.JSON(403, gin.H{
		"code": 403,
		"msg":  msg,
	})
}
