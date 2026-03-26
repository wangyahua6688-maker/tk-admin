package constants

import "errors"

// 常用错误
var (
	// 更新当前变量或字段值。
	ErrInvalidToken = errors.New("token is invalid")
	// 更新当前变量或字段值。
	ErrExpiredToken = errors.New("token is expired")
	// 更新当前变量或字段值。
	ErrSigningToken = errors.New("failed to sign token")
	// 更新当前变量或字段值。
	ErrParsingToken = errors.New("failed to parse token")
	// 更新当前变量或字段值。
	ErrTokenNotFound = errors.New("token not found in store")
	// 更新当前变量或字段值。
	ErrTokenStoreFailed = errors.New("failed to store token")
)
