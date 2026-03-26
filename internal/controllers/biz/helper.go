package biz

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

// parseUintID 解析路由参数中的主键ID，并统一校验为正整数。
func parseUintID(c *gin.Context) (uint, error) {
	// 定义并初始化当前变量。
	id, err := strconv.Atoi(c.Param("id"))
	// 判断条件并进入对应分支逻辑。
	if err != nil || id <= 0 {
		// 返回当前处理结果。
		return 0, errors.New("invalid id")
	}
	// 返回当前处理结果。
	return uint(id), nil
}
