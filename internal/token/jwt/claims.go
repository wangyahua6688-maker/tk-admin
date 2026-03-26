package jwt

// 处理当前语句逻辑。
import gjwt "github.com/golang-jwt/jwt/v5"

// Claims 扩展标准 JWT 声明，承载业务字段。
type Claims struct {
	// 处理当前语句逻辑。
	UserID    uint   `json:"user_id"`
	TokenType string `json:"token_type"`          // access | refresh
	DeviceID  string `json:"device_id,omitempty"` // refresh token device scope
	// 处理当前语句逻辑。
	gjwt.RegisteredClaims
}
