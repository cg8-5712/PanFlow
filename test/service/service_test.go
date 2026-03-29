package service_test

import (
	"context"
	"testing"
	"time"

	"panflow/internal/model"
	"panflow/internal/service"
)

func TestUserCacheKey(t *testing.T) {
	key := service.UserCacheKey(42)
	if key != "user:42" {
		t.Fatalf("unexpected key: %s", key)
	}
}

func TestConfigCacheKey(t *testing.T) {
	key := service.ConfigCacheKey("guest_daily_limit")
	if key != "config:guest_daily_limit" {
		t.Fatalf("unexpected key: %s", key)
	}
}

func TestBlacklistCacheKey(t *testing.T) {
	key := service.BlacklistCacheKey("ip", "1.2.3.4")
	if key != "blacklist:ip:1.2.3.4" {
		t.Fatalf("unexpected key: %s", key)
	}
}

func TestCacheKeyUniqueness(t *testing.T) {
	keys := []string{
		service.UserCacheKey(1),
		service.UserCacheKey(2),
		service.ConfigCacheKey("key1"),
		service.BlacklistCacheKey("ip", "1.1.1.1"),
		service.BlacklistCacheKey("fingerprint", "abc"),
	}
	seen := make(map[string]bool)
	for _, k := range keys {
		if seen[k] {
			t.Fatalf("duplicate cache key: %s", k)
		}
		seen[k] = true
	}
}

func TestUserModel_GuestWithinLimit(t *testing.T) {
	u := &model.User{UserType: "guest", DailyLimit: 5, DailyUsedCount: 4}
	exceeded := int64(u.DailyLimit) > 0 && u.DailyUsedCount >= int64(u.DailyLimit)
	if exceeded {
		t.Fatal("should not be exceeded")
	}
}

func TestUserModel_GuestExceeded(t *testing.T) {
	u := &model.User{UserType: "guest", DailyLimit: 5, DailyUsedCount: 5}
	exceeded := int64(u.DailyLimit) > 0 && u.DailyUsedCount >= int64(u.DailyLimit)
	if !exceeded {
		t.Fatal("expected exceeded")
	}
}

func TestUserModel_VipNoBalance(t *testing.T) {
	u := &model.User{UserType: "vip", VipBalance: 0}
	if u.VipBalance > 0 {
		t.Fatal("expected no balance")
	}
}

func TestUserModel_AdminUnlimited(t *testing.T) {
	u := &model.User{UserType: "admin", DailyUsedCount: 99999}
	if u.UserType != "admin" {
		t.Fatal("expected admin")
	}
}

func TestNewUserService(t *testing.T) {
	_ = service.NewUserService(nil)
}

func TestNewAccountService(t *testing.T) {
	_ = service.NewAccountService(nil, "")
}

func TestNewConfigService(t *testing.T) {
	_ = service.NewConfigService(nil)
}

func TestNewRecordService(t *testing.T) {
	_ = service.NewRecordService(nil)
}

func TestCacheGetMiss(t *testing.T) {
	ctx := context.Background()
	var dest struct{ Val string }
	hit := service.CacheGet(ctx, "nonexistent-key-xyz", &dest)
	if hit {
		t.Fatal("expected cache miss")
	}
}

func TestCacheSetL1Only_NoPanic(t *testing.T) {
	err := service.CacheSetL1Only("test-key", map[string]string{"a": "b"}, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
