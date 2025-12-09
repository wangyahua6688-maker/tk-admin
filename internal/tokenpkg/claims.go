package tokenpkg

import "github.com/golang-jwt/jwt/v5"

// Claims extends jwt.RegisteredClaims with app-specific fields.
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}
