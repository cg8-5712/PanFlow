package middleware

import (
	"net/http"

	"panflow/internal/config"
	"panflow/internal/handler"

	"github.com/gin-gonic/gin"
)

// PassFilterAdmin validates the admin password from header, query, or body
func PassFilterAdmin(cfg *config.HklistConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.Debug {
			c.Next()
			return
		}

		password := c.GetHeader("admin_password")
		if password == "" {
			password = c.Query("admin_password")
		}
		if password == "" {
			var body struct {
				AdminPassword string `json:"admin_password" form:"admin_password"`
			}
			_ = c.ShouldBind(&body)
			password = body.AdminPassword
		}

		if password != cfg.AdminPassword {
			handler.FailForbidden(c, 20001, "admin password error")
			c.Abort()
			return
		}

		c.Next()
	}
}

// PassFilterUser validates the parse password from query or body
func PassFilterUser(cfg *config.HklistConfig) gin.HandlerFunc {
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
