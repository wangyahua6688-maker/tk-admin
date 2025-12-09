package tokenpkg

import "errors"

// 常用错误
var (
	ErrInvalidToken     = errors.New("token is invalid")
	ErrExpiredToken     = errors.New("token is expired")
	ErrSigningToken     = errors.New("failed to sign token")
	ErrParsingToken     = errors.New("failed to parse token")
	ErrTokenNotFound    = errors.New("token not found in store")
	ErrTokenStoreFailed = errors.New("failed to store token")
)
