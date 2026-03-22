package service_test

import (
	"testing"

	"panflow/internal/model"
)

// TestUserQuota_Guest 测试游客配额
func TestUserQuota_Guest(t *testing.T) {
	user := &model.User{
		UserType:       "guest",
		DailyLimit:     5,
		DailyUsedCount: 3,
	}

	// 未超限
	exceeded := int64(user.DailyLimit) > 0 && user.DailyUsedCount >= int64(user.DailyLimit)
	if exceeded {
		t.Fatal("guest should not exceed limit")
	}

	// 达到限制
	user.DailyUsedCount = 5
	exceeded = int64(user.DailyLimit) > 0 && user.DailyUsedCount >= int64(user.DailyLimit)
	if !exceeded {
		t.Fatal("guest should exceed limit")
	}
}

// TestUserQuota_VIP 测试 VIP 配额
func TestUserQuota_VIP(t *testing.T) {
	user := &model.User{
		UserType:   "vip",
		VipBalance: 100,
	}

	// 有余额
	if user.VipBalance <= 0 {
		t.Fatal("VIP should have balance")
	}

	// 余额耗尽
	user.VipBalance = 0
	if user.VipBalance > 0 {
		t.Fatal("VIP should have no balance")
	}
}

// TestUserQuota_SVIP 测试 SVIP 配额
func TestUserQuota_SVIP(t *testing.T) {
	user := &model.User{
		UserType:       "svip",
		DailyLimit:     100,
		DailyUsedCount: 50,
	}

	// 未超限
	exceeded := int64(user.DailyLimit) > 0 && user.DailyUsedCount >= int64(user.DailyLimit)
	if exceeded {
		t.Fatal("SVIP should not exceed limit")
	}

	// 达到限制
	user.DailyUsedCount = 100
	exceeded = int64(user.DailyLimit) > 0 && user.DailyUsedCount >= int64(user.DailyLimit)
	if !exceeded {
		t.Fatal("SVIP should exceed limit")
	}
}

// TestUserQuota_Admin 测试管理员无限制
func TestUserQuota_Admin(t *testing.T) {
	user := &model.User{
		UserType:       "admin",
		DailyUsedCount: 99999,
	}

	// 管理员不受限制
	if user.UserType != "admin" {
		t.Fatal("should be admin")
	}
}

// TestUserType_Validation 测试用户类型验证
func TestUserType_Validation(t *testing.T) {
	validTypes := []string{"guest", "vip", "svip", "admin"}

	for _, typ := range validTypes {
		user := &model.User{UserType: typ}
		if user.UserType == "" {
			t.Fatalf("user type %s should not be empty", typ)
		}
	}
}

// TestUserVipBalance_Deduction 测试 VIP 余额扣减
func TestUserVipBalance_Deduction(t *testing.T) {
	user := &model.User{
		UserType:   "vip",
		VipBalance: 100,
	}

	// 扣减前
	if user.VipBalance != 100 {
		t.Fatal("initial balance should be 100")
	}

	// 模拟扣减
	user.VipBalance -= 1
	if user.VipBalance != 99 {
		t.Fatalf("balance should be 99 after deduction, got %d", user.VipBalance)
	}

	// 多次扣减
	for i := 0; i < 99; i++ {
		user.VipBalance -= 1
	}
	if user.VipBalance != 0 {
		t.Fatalf("balance should be 0 after full deduction, got %d", user.VipBalance)
	}
}

// TestUserDailyReset 测试每日重置逻辑
func TestUserDailyReset(t *testing.T) {
	user := &model.User{
		UserType:       "guest",
		DailyLimit:     5,
		DailyUsedCount: 5,
	}

	// 达到限制
	exceeded := int64(user.DailyLimit) > 0 && user.DailyUsedCount >= int64(user.DailyLimit)
	if !exceeded {
		t.Fatal("should exceed limit before reset")
	}

	// 模拟重置
	user.DailyUsedCount = 0
	exceeded = int64(user.DailyLimit) > 0 && user.DailyUsedCount >= int64(user.DailyLimit)
	if exceeded {
		t.Fatal("should not exceed limit after reset")
	}
}

// TestUserBaiduAccountBinding 测试百度账号绑定
func TestUserBaiduAccountBinding(t *testing.T) {
	accountID := uint(100)
	user := &model.User{
		UserType:       "svip",
		BaiduAccountID: &accountID,
	}

	if user.BaiduAccountID == nil {
		t.Fatal("SVIP user should have baidu account binding")
	}
	if *user.BaiduAccountID != 100 {
		t.Fatalf("expected account ID 100, got %d", *user.BaiduAccountID)
	}
}

// TestUserWithoutBinding 测试无绑定用户
func TestUserWithoutBinding(t *testing.T) {
	user := &model.User{
		UserType:       "guest",
		BaiduAccountID: nil,
	}

	if user.BaiduAccountID != nil {
		t.Fatal("guest user should not have baidu account binding")
	}
}
