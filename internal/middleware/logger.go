package middleware

import (
	"bytes"
	"fmt"
	"io"
	"strings"
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
		if c.Request.Body != nil {
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

	// 这里可以添加更复杂的脱敏逻辑
	// 例如：将密码、token、身份证号等敏感信息替换为***

	// 简单的JSON密码字段脱敏
	patterns := map[string]string{
		`"password":\s*".*?"`:      `"password":"***"`,
		`"token":\s*".*?"`:         `"token":"***"`,
		`"refresh_token":\s*".*?"`: `"refresh_token":"***"`,
		`"secret":\s*".*?"`:        `"secret":"***"`,
		`"api_key":\s*".*?"`:       `"api_key":"***"`,
	}

	maskedBody := body
	for pattern, replacement := range patterns {
		maskedBody = strings.ReplaceAll(maskedBody, pattern, replacement)
	}

	return maskedBody
}

// toJSONString 将map转换为JSON字符串
func toJSONString(data map[string]interface{}) string {
	// 这里可以使用json.Marshal，简化示例使用fmt
	jsonStr := "{"
	first := true
	for k, v := range data {
		if !first {
			jsonStr += ", "
		}
		jsonStr += fmt.Sprintf(`"%s": "%v"`, k, v)
		first = false
	}
	jsonStr += "}"
	return jsonStr
}
