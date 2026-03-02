package services

import (
	"errors"
	"unicode"
)

// ValidatePasswordStrength 校验密码强度。
// 规则：
// 1. 长度 8-72（bcrypt 有 72 字节上限）
// 2. 必须包含大小写字母和数字
func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("密码长度至少8位")
	}
	if len(password) > 72 {
		return errors.New("密码长度不能超过72位")
	}

	var hasUpper, hasLower, hasDigit bool
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit {
		return errors.New("密码需包含大小写字母和数字")
	}
	return nil
}
