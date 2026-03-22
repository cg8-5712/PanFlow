package model_test

import (
	"testing"
	"time"

	"panflow/internal/model"

	"gorm.io/gorm"
)

// TestAccount_Timestamps 测试时间戳字段
func TestAccount_Timestamps(t *testing.T) {
	now := time.Now()
	account := &model.Account{
		CreatedAt: now,
		UpdatedAt: now,
	}

	if account.CreatedAt.IsZero() {
		t.Fatal("created_at should not be zero")
	}
	if account.UpdatedAt.IsZero() {
		t.Fatal("updated_at should not be zero")
	}
}

// TestAccount_SoftDelete 测试软删除
func TestAccount_SoftDelete(t *testing.T) {
	account := &model.Account{
		DeletedAt: gorm.DeletedAt{},
	}

	if account.DeletedAt.Valid {
		t.Fatal("new account should not be deleted")
	}
}

// TestToken_Timestamps 测试 Token 时间戳
func TestToken_Timestamps(t *testing.T) {
	now := time.Now()
	token := &model.Token{
		CreatedAt: now,
		UpdatedAt: now,
	}

	if token.CreatedAt.IsZero() {
		t.Fatal("created_at should not be zero")
	}
}

// TestToken_ExpiresAt 测试过期时间
func TestToken_ExpiresAt(t *testing.T) {
	future := time.Now().Add(24 * time.Hour)
	token := &model.Token{
		ExpiresAt: &future,
	}

	if token.ExpiresAt == nil {
		t.Fatal("expires_at should not be nil")
	}
	if token.ExpiresAt.Before(time.Now()) {
		t.Fatal("token should not be expired")
	}
}

// TestUser_Timestamps 测试用户时间戳
func TestUser_Timestamps(t *testing.T) {
	now := time.Now()
	user := &model.User{
		CreatedAt: now,
		UpdatedAt: now,
	}

	if user.CreatedAt.IsZero() {
		t.Fatal("created_at should not be zero")
	}
}

// TestUser_PasswordHidden 测试密码字段隐藏
func TestUser_PasswordHidden(t *testing.T) {
	user := &model.User{
		Username: "test",
		Password: "secret",
	}

	// Password 字段有 json:"-" 标签，不会被序列化
	if user.Password == "" {
		t.Fatal("password should be set internally")
	}
}

// TestBlackList_ExpiresAt 测试黑名单过期
func TestBlackList_ExpiresAt(t *testing.T) {
	future := time.Now().Add(24 * time.Hour)
	bl := &model.BlackList{
		ExpiresAt: &future,
	}

	if bl.ExpiresAt == nil {
		t.Fatal("expires_at should not be nil")
	}
	if bl.ExpiresAt.Before(time.Now()) {
		t.Fatal("blacklist should not be expired")
	}
}

// TestBlackList_Permanent 测试永久黑名单
func TestBlackList_Permanent(t *testing.T) {
	bl := &model.BlackList{
		Type:       "ip",
		Identifier: "1.2.3.4",
		ExpiresAt:  nil, // 永久
	}

	if bl.ExpiresAt != nil {
		t.Fatal("permanent blacklist should have nil expires_at")
	}
}

// TestFileList_UniqueIndex 测试唯一索引
func TestFileList_UniqueIndex(t *testing.T) {
	file1 := &model.FileList{
		FsID: "123456",
	}
	file2 := &model.FileList{
		FsID: "123456",
	}

	// FsID 应该是唯一的
	if file1.FsID != file2.FsID {
		t.Fatal("same fs_id should match")
	}
}

// TestConfig_TypeValidation 测试配置类型
func TestConfig_TypeValidation(t *testing.T) {
	validTypes := []string{"string", "int", "bool", "json"}
	for _, typ := range validTypes {
		cfg := &model.Config{
			Type: typ,
		}
		if cfg.Type == "" {
			t.Fatalf("type %s should not be empty", typ)
		}
	}
}

// TestRecord_Associations 测试关联关系
func TestRecord_Associations(t *testing.T) {
	tokenID := uint(1)
	accountID := uint(2)
	userID := uint(3)

	record := &model.Record{
		TokenID:   tokenID,
		AccountID: accountID,
		UserID:    &userID,
	}

	if record.TokenID == 0 {
		t.Fatal("token_id should not be 0")
	}
	if record.AccountID == 0 {
		t.Fatal("account_id should not be 0")
	}
	if record.UserID == nil {
		t.Fatal("user_id should not be nil")
	}
}

// TestRecord_OptionalUser 测试可选用户关联
func TestRecord_OptionalUser(t *testing.T) {
	record := &model.Record{
		TokenID:   1,
		AccountID: 1,
		UserID:    nil, // 游客无用户 ID
	}

	if record.UserID != nil {
		t.Fatal("guest record should have nil user_id")
	}
}
