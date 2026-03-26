package jwt

import (
	"errors"
	"go-admin/internal/constants"

	gjwt "github.com/golang-jwt/jwt/v5"
)

// ParseTokenClaims 解析TokenClaims。
func ParseTokenClaims(tokenString string, signingKey string) (*Claims, error) {
	// 定义并初始化当前变量。
	token, err := gjwt.ParseWithClaims(tokenString, &Claims{}, func(t *gjwt.Token) (interface{}, error) {
		// 只允许 HMAC 家族算法，避免 alg 混淆攻击。
		if _, ok := t.Method.(*gjwt.SigningMethodHMAC); !ok {
			// 返回当前处理结果。
			return nil, errors.New("unexpected signing method")
		}
		// 返回当前处理结果。
		return []byte(signingKey), nil
	})
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 将 jwt 的错误语义映射到项目内统一错误码，便于上层控制器处理。
		if errors.Is(err, gjwt.ErrTokenExpired) {
			// 返回当前处理结果。
			return nil, constants.ErrExpiredToken
		}
		// 返回当前处理结果。
		return nil, constants.ErrParsingToken
	}
	// 定义并初始化当前变量。
	claims, ok := token.Claims.(*Claims)
	// 判断条件并进入对应分支逻辑。
	if !ok || !token.Valid {
		// 返回当前处理结果。
		return nil, constants.ErrInvalidToken
	}
	// 返回当前处理结果。
	return claims, nil
}
