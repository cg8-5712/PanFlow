package service_test

import (
	"testing"
	"time"

	"panflow/internal/service"
)

func TestJWTIssueAndVerify(t *testing.T) {
	svc := service.NewJWTService("my-test-secret", 1)

	tokenStr, exp, err := svc.Issue()
	if err != nil {
		t.Fatalf("Issue failed: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("expected non-empty token")
	}
	if exp.Before(time.Now()) {
		t.Fatal("expiry should be in the future")
	}

	claims, err := svc.Verify(tokenStr)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
	if claims.Role != "admin" {
		t.Fatalf("expected role=admin, got %s", claims.Role)
	}
}

func TestJWTVerify_InvalidToken(t *testing.T) {
	svc := service.NewJWTService("secret", 1)
	_, err := svc.Verify("this.is.not.valid")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestJWTVerify_WrongSecret(t *testing.T) {
	issuer := service.NewJWTService("secret-A", 1)
	verifier := service.NewJWTService("secret-B", 1)

	tokenStr, _, _ := issuer.Issue()
	_, err := verifier.Verify(tokenStr)
	if err == nil {
		t.Fatal("expected error when verifying with wrong secret")
	}
}

func TestJWTDefaultExpireHours(t *testing.T) {
	// expireHours=0 should default to 24
	svc := service.NewJWTService("secret", 0)
	_, exp, err := svc.Issue()
	if err != nil {
		t.Fatal(err)
	}
	// Should expire ~24h from now (allow 1min tolerance)
	diff := time.Until(exp)
	if diff < 23*time.Hour || diff > 25*time.Hour {
		t.Fatalf("expected ~24h expiry, got %v", diff)
	}
}

func TestJWTEmptyStringRejected(t *testing.T) {
	svc := service.NewJWTService("secret", 1)
	_, err := svc.Verify("")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}
