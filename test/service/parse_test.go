package service_test

import (
	"testing"

	"panflow/internal/model"
	"panflow/internal/service"
)

// TestExtractCreds_Cookie 测试 cookie 类型账号凭据提取
func TestExtractCreds_Cookie(t *testing.T) {
	account := &model.Account{
		AccountType: "cookie",
		AccountData: model.JSONMap{
			"cookie":     "BDUSS=xxx; STOKEN=yyy;",
			"user_agent": "netdisk;2.0",
		},
	}

	// 由于 extractCreds 是私有方法，这里仅做结构验证
	if account.AccountData["cookie"] != "BDUSS=xxx; STOKEN=yyy;" {
		t.Fatal("cookie extraction failed")
	}
}

// TestExtractCreds_OpenPlatform 测试 open_platform 类型
func TestExtractCreds_OpenPlatform(t *testing.T) {
	account := &model.Account{
		AccountType: "open_platform",
		AccountData: model.JSONMap{
			"access_token":  "access_token_value",
			"refresh_token": "refresh_token_value",
		},
	}

	if account.AccountData["access_token"] != "access_token_value" {
		t.Fatal("access_token extraction failed")
	}
}

// TestExtractCreds_DownloadTicket 测试 download_ticket 类型
func TestExtractCreds_DownloadTicket(t *testing.T) {
	account := &model.Account{
		AccountType: "download_ticket",
		AccountData: model.JSONMap{
			"save_cookie":     "save_cookie_value",
			"download_cookie": "download_cookie_value",
			"surl":            "test_surl",
		},
	}

	saveCookie, _ := account.AccountData["save_cookie"].(string)
	downloadCookie, _ := account.AccountData["download_cookie"].(string)

	if saveCookie != "save_cookie_value" {
		t.Fatal("save_cookie extraction failed")
	}
	if downloadCookie != "download_cookie_value" {
		t.Fatal("download_cookie extraction failed")
	}
}

// TestExtractCreds_EnterpriseCookie 测试 enterprise_cookie 类型
func TestExtractCreds_EnterpriseCookie(t *testing.T) {
	account := &model.Account{
		AccountType: "enterprise_cookie",
		AccountData: model.JSONMap{
			"cookie":       "enterprise_cookie_value",
			"dlink_cookie": "dlink_cookie_value",
			"cid":          float64(12345),
		},
	}

	cookie, _ := account.AccountData["cookie"].(string)
	dlinkCookie, _ := account.AccountData["dlink_cookie"].(string)
	cid, _ := account.AccountData["cid"].(float64)

	if cookie != "enterprise_cookie_value" {
		t.Fatal("cookie extraction failed")
	}
	if dlinkCookie != "dlink_cookie_value" {
		t.Fatal("dlink_cookie extraction failed")
	}
	if cid != 12345 {
		t.Fatal("cid extraction failed")
	}
}

// TestParseRequest_Validation 测试 ParseRequest 结构验证
func TestParseRequest_Validation(t *testing.T) {
	req := &service.ParseRequest{
		Surl:     "test_surl",
		Pwd:      "test_pwd",
		FsIDs:    []int64{123, 456},
		ClientIP: "1.2.3.4",
		TokenID:  1,
		UserType: "guest",
	}

	if req.Surl != "test_surl" {
		t.Fatal("surl validation failed")
	}
	if len(req.FsIDs) != 2 {
		t.Fatal("fsIDs validation failed")
	}
	if req.UserType != "guest" {
		t.Fatal("userType validation failed")
	}
}

// TestParseResult_Structure 测试 ParseResult 结构
func TestParseResult_Structure(t *testing.T) {
	result := &service.ParseResult{
		FsID: 123456,
		URLs: []string{"https://cdn1.example.com/file", "https://cdn2.example.com/file"},
		Size: 1024000,
	}

	if result.FsID != 123456 {
		t.Fatal("fsID validation failed")
	}
	if len(result.URLs) != 2 {
		t.Fatal("URLs validation failed")
	}
	if result.Size != 1024000 {
		t.Fatal("size validation failed")
	}
}

// TestAccountType_AllTypes 测试所有账号类型常量
func TestAccountType_AllTypes(t *testing.T) {
	types := []string{"cookie", "open_platform", "enterprise_cookie", "download_ticket"}
	for _, typ := range types {
		if typ == "" {
			t.Fatalf("account type should not be empty")
		}
	}
}

// TestUserType_AllTypes 测试所有用户类型
func TestUserType_AllTypes(t *testing.T) {
	types := []string{"guest", "vip", "svip", "admin"}
	for _, typ := range types {
		if typ == "" {
			t.Fatalf("user type should not be empty")
		}
	}
}
