package jwt

import (
	"errors"
	"go-admin-full/internal/constants"

	gjwt "github.com/golang-jwt/jwt/v5"
)

func ParseTokenClaims(tokenString string, signingKey string) (*Claims, error) {
	token, err := gjwt.ParseWithClaims(tokenString, &Claims{}, func(t *gjwt.Token) (interface{}, error) {
		// 只允许 HMAC 家族算法，避免 alg 混淆攻击。
		if _, ok := t.Method.(*gjwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		// 将 jwt 的错误语义映射到项目内统一错误码，便于上层控制器处理。
		if errors.Is(err, gjwt.ErrTokenExpired) {
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

// ParseToken 兼容方法：直接返回 userID。
// 新代码建议优先调用 ParseTokenClaims，以便拿到 token_type / device_id / jti 等上下文。
func ParseToken(tokenString string, signingKey string) (uint, error) {
	claims, err := ParseTokenClaims(tokenString, signingKey)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
