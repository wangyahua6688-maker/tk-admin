package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/utils"
)

// RequestInfo 请求信息结构
type RequestInfo struct {
	Timestamp    time.Time     // 请求时间
	Latency      time.Duration // 处理耗时
	ClientIP     string        // 客户端IP
	Method       string        // HTTP方法
	Path         string        // 请求路径
	StatusCode   int           // 状态码
	RequestSize  int64         // 请求大小
	ResponseSize int           // 响应大小
	RequestBody  string        // 请求体（脱敏后）
	Error        string        // 错误信息
	UserAgent    string        // 用户代理
	Referer      string        // 来源
}

// JSONLoggerMiddleware JSON格式日志中间件
func JSONLoggerMiddleware() gin.HandlerFunc {
	// 返回当前处理结果。
	return func(c *gin.Context) {
		// 定义并初始化当前变量。
		startTime := time.Now()
		// 定义并初始化当前变量。
		path := c.Request.URL.Path
		// 定义并初始化当前变量。
		query := c.Request.URL.RawQuery

		// 读取请求体
		var requestBody []byte
		// 判断条件并进入对应分支逻辑。
		if c.Request.Body != nil && shouldLogRequestBody(c) {
			// 更新当前变量或字段值。
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 更新当前变量或字段值。
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 调用c.Next完成当前处理。
		c.Next()

		// 定义并初始化当前变量。
		endTime := time.Now()
		// 定义并初始化当前变量。
		latency := endTime.Sub(startTime)

		// 定义并初始化当前变量。
		logger := utils.LoggerFromContext(c.Request.Context())

		// 构建JSON格式的日志
		logEntry := map[string]interface{}{
			// 调用startTime.Format完成当前处理。
			"timestamp": startTime.Format(time.RFC3339),
			// 调用latency.String完成当前处理。
			"latency": latency.String(),
			// 调用latency.Milliseconds完成当前处理。
			"latency_ms": latency.Milliseconds(),
			// 调用c.ClientIP完成当前处理。
			"client_ip": c.ClientIP(),
			// 处理当前语句逻辑。
			"method": c.Request.Method,
			// 处理当前语句逻辑。
			"path": path,
			// 处理当前语句逻辑。
			"query": query,
			// 调用c.Writer.Status完成当前处理。
			"status_code": c.Writer.Status(),
			// 调用c.Request.UserAgent完成当前处理。
			"user_agent": c.Request.UserAgent(),
			// 调用c.Request.Referer完成当前处理。
			"referer": c.Request.Referer(),
		}

		// 添加错误信息
		if errStr := c.Errors.ByType(gin.ErrorTypePrivate).String(); errStr != "" {
			// 更新当前变量或字段值。
			logEntry["error"] = errStr
		}

		// 添加请求体（脱敏后）
		if len(requestBody) > 0 {
			// 更新当前变量或字段值。
			logEntry["request_body"] = maskSensitiveData(string(requestBody))
		}

		// 记录JSON格式日志
		logger.Info("%s", toJSONString(logEntry))
	}
}

// maskSensitiveData 脱敏敏感数据
func maskSensitiveData(body string) string {
	// 判断条件并进入对应分支逻辑。
	if body == "" {
		// 返回当前处理结果。
		return body
	}

	// 优先按JSON结构脱敏
	var any interface{}
	// 判断条件并进入对应分支逻辑。
	if err := json.Unmarshal([]byte(body), &any); err == nil {
		// 调用maskJSON完成当前处理。
		maskJSON(any)
		// 判断条件并进入对应分支逻辑。
		if out, err := json.Marshal(any); err == nil {
			// 返回当前处理结果。
			return string(out)
		}
	}

	// JSON解析失败时，使用正则兜底
	replacements := []struct {
		// 处理当前语句逻辑。
		re *regexp.Regexp
		// 处理当前语句逻辑。
		repl string
	}{
		// 调用regexp.MustCompile完成当前处理。
		{regexp.MustCompile(`(?i)"password"\s*:\s*"[^"]*"`), `"password":"***"`},
		// 调用regexp.MustCompile完成当前处理。
		{regexp.MustCompile(`(?i)"token"\s*:\s*"[^"]*"`), `"token":"***"`},
		// 调用regexp.MustCompile完成当前处理。
		{regexp.MustCompile(`(?i)"refresh_token"\s*:\s*"[^"]*"`), `"refresh_token":"***"`},
		// 调用regexp.MustCompile完成当前处理。
		{regexp.MustCompile(`(?i)"secret"\s*:\s*"[^"]*"`), `"secret":"***"`},
		// 调用regexp.MustCompile完成当前处理。
		{regexp.MustCompile(`(?i)"api_key"\s*:\s*"[^"]*"`), `"api_key":"***"`},
	}

	// 定义并初始化当前变量。
	masked := body
	// 循环处理当前数据集合。
	for _, item := range replacements {
		// 更新当前变量或字段值。
		masked = item.re.ReplaceAllString(masked, item.repl)
	}
	// 返回当前处理结果。
	return masked
}

// toJSONString 将map转换为JSON字符串
func toJSONString(data map[string]interface{}) string {
	// 定义并初始化当前变量。
	b, err := json.Marshal(data)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return `{"msg":"failed to marshal log"}`
	}
	// 返回当前处理结果。
	return string(b)
}

// shouldLogRequestBody 处理shouldLogRequestBody相关逻辑。
func shouldLogRequestBody(c *gin.Context) bool {
	// 定义并初始化当前变量。
	path := c.Request.URL.Path
	// 根据表达式进入多分支处理。
	switch path {
	case "/auth/login", "/auth/register", "/auth/refresh":
		// 返回当前处理结果。
		return false
	}
	// 返回当前处理结果。
	return true
}

// maskJSON 处理maskJSON相关逻辑。
func maskJSON(v interface{}) {
	// 根据表达式进入多分支处理。
	switch t := v.(type) {
	case map[string]interface{}:
		// 循环处理当前数据集合。
		for k, val := range t {
			// 定义并初始化当前变量。
			lk := lower(k)
			// 判断条件并进入对应分支逻辑。
			if lk == "password" || lk == "token" || lk == "refresh_token" || lk == "secret" || lk == "api_key" {
				// 更新当前变量或字段值。
				t[k] = "***"
				// 处理当前语句逻辑。
				continue
			}
			// 调用maskJSON完成当前处理。
			maskJSON(val)
		}
	case []interface{}:
		// 循环处理当前数据集合。
		for _, item := range t {
			// 调用maskJSON完成当前处理。
			maskJSON(item)
		}
	}
}

// lower 处理lower相关逻辑。
func lower(s string) string {
	// 定义并初始化当前变量。
	out := make([]byte, 0, len(s))
	// 循环处理当前数据集合。
	for i := 0; i < len(s); i++ {
		// 定义并初始化当前变量。
		ch := s[i]
		// 判断条件并进入对应分支逻辑。
		if ch >= 'A' && ch <= 'Z' {
			// 更新当前变量或字段值。
			out = append(out, ch+('a'-'A'))
			// 进入新的代码块进行处理。
		} else {
			// 更新当前变量或字段值。
			out = append(out, ch)
		}
	}
	// 返回当前处理结果。
	return string(out)
}
