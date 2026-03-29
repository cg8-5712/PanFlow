package middleware

import (
	"context"
	"net/http"

	"panflow/internal/handler"
	"panflow/internal/repository"
	"panflow/pkg/cache"

	"github.com/gin-gonic/gin"
)

// IdentifierFilter blocks blacklisted IPs and browser fingerprints
func IdentifierFilter(blackListRepo *repository.BlackListRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		// Check IP blacklist (L2 cache first)
		ipKey := "blacklist:ip:" + ip
		if val, _ := cache.RedisGet(context.Background(), ipKey); val == "1" {
			handler.Fail(c, http.StatusForbidden, 20014, "your ip is blocked")
			c.Abort()
			return
		}

		blocked, err := blackListRepo.IsBlocked("ip", ip)
		if err == nil && blocked {
			_ = cache.RedisSet(context.Background(), ipKey, "1", 0)
			handler.Fail(c, http.StatusForbidden, 20014, "your ip is blocked")
			c.Abort()
			return
		}

		// Check fingerprint blacklist
		fingerprint := c.GetHeader("rand2")
		if fingerprint == "" {
			fingerprint = c.Query("rand2")
		}

		if fingerprint != "" {
			fpKey := "blacklist:fp:" + fingerprint
			if val, _ := cache.RedisGet(context.Background(), fpKey); val == "1" {
				handler.Fail(c, http.StatusForbidden, 20014, "your fingerprint is blocked")
				c.Abort()
				return
			}

			blocked, err := blackListRepo.IsBlocked("fingerprint", fingerprint)
			if err == nil && blocked {
				_ = cache.RedisSet(context.Background(), fpKey, "1", 0)
				handler.Fail(c, http.StatusForbidden, 20014, "your fingerprint is blocked")
				c.Abort()
				return
			}
		}

		c.Set("client_ip", ip)
		c.Set("fingerprint", fingerprint)
		c.Next()
	}
}
