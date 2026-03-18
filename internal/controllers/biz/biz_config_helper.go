package biz

import (
	"fmt"
	"net/url"
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

// normalizeSafeURL 归一化并校验 URL，仅允许 http/https 和可选站内相对路径。
func normalizeSafeURL(raw string, allowRelativePath bool) (string, error) {
	// 定义并初始化当前变量。
	trimmed := strings.TrimSpace(raw)
	// 判断条件并进入对应分支逻辑。
	if trimmed == "" {
		// 返回当前处理结果。
		return "", nil
	}
	// 判断条件并进入对应分支逻辑。
	if strings.HasPrefix(trimmed, "//") {
		// 返回当前处理结果。
		return "", fmt.Errorf("url scheme not allowed")
	}
	// 判断条件并进入对应分支逻辑。
	if strings.HasPrefix(trimmed, "/") {
		// 判断条件并进入对应分支逻辑。
		if allowRelativePath {
			// 返回当前处理结果。
			return trimmed, nil
		}
		// 返回当前处理结果。
		return "", fmt.Errorf("relative url not allowed")
	}
	// 定义并初始化当前变量。
	parsed, err := url.Parse(trimmed)
	// 判断条件并进入对应分支逻辑。
	if err != nil || parsed == nil {
		// 返回当前处理结果。
		return "", fmt.Errorf("invalid url")
	}
	// 定义并初始化当前变量。
	scheme := strings.ToLower(strings.TrimSpace(parsed.Scheme))
	// 判断条件并进入对应分支逻辑。
	if scheme != "http" && scheme != "https" {
		// 返回当前处理结果。
		return "", fmt.Errorf("url scheme not allowed")
	}
	// 判断条件并进入对应分支逻辑。
	if strings.TrimSpace(parsed.Host) == "" {
		// 返回当前处理结果。
		return "", fmt.Errorf("invalid url")
	}
	// 返回当前处理结果。
	return trimmed, nil
}
