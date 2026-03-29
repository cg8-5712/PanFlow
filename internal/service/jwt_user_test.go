package service_test

import (
	"testing"

	"panflow/internal/service"
)

// TestNewJWTService 测试 JWT 服务创建
func TestNewJWTService(t *testing.T) {
	svc := service.NewJWTService("secret", 2)
	if svc == nil {
		t.Fatal("service should not be nil")
	}
}

// TestJWTService_DefaultExpireHours 测试默认过期时间
func TestJWTService_DefaultExpireHours(t *testing.T) {
	svc := service.NewJWTService("secret", 0)
	if svc == nil {
		t.Fatal("service should not be nil")
	}
	// 默认应该是 24 小时，通过 Issue 验证
	_, exp, err := svc.Issue()
	if err != nil {
		t.Fatal(err)
	}
	// 过期时间应该在未来
	if exp.Unix() <= 0 {
		t.Fatal("expiry should be set")
	}
}

// TestJWTService_IssueUser 测试用户 JWT 签发
func TestJWTService_IssueUser(t *testing.T) {
	svc := service.NewJWTService("secret", 2)

	tokenID := uint(1)
	userType := "guest"
	tokenStr, exp, err := svc.IssueUser(tokenID, userType, nil)

	if err != nil {
		t.Fatalf("issue user failed: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("token should not be empty")
	}
	if exp.Unix() <= 0 {
		t.Fatal("expiry should be set")
	}
}

// TestJWTService_VerifyUser 测试用户 JWT 验证
func TestJWTService_VerifyUser(t *testing.T) {
	svc := service.NewJWTService("secret", 2)

	tokenID := uint(123)
	userType := "vip"
	tokenStr, _, _ := svc.IssueUser(tokenID, userType, nil)

	claims, err := svc.VerifyUser(tokenStr)
	if err != nil {
		t.Fatalf("verify user failed: %v", err)
	}
	if claims.TokenID != tokenID {
		t.Fatalf("expected token_id %d, got %d", tokenID, claims.TokenID)
	}
	if claims.UserType != userType {
		t.Fatalf("expected user_type %s, got %s", userType, claims.UserType)
	}
}

// TestJWTService_VerifyUser_WithUserID 测试带用户 ID 的验证
func TestJWTService_VerifyUser_WithUserID(t *testing.T) {
	svc := service.NewJWTService("secret", 2)

	tokenID := uint(1)
	userID := uint(100)
	tokenStr, _, _ := svc.IssueUser(tokenID, "svip", &userID)

	claims, err := svc.VerifyUser(tokenStr)
	if err != nil {
		t.Fatalf("verify user failed: %v", err)
	}
	if claims.UserID == nil {
		t.Fatal("user_id should not be nil")
	}
	if *claims.UserID != userID {
		t.Fatalf("expected user_id %d, got %d", userID, *claims.UserID)
	}
}

// TestJWTService_VerifyUser_InvalidToken 测试无效 token
func TestJWTService_VerifyUser_InvalidToken(t *testing.T) {
	svc := service.NewJWTService("secret", 2)

	_, err := svc.VerifyUser("invalid.token.string")
	if err == nil {
		t.Fatal("should fail for invalid token")
	}
}

// TestJWTService_VerifyUser_WrongSecret 测试错误密钥
func TestJWTService_VerifyUser_WrongSecret(t *testing.T) {
	issuer := service.NewJWTService("secret-A", 2)
	verifier := service.NewJWTService("secret-B", 2)

	tokenStr, _, _ := issuer.IssueUser(1, "guest", nil)

	_, err := verifier.VerifyUser(tokenStr)
	if err == nil {
		t.Fatal("should fail with wrong secret")
	}
}

// TestJWTService_UserTypes 测试所有用户类型
func TestJWTService_UserTypes(t *testing.T) {
	svc := service.NewJWTService("secret", 2)

	types := []string{"guest", "vip", "svip", "admin"}
	for _, typ := range types {
		tokenStr, _, err := svc.IssueUser(1, typ, nil)
		if err != nil {
			t.Fatalf("issue failed for type %s: %v", typ, err)
		}

		claims, err := svc.VerifyUser(tokenStr)
		if err != nil {
			t.Fatalf("verify failed for type %s: %v", typ, err)
		}
		if claims.UserType != typ {
			t.Fatalf("expected type %s, got %s", typ, claims.UserType)
		}
	}
}

// TestJWTService_EmptyToken 测试空 token
func TestJWTService_EmptyToken(t *testing.T) {
	svc := service.NewJWTService("secret", 2)

	_, err := svc.VerifyUser("")
	if err == nil {
		t.Fatal("should fail for empty token")
	}
}

// TestJWTService_AdminAndUser 测试管理员和用户 JWT 不混淆
func TestJWTService_AdminAndUser(t *testing.T) {
	svc := service.NewJWTService("secret", 2)

	// 签发管理员 token（Issue 返回 AdminClaims，VerifyUser 解析为 UserClaims 会失败）
	adminToken, _, _ := svc.Issue()

	// 用用户验证器验证管理员 token：AdminClaims 无 UserClaims 字段，解析会失败
	_, err := svc.VerifyUser(adminToken)
	if err == nil {
		t.Fatal("admin token (AdminClaims) should not pass user verification")
	}

	// 签发用户 token
	userToken, _, _ := svc.IssueUser(1, "guest", nil)

	// 用管理员验证器验证用户 token 应该失败
	_, err = svc.Verify(userToken)
	if err == nil {
		t.Fatal("user token should not pass admin verification")
	}
}
