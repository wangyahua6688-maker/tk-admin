package utils

import "github.com/gin-gonic/gin"

func JSONOK(c *gin.Context, data interface{}) {
    c.JSON(200, gin.H{"code": 0, "msg": "ok", "data": data})
}

func JSONError(c *gin.Context, code int, msg string) {
    c.JSON(200, gin.H{"code": code, "msg": msg})
}
