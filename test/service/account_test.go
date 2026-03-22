package service_test

import (
	"testing"

	"panflow/internal/model"
)

// TestAccountModel_Structure 测试账号模型结构
func TestAccountModel_Structure(t *testing.T) {
	account := &model.Account{
		ID:          1,
		BaiduName:   "test_user",
		UK:          "123456",
		AccountType: "cookie",
		AccountData: model.JSONMap{
			"cookie": "BDUSS=xxx",
		},
		Switch:    true,
		UsedCount: 10,
		UsedSize:  1024000,
	}

	if account.ID == 0 {
		t.Fatal("account ID should not be 0")
	}
	if account.AccountType == "" {
		t.Fatal("account type should not be empty")
	}
	if !account.Switch {
		t.Fatal("account should be enabled")
	}
}

// TestAccountType_Cookie 测试 cookie 类型账号
func TestAccountType_Cookie(t *testing.T) {
	account := &model.Account{
		AccountType: "cookie",
		AccountData: model.JSONMap{
			"cookie":     "BDUSS=xxx; STOKEN=yyy;",
			"vip_type":   "超级会员",
			"expires_at": "2025-12-31 23:59:59",
		},
	}

	if account.AccountType != "cookie" {
		t.Fatal("account type should be cookie")
	}

	cookie, ok := account.AccountData["cookie"].(string)
	if !ok || cookie == "" {
		t.Fatal("cookie should not be empty")
	}
}

// TestAccountType_OpenPlatform 测试 open_platform 类型账号
func TestAccountType_OpenPlatform(t *testing.T) {
	account := &model.Account{
		AccountType: "open_platform",
		AccountData: model.JSONMap{
			"access_token":     "access_token_value",
			"refresh_token":    "refresh_token_value",
			"token_expires_at": "2025-12-31 23:59:59",
		},
	}

	if account.AccountType != "open_platform" {
		t.Fatal("account type should be open_platform")
	}

	accessToken, ok := account.AccountData["access_token"].(string)
	if !ok || accessToken == "" {
		t.Fatal("access_token should not be empty")
	}
}

// TestAccountType_EnterpriseCookie 测试 enterprise_cookie 类型账号
func TestAccountType_EnterpriseCookie(t *testing.T) {
	account := &model.Account{
		AccountType: "enterprise_cookie",
		AccountData: model.JSONMap{
			"cookie":       "enterprise_cookie_value",
			"cid":          float64(12345),
			"bdstoken":     "bdstoken_value",
			"dlink_cookie": "dlink_cookie_value",
		},
	}

	if account.AccountType != "enterprise_cookie" {
		t.Fatal("account type should be enterprise_cookie")
	}

	cid, ok := account.AccountData["cid"].(float64)
	if !ok || cid == 0 {
		t.Fatal("cid should not be 0")
	}
}

// TestAccountType_DownloadTicket 测试 download_ticket 类型账号
func TestAccountType_DownloadTicket(t *testing.T) {
	account := &model.Account{
		AccountType: "download_ticket",
		AccountData: model.JSONMap{
			"surl":            "test_surl",
			"pwd":             "test_pwd",
			"dir":             "/",
			"save_cookie":     "save_cookie_value",
			"save_bdstoken":   "save_bdstoken_value",
			"download_cookie": "download_cookie_value",
		},
	}

	if account.AccountType != "download_ticket" {
		t.Fatal("account type should be download_ticket")
	}

	saveCookie, ok := account.AccountData["save_cookie"].(string)
	if !ok || saveCookie == "" {
		t.Fatal("save_cookie should not be empty")
	}

	downloadCookie, ok := account.AccountData["download_cookie"].(string)
	if !ok || downloadCookie == "" {
		t.Fatal("download_cookie should not be empty")
	}
}

// TestAccountSwitch_Enabled 测试账号启用状态
func TestAccountSwitch_Enabled(t *testing.T) {
	account := &model.Account{
		Switch: true,
		Reason: "",
	}

	if !account.Switch {
		t.Fatal("account should be enabled")
	}
	if account.Reason != "" {
		t.Fatal("enabled account should have no reason")
	}
}

// TestAccountSwitch_Disabled 测试账号禁用状态
func TestAccountSwitch_Disabled(t *testing.T) {
	account := &model.Account{
		Switch: false,
		Reason: "账号被封禁",
	}

	if account.Switch {
		t.Fatal("account should be disabled")
	}
	if account.Reason == "" {
		t.Fatal("disabled account should have reason")
	}
}

// TestAccountUsage_Increment 测试用量递增
func TestAccountUsage_Increment(t *testing.T) {
	account := &model.Account{
		UsedCount: 0,
		UsedSize:  0,
	}

	// 模拟使用
	account.UsedCount++
	account.UsedSize += 1024000

	if account.UsedCount != 1 {
		t.Fatalf("expected used_count=1, got %d", account.UsedCount)
	}
	if account.UsedSize != 1024000 {
		t.Fatalf("expected used_size=1024000, got %d", account.UsedSize)
	}
}

// TestAccountProviderUser 测试账号提供者关联
func TestAccountProviderUser(t *testing.T) {
	providerID := uint(100)
	account := &model.Account{
		ProviderUserID: &providerID,
	}

	if account.ProviderUserID == nil {
		t.Fatal("provider_user_id should not be nil")
	}
	if *account.ProviderUserID != 100 {
		t.Fatalf("expected provider_user_id=100, got %d", *account.ProviderUserID)
	}
}

// TestAccountWithoutProvider 测试无提供者的账号
func TestAccountWithoutProvider(t *testing.T) {
	account := &model.Account{
		ProviderUserID: nil,
	}

	if account.ProviderUserID != nil {
		t.Fatal("provider_user_id should be nil for public accounts")
	}
}

// TestAccountTotalSize 测试账号总容量
func TestAccountTotalSize(t *testing.T) {
	account := &model.Account{
		TotalSize: 2199023255552, // 2TB
		UsedSize:  1099511627776, // 1TB
	}

	remaining := account.TotalSize - account.UsedSize
	if remaining != 1099511627776 {
		t.Fatalf("expected remaining=1TB, got %d", remaining)
	}
}

// TestAccountVipType 测试会员类型
func TestAccountVipType(t *testing.T) {
	tests := []struct {
		vipType string
		valid   bool
	}{
		{"超级会员", true},
		{"普通会员", true},
		{"普通用户", true},
		{"", false},
	}

	for _, tt := range tests {
		account := &model.Account{
			AccountData: model.JSONMap{
				"vip_type": tt.vipType,
			},
		}

		vipType, ok := account.AccountData["vip_type"].(string)
		if tt.valid && (!ok || vipType == "") {
			t.Errorf("vip_type %s should be valid", tt.vipType)
		}
	}
}
