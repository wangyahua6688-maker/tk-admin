package biz

import (
	"strings"
	"time"
)

// parseRFC3339Ptr 将字符串指针解析为 RFC3339 时间，解析失败时返回 nil。
func parseRFC3339Ptr(raw *string) *time.Time {
	// 判断条件并进入对应分支逻辑。
	if raw == nil {
		// 返回当前处理结果。
		return nil
	}
	// 定义并初始化当前变量。
	v := strings.TrimSpace(*raw)
	// 判断条件并进入对应分支逻辑。
	if v == "" {
		// 返回当前处理结果。
		return nil
	}
	// 定义并初始化当前变量。
	t, err := time.Parse(time.RFC3339, v)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return nil
	}
	// 返回当前处理结果。
	return &t
}
