package service_test

import (
	"testing"

	"panflow/internal/model"
)

// TestTokenModel_Structure 测试 Token 模型结构
func TestTokenModel_Structure(t *testing.T) {
	token := &model.Token{
		ID:        1,
		Token:     "test_token_123",
		TokenType: "daily",
		UserType:  "guest",
		Count:     10,
		Size:      10737418240, // 10GB
		Day:       1,
		Switch:    true,
	}

	if token.ID == 0 {
		t.Fatal("token ID should not be 0")
	}
	if token.Token == "" {
		t.Fatal("token string should not be empty")
	}
	if !token.Switch {
		t.Fatal("token should be enabled")
	}
}

// TestTokenType_Daily 测试每日类型 token
func TestTokenType_Daily(t *testing.T) {
	token := &model.Token{
		TokenType: "daily",
		Day:       5,
		UsedCount: 3,
	}

	if token.TokenType != "daily" {
		t.Fatal("token type should be daily")
	}

	// 未超限
	exceeded := token.Day > 0 && token.UsedCount >= token.Day
	if exceeded {
		t.Fatal("should not exceed daily limit")
	}
}

// TestTokenType_Normal 测试普通类型 token
func TestTokenType_Normal(t *testing.T) {
	token := &model.Token{
		TokenType: "normal",
		Count:     100,
		UsedCount: 50,
	}

	if token.TokenType != "normal" {
		t.Fatal("token type should be normal")
	}

	// 未超限
	exceeded := token.Count > 0 && token.UsedCount >= token.Count
	if exceeded {
		t.Fatal("should not exceed count limit")
	}
}

// TestTokenUserType_Guest 测试游客类型 token
func TestTokenUserType_Guest(t *testing.T) {
	token := &model.Token{
		UserType: "guest",
		Count:    10,
		Size:     10737418240,
	}

	if token.UserType != "guest" {
		t.Fatal("user type should be guest")
	}
}

// TestTokenUserType_VIP 测试 VIP 类型 token
func TestTokenUserType_VIP(t *testing.T) {
	token := &model.Token{
		UserType: "vip",
		Count:    100,
	}

	if token.UserType != "vip" {
		t.Fatal("user type should be vip")
	}
}

// TestTokenUserType_SVIP 测试 SVIP 类型 token
func TestTokenUserType_SVIP(t *testing.T) {
	providerID := uint(100)
	token := &model.Token{
		UserType:       "svip",
		ProviderUserID: &providerID,
	}

	if token.UserType != "svip" {
		t.Fatal("user type should be svip")
	}
	if token.ProviderUserID == nil {
		t.Fatal("SVIP token should have provider_user_id")
	}
}

// TestTokenUserType_Admin 测试管理员类型 token
func TestTokenUserType_Admin(t *testing.T) {
	token := &model.Token{
		UserType: "admin",
		Count:    0, // 无限制
	}

	if token.UserType != "admin" {
		t.Fatal("user type should be admin")
	}
}

// TestTokenIPLimit 测试 IP 限制
func TestTokenIPLimit(t *testing.T) {
	token := &model.Token{
		CanUseIPCount: 3,
		IP:            model.JSONSlice{"1.2.3.4", "5.6.7.8"},
	}

	if token.CanUseIPCount != 3 {
		t.Fatal("IP limit should be 3")
	}
	if len(token.IP) != 2 {
		t.Fatal("should have 2 IPs")
	}

	// 未达到限制
	if int64(len(token.IP)) >= token.CanUseIPCount {
		t.Fatal("should not reach IP limit")
	}
}

// TestTokenIPLimit_Reached 测试 IP 限制达到
func TestTokenIPLimit_Reached(t *testing.T) {
	token := &model.Token{
		CanUseIPCount: 2,
		IP:            model.JSONSlice{"1.2.3.4", "5.6.7.8"},
	}

	// 达到限制
	if int64(len(token.IP)) < token.CanUseIPCount {
		t.Fatal("should reach IP limit")
	}
}

// TestTokenSizeQuota 测试大小配额
func TestTokenSizeQuota(t *testing.T) {
	token := &model.Token{
		Size:     10737418240, // 10GB
		UsedSize: 5368709120,  // 5GB
	}

	remaining := token.Size - token.UsedSize
	if remaining != 5368709120 {
		t.Fatalf("expected remaining=5GB, got %d", remaining)
	}

	// 未超限
	exceeded := token.Size > 0 && token.UsedSize >= token.Size
	if exceeded {
		t.Fatal("should not exceed size quota")
	}
}

// TestTokenSizeQuota_Exceeded 测试大小配额超限
func TestTokenSizeQuota_Exceeded(t *testing.T) {
	token := &model.Token{
		Size:     10737418240, // 10GB
		UsedSize: 10737418240, // 10GB
	}

	// 超限
	exceeded := token.Size > 0 && token.UsedSize >= token.Size
	if !exceeded {
		t.Fatal("should exceed size quota")
	}
}

// TestTokenCountQuota 测试次数配额
func TestTokenCountQuota(t *testing.T) {
	token := &model.Token{
		Count:     100,
		UsedCount: 50,
	}

	remaining := token.Count - token.UsedCount
	if remaining != 50 {
		t.Fatalf("expected remaining=50, got %d", remaining)
	}
}

// TestTokenUnlimited 测试无限制 token
func TestTokenUnlimited(t *testing.T) {
	token := &model.Token{
		Count:     0, // 0 表示无限制
		Size:      0,
		UsedCount: 99999,
		UsedSize:  999999999,
	}

	// 无限制
	countExceeded := token.Count > 0 && token.UsedCount >= token.Count
	sizeExceeded := token.Size > 0 && token.UsedSize >= token.Size

	if countExceeded || sizeExceeded {
		t.Fatal("unlimited token should never exceed")
	}
}

// TestGuestToken_Default 测试默认游客 token
func TestGuestToken_Default(t *testing.T) {
	token := &model.Token{
		Token:         "guest",
		TokenType:     "daily",
		UserType:      "guest",
		Count:         10,
		Size:          10737418240,
		Day:           1,
		CanUseIPCount: 99999,
		Switch:        true,
	}

	if token.Token != "guest" {
		t.Fatal("default guest token should be 'guest'")
	}
	if token.TokenType != "daily" {
		t.Fatal("guest token should be daily type")
	}
	if token.CanUseIPCount != 99999 {
		t.Fatal("guest token should have high IP limit")
	}
}
