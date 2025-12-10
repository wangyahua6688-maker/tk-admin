package tokenpkg

import "github.com/golang-jwt/jwt/v5"

// Claims extends jwt.RegisteredClaims with app-specific fields.
type Claims struct {
	UserID    uint   `json:"user_id"`
	TokenType string `json:"token_type"`          // access | refresh
	DeviceID  string `json:"device_id,omitempty"` // refresh token device scope
	jwt.RegisteredClaims
}
