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
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil && shouldLogRequestBody(c) {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		c.Next()

		endTime := time.Now()
		latency := endTime.Sub(startTime)

		logger := utils.LoggerFromContext(c.Request.Context())

		// 构建JSON格式的日志
		logEntry := map[string]interface{}{
			"timestamp":   startTime.Format(time.RFC3339),
			"latency":     latency.String(),
			"latency_ms":  latency.Milliseconds(),
			"client_ip":   c.ClientIP(),
			"method":      c.Request.Method,
			"path":        path,
			"query":       query,
			"status_code": c.Writer.Status(),
			"user_agent":  c.Request.UserAgent(),
			"referer":     c.Request.Referer(),
		}

		// 添加错误信息
		if errStr := c.Errors.ByType(gin.ErrorTypePrivate).String(); errStr != "" {
			logEntry["error"] = errStr
		}

		// 添加请求体（脱敏后）
		if len(requestBody) > 0 {
			logEntry["request_body"] = maskSensitiveData(string(requestBody))
		}

		// 记录JSON格式日志
		logger.Info("%s", toJSONString(logEntry))
	}
}

// maskSensitiveData 脱敏敏感数据
func maskSensitiveData(body string) string {
	if body == "" {
		return body
	}

	// 优先按JSON结构脱敏
	var any interface{}
	if err := json.Unmarshal([]byte(body), &any); err == nil {
		maskJSON(any)
		if out, err := json.Marshal(any); err == nil {
			return string(out)
		}
	}

	// JSON解析失败时，使用正则兜底
	replacements := []struct {
		re   *regexp.Regexp
		repl string
	}{
		{regexp.MustCompile(`(?i)"password"\s*:\s*"[^"]*"`), `"password":"***"`},
		{regexp.MustCompile(`(?i)"token"\s*:\s*"[^"]*"`), `"token":"***"`},
		{regexp.MustCompile(`(?i)"refresh_token"\s*:\s*"[^"]*"`), `"refresh_token":"***"`},
		{regexp.MustCompile(`(?i)"secret"\s*:\s*"[^"]*"`), `"secret":"***"`},
		{regexp.MustCompile(`(?i)"api_key"\s*:\s*"[^"]*"`), `"api_key":"***"`},
	}

	masked := body
	for _, item := range replacements {
		masked = item.re.ReplaceAllString(masked, item.repl)
	}
	return masked
}

// toJSONString 将map转换为JSON字符串
func toJSONString(data map[string]interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		return `{"msg":"failed to marshal log"}`
	}
	return string(b)
}

func shouldLogRequestBody(c *gin.Context) bool {
	path := c.Request.URL.Path
	switch path {
	case "/auth/login", "/auth/register", "/auth/refresh":
		return false
	}
	return true
}

func maskJSON(v interface{}) {
	switch t := v.(type) {
	case map[string]interface{}:
		for k, val := range t {
			lk := lower(k)
			if lk == "password" || lk == "token" || lk == "refresh_token" || lk == "secret" || lk == "api_key" {
				t[k] = "***"
				continue
			}
			maskJSON(val)
		}
	case []interface{}:
		for _, item := range t {
			maskJSON(item)
		}
	}
}

func lower(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch >= 'A' && ch <= 'Z' {
			out = append(out, ch+('a'-'A'))
		} else {
			out = append(out, ch)
		}
	}
	return string(out)
}
