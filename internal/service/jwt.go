package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrTokenInvalid = errors.New("invalid or expired token")

type JWTService struct {
	secret      []byte
	expireHours int
}

func NewJWTService(secret string, expireHours int) *JWTService {
	if expireHours <= 0 {
		expireHours = 24
	}
	return &JWTService{secret: []byte(secret), expireHours: expireHours}
}

// ── Admin claims ──────────────────────────────────────────────────────────────

type AdminClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// Issue signs a new admin JWT
func (s *JWTService) Issue() (string, time.Time, error) {
	exp := time.Now().Add(time.Duration(s.expireHours) * time.Hour)
	claims := AdminClaims{
		Role: "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	return signed, exp, err
}

// Verify parses and validates an admin JWT
func (s *JWTService) Verify(tokenStr string) (*AdminClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AdminClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		return s.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrTokenInvalid
	}
	claims, ok := token.Claims.(*AdminClaims)
	if !ok || claims.Role != "admin" {
		return nil, ErrTokenInvalid
	}
	return claims, nil
}

// ── User claims ───────────────────────────────────────────────────────────────

// UserClaims is embedded in the JWT issued to regular users after login.
// UserType mirrors the token's user_type: guest | vip | svip | admin.
type UserClaims struct {
	TokenID  uint   `json:"token_id"`
	UserType string `json:"user_type"` // guest | vip | svip | admin
	UserID   *uint  `json:"user_id,omitempty"`
	jwt.RegisteredClaims
}

// IssueUser signs a JWT for a user identified by their token record
func (s *JWTService) IssueUser(tokenID uint, userType string, userID *uint) (string, time.Time, error) {
	exp := time.Now().Add(time.Duration(s.expireHours) * time.Hour)
	claims := UserClaims{
		TokenID:  tokenID,
		UserType: userType,
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	return signed, exp, err
}

// VerifyUser parses and validates a user JWT
func (s *JWTService) VerifyUser(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		return s.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrTokenInvalid
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, ErrTokenInvalid
	}
	return claims, nil
}
