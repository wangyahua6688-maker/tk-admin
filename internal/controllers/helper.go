package controllers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

// parseUintID 解析路由参数中的主键ID，并统一校验为正整数。
func parseUintID(c *gin.Context) (uint, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return uint(id), nil
}
