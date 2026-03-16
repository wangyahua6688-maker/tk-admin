package biz

import (
	"strings"
	"time"
)

func parseRFC3339Ptr(raw *string) *time.Time {
	if raw == nil {
		return nil
	}
	v := strings.TrimSpace(*raw)
	if v == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return nil
	}
	return &t
}
