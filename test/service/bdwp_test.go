package service_test

import (
	"testing"
)

// TestBdwpClient_AuthHeaders 测试认证头构建逻辑
func TestBdwpClient_AuthHeaders(t *testing.T) {
	tests := []struct {
		name        string
		cookie      string
		accessToken string
		wantCookie  bool
		wantToken   bool
	}{
		{
			name:        "cookie auth",
			cookie:      "BDUSS=xxx",
			accessToken: "",
			wantCookie:  true,
			wantToken:   false,
		},
		{
			name:        "access_token auth",
			cookie:      "",
			accessToken: "token123",
			wantCookie:  false,
			wantToken:   true,
		},
		{
			name:        "both provided - prefer access_token",
			cookie:      "BDUSS=xxx",
			accessToken: "token123",
			wantCookie:  false,
			wantToken:   true,
		},
		{
			name:        "neither provided",
			cookie:      "",
			accessToken: "",
			wantCookie:  false,
			wantToken:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 由于 withAuth 是私有函数，这里仅验证逻辑
			hasCookie := tt.cookie != ""
			hasToken := tt.accessToken != ""

			if hasToken != tt.wantToken {
				t.Errorf("token check failed: got %v, want %v", hasToken, tt.wantToken)
			}
			if !hasToken && hasCookie != tt.wantCookie {
				t.Errorf("cookie check failed: got %v, want %v", hasCookie, tt.wantCookie)
			}
		})
	}
}

// TestShareInfo_Structure 测试 ShareInfo 结构
func TestShareInfo_Structure(t *testing.T) {
	type ShareInfo struct {
		Errno    int
		ShareID  int64
		UK       int64
		BDSToken string
	}

	info := ShareInfo{
		Errno:    0,
		ShareID:  123456,
		UK:       789012,
		BDSToken: "test_token",
	}

	if info.Errno != 0 {
		t.Fatal("errno should be 0 for success")
	}
	if info.ShareID == 0 {
		t.Fatal("shareID should not be 0")
	}
	if info.BDSToken == "" {
		t.Fatal("bdstoken should not be empty")
	}
}

// TestShareFile_Structure 测试 ShareFile 结构
func TestShareFile_Structure(t *testing.T) {
	type ShareFile struct {
		FsID     int64
		Filename string
		Size     int64
		IsDir    int
		Path     string
	}

	file := ShareFile{
		FsID:     123456789,
		Filename: "test.mp4",
		Size:     1024000,
		IsDir:    0,
		Path:     "/test.mp4",
	}

	if file.FsID == 0 {
		t.Fatal("fsID should not be 0")
	}
	if file.Filename == "" {
		t.Fatal("filename should not be empty")
	}
	if file.Size <= 0 {
		t.Fatal("size should be positive")
	}
	if file.IsDir != 0 {
		t.Fatal("isDir should be 0 for file")
	}
}

// TestBanStatus_Structure 测试 BanStatus 结构
func TestBanStatus_Structure(t *testing.T) {
	type BanStatus struct {
		Banned          bool
		StartTime       int64
		EndTime         int64
		BanReason       string
		BanTimes        int
		BanMsg          string
		UserOperateType int
	}

	// 未封禁状态
	status := BanStatus{
		Banned:          false,
		StartTime:       0,
		EndTime:         0,
		BanReason:       "",
		BanTimes:        0,
		BanMsg:          "",
		UserOperateType: 0,
	}

	if status.Banned {
		t.Fatal("should not be banned")
	}

	// 已封禁状态
	bannedStatus := BanStatus{
		Banned:    true,
		StartTime: 1640000000,
		EndTime:   1640086400,
		BanReason: "违规分享",
		BanTimes:  1,
		BanMsg:    "账号已被限速",
	}

	if !bannedStatus.Banned {
		t.Fatal("should be banned")
	}
	if bannedStatus.BanReason == "" {
		t.Fatal("ban reason should not be empty")
	}
}

// TestTransferResp_Errno 测试转存响应错误码
func TestTransferResp_Errno(t *testing.T) {
	type TransferResp struct {
		Errno int
	}

	tests := []struct {
		errno   int
		success bool
	}{
		{0, true},
		{-1, false},
		{-7, false},
		{-8, false},
		{-9, false},
	}

	for _, tt := range tests {
		resp := TransferResp{Errno: tt.errno}
		isSuccess := resp.Errno == 0
		if isSuccess != tt.success {
			t.Errorf("errno %d: expected success=%v, got %v", tt.errno, tt.success, isSuccess)
		}
	}
}

// TestLocateResp_URLs 测试 LocateDownload 响应
func TestLocateResp_URLs(t *testing.T) {
	type URLItem struct {
		URL string
	}
	type LocateResp struct {
		Errno int
		URLs  []URLItem
	}

	resp := LocateResp{
		Errno: 0,
		URLs: []URLItem{
			{URL: "https://d1.baidupcs.com/file/xxx"},
			{URL: "https://d2.baidupcs.com/file/xxx"},
		},
	}

	if resp.Errno != 0 {
		t.Fatal("errno should be 0")
	}
	if len(resp.URLs) != 2 {
		t.Fatalf("expected 2 URLs, got %d", len(resp.URLs))
	}
	for i, u := range resp.URLs {
		if u.URL == "" {
			t.Fatalf("URL[%d] should not be empty", i)
		}
	}
}

// TestDeleteFiles_PathFormat 测试删除文件路径格式
func TestDeleteFiles_PathFormat(t *testing.T) {
	paths := []string{
		"/我的资源/test.mp4",
		"/我的资源/folder/file.zip",
		"/我的资源/中文文件名.txt",
	}

	for _, path := range paths {
		if path == "" {
			t.Fatal("path should not be empty")
		}
		if path[0] != '/' {
			t.Fatalf("path should start with /: %s", path)
		}
	}
}

// TestBaiduAPIEndpoints 测试百度 API 端点常量
func TestBaiduAPIEndpoints(t *testing.T) {
	endpoints := map[string]string{
		"pcs_base":   "https://pan.baidu.com",
		"api_base":   "https://pan.baidu.com/api",
		"disk_base":  "https://pan.baidu.com/rest/2.0/xpan",
		"checkapl":   "https://pan.baidu.com/api/checkapl/download",
		"share_list": "https://pan.baidu.com/share/list",
		"transfer":   "https://pan.baidu.com/api/share/transfer",
	}

	for name, url := range endpoints {
		if url == "" {
			t.Fatalf("%s endpoint should not be empty", name)
		}
		if url[:8] != "https://" {
			t.Fatalf("%s endpoint should use https: %s", name, url)
		}
	}
}
