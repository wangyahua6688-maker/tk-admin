package rbac

import (
	"errors"
	"unicode"
)

// ValidatePasswordStrength 校验密码强度。
// 规则：
// 1. 长度 8-72（bcrypt 有 72 字节上限）
// 2. 必须包含大小写字母和数字
func ValidatePasswordStrength(password string) error {
	// 判断条件并进入对应分支逻辑。
	if len(password) < 8 {
		// 返回当前处理结果。
		return errors.New("密码长度至少8位")
	}
	// 判断条件并进入对应分支逻辑。
	if len(password) > 72 {
		// 返回当前处理结果。
		return errors.New("密码长度不能超过72位")
	}

	// 声明当前变量。
	var hasUpper, hasLower, hasDigit bool
	// 循环处理当前数据集合。
	for _, r := range password {
		// 根据表达式进入多分支处理。
		switch {
		case unicode.IsUpper(r):
			// 更新当前变量或字段值。
			hasUpper = true
		case unicode.IsLower(r):
			// 更新当前变量或字段值。
			hasLower = true
		case unicode.IsDigit(r):
			// 更新当前变量或字段值。
			hasDigit = true
		}
	}

	// 判断条件并进入对应分支逻辑。
	if !hasUpper || !hasLower || !hasDigit {
		// 返回当前处理结果。
		return errors.New("密码需包含大小写字母和数字")
	}
	// 返回当前处理结果。
	return nil
}
