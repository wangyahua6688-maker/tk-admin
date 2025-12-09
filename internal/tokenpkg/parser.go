package tokenpkg

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// ParseToken parses tokenString using signingKey and returns userID if valid.
// It maps jwt.ErrTokenExpired to ErrExpiredToken and returns ErrParsingToken for other parse errors.
func ParseToken(tokenString string, signingKey string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		// enforce HMAC signing method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		// jwt v5 exports jwt.ErrTokenExpired for expired tokens
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, ErrExpiredToken
		}
		return 0, ErrParsingToken
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return 0, ErrInvalidToken
	}
	return claims.UserID, nil
}
