package middleware

import (
	"net/http"
	"strings"

	"panflow/internal/config"
	"panflow/internal/handler"
	"panflow/internal/service"

	"github.com/gin-gonic/gin"
)

// Context keys for JWT claims injected by middleware
const (
	CtxTokenID  = "jwt_token_id"
	CtxUserType = "jwt_user_type"
	CtxUserID   = "jwt_user_id"
)

// JWTAuth validates the admin Bearer JWT in the Authorization header
func JWTAuth(jwtSvc *service.JWTService, debug bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if debug {
			c.Next()
			return
		}

		tokenStr := bearerToken(c)
		if tokenStr == "" {
			handler.Fail(c, http.StatusUnauthorized, 20001, "missing or invalid Authorization header")
			c.Abort()
			return
		}

		claims, err := jwtSvc.Verify(tokenStr)
		if err != nil {
			handler.Fail(c, http.StatusUnauthorized, 20001, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set("admin_role", claims.Role)
		c.Next()
	}
}

// UserJWTAuth validates the user Bearer JWT and injects token_id / user_type / user_id into context
func UserJWTAuth(jwtSvc *service.JWTService, debug bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if debug {
			// inject guest defaults so downstream handlers don't panic
			c.Set(CtxTokenID, uint(1))
			c.Set(CtxUserType, "guest")
			c.Next()
			return
		}

		tokenStr := bearerToken(c)
		if tokenStr == "" {
			handler.Fail(c, http.StatusUnauthorized, 20003, "login required")
			c.Abort()
			return
		}

		claims, err := jwtSvc.VerifyUser(tokenStr)
		if err != nil {
			handler.Fail(c, http.StatusUnauthorized, 20003, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set(CtxTokenID, claims.TokenID)
		c.Set(CtxUserType, claims.UserType)
		if claims.UserID != nil {
			c.Set(CtxUserID, *claims.UserID)
		}
		c.Next()
	}
}

// PassFilterUser validates the parse password from query or body
func PassFilterUser(cfg *config.PanflowConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.ParsePassword == "" {
			c.Next()
			return
		}

		password := c.Query("parse_password")
		if password == "" {
			var body struct {
				ParsePassword string `json:"parse_password" form:"parse_password"`
			}
			_ = c.ShouldBind(&body)
			password = body.ParsePassword
		}

		if password != cfg.ParsePassword {
			handler.Fail(c, http.StatusForbidden, 20002, "parse password error")
			c.Abort()
			return
		}

		c.Next()
	}
}

// bearerToken extracts the Bearer token from Authorization header
func bearerToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return ""
}
