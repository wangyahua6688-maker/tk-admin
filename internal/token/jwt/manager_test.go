package jwt

import (
	"testing"
	"time"

	"go-admin-full/internal/constants"
	tokenstore "go-admin-full/internal/token/store"
)

func TestNormalizeDeviceID_IdempotentForHashedInput(t *testing.T) {
	raw := normalizeDeviceID("my-device")
	again := normalizeDeviceID(raw)
	if raw != again {
		t.Fatalf("normalizeDeviceID should be idempotent for normalized token, got %s != %s", raw, again)
	}
}

func TestRefreshTokenPairFailsAfterSessionIdleTimeout(t *testing.T) {
	cfg := DefaultConfig()
	cfg.SigningKey = "this-is-a-test-signing-key-with-enough-length-123456"
	cfg.AccessExpire = 2 * time.Second
	cfg.RefreshExpire = 10 * time.Second
	cfg.SessionIdleTimeout = 200 * time.Millisecond

	mgr := NewManager(cfg, tokenstore.NewMemoryStore())
	_, refresh, err := mgr.GenerateTokensWithDevice(1001, "device-A")
	if err != nil {
		t.Fatalf("GenerateTokensWithDevice failed: %v", err)
	}

	time.Sleep(320 * time.Millisecond)
	_, _, err = mgr.RefreshTokenPair(refresh)
	if err != constants.ErrExpiredToken {
		t.Fatalf("expected ErrExpiredToken, got: %v", err)
	}
}

func TestValidateAccessTokenTouchesSessionTTL(t *testing.T) {
	cfg := DefaultConfig()
	cfg.SigningKey = "this-is-a-test-signing-key-with-enough-length-654321"
	cfg.AccessExpire = 2 * time.Second
	cfg.RefreshExpire = 10 * time.Second
	cfg.SessionIdleTimeout = 250 * time.Millisecond

	mgr := NewManager(cfg, tokenstore.NewMemoryStore())
	access, refresh, err := mgr.GenerateTokensWithDevice(2002, "device-B")
	if err != nil {
		t.Fatalf("GenerateTokensWithDevice failed: %v", err)
	}

	time.Sleep(150 * time.Millisecond)
	if _, err := mgr.ValidateAccessToken(access); err != nil {
		t.Fatalf("ValidateAccessToken failed unexpectedly: %v", err)
	}

	time.Sleep(150 * time.Millisecond)
	if _, _, err := mgr.RefreshTokenPair(refresh); err != nil {
		t.Fatalf("RefreshTokenPair should succeed after session touch, got: %v", err)
	}
}
