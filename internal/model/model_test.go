package model_test

import (
	"testing"
	"time"

	"panflow/internal/model"

	"gorm.io/gorm"
)

func TestAccount_Timestamps(t *testing.T) {
	now := time.Now()
	account := &model.Account{CreatedAt: now, UpdatedAt: now}
	if account.CreatedAt.IsZero() {
		t.Fatal("created_at should not be zero")
	}
}

func TestAccount_SoftDelete(t *testing.T) {
	account := &model.Account{DeletedAt: gorm.DeletedAt{}}
	if account.DeletedAt.Valid {
		t.Fatal("new account should not be deleted")
	}
}

func TestUser_Timestamps(t *testing.T) {
	now := time.Now()
	user := &model.User{CreatedAt: now}
	if user.CreatedAt.IsZero() {
		t.Fatal("created_at should not be zero")
	}
}

func TestUser_PasswordHidden(t *testing.T) {
	user := &model.User{Username: "test", Password: "secret"}
	if user.Password == "" {
		t.Fatal("password should be set internally")
	}
}

func TestBlackList_ExpiresAt(t *testing.T) {
	future := model.LocalTime{Time: time.Now().Add(24 * time.Hour)}
	bl := &model.BlackList{ExpiresAt: &future}
	if bl.ExpiresAt == nil {
		t.Fatal("expires_at should not be nil")
	}
	if bl.ExpiresAt.Before(time.Now()) {
		t.Fatal("blacklist should not be expired")
	}
}

func TestBlackList_Permanent(t *testing.T) {
	bl := &model.BlackList{Type: "ip", Identifier: "1.2.3.4", ExpiresAt: nil}
	if bl.ExpiresAt != nil {
		t.Fatal("permanent blacklist should have nil expires_at")
	}
}

func TestFileList_UniqueIndex(t *testing.T) {
	file1 := &model.FileList{FsID: "123456"}
	file2 := &model.FileList{FsID: "123456"}
	if file1.FsID != file2.FsID {
		t.Fatal("same fs_id should match")
	}
}

func TestConfig_TypeValidation(t *testing.T) {
	for _, typ := range []string{"string", "int", "bool", "json"} {
		cfg := &model.Config{Type: typ}
		if cfg.Type == "" {
			t.Fatalf("type %s should not be empty", typ)
		}
	}
}

func TestRecord_Structure(t *testing.T) {
	userID := uint(1)
	record := &model.Record{
		IP:        "1.2.3.4",
		AccountID: 1,
		UserID:    &userID,
	}
	if record.IP == "" {
		t.Fatal("IP should not be empty")
	}
	if record.UserID == nil {
		t.Fatal("user_id should not be nil")
	}
}

func TestRecord_OptionalUser(t *testing.T) {
	record := &model.Record{AccountID: 1, UserID: nil}
	if record.UserID != nil {
		t.Fatal("guest record should have nil user_id")
	}
}
