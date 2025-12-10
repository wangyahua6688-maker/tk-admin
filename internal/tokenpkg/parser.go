package tokenpkg

import (
	"errors"
	"go-admin-full/internal/constants"

	"github.com/golang-jwt/jwt/v5"
)

func ParseTokenClaims(tokenString string, signingKey string) (*Claims, error) {
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
			return nil, constants.ErrExpiredToken
		}
		return nil, constants.ErrParsingToken
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, constants.ErrInvalidToken
	}
	return claims, nil
}

// ParseToken parses tokenString using signingKey and returns userID if valid.
// It maps jwt.ErrTokenExpired to ErrExpiredToken and returns ErrParsingToken for other parse errors.
func ParseToken(tokenString string, signingKey string) (uint, error) {
	claims, err := ParseTokenClaims(tokenString, signingKey)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
