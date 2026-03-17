package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.Use(middlewares...)
	return r
}

func TestCorsMiddleware_AllowsGET(t *testing.T) {
	r := newRouter(corsMiddleware())
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("missing CORS header")
	}
}

func TestCorsMiddleware_Preflight(t *testing.T) {
	r := newRouter(corsMiddleware())
	r.OPTIONS("/test", func(c *gin.Context) {})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != 204 {
		t.Fatalf("expected 204 for OPTIONS, got %d", w.Code)
	}
}

func TestPassFilterAdmin_NoPassword(t *testing.T) {
	r := newRouter(passFilterAdmin("secret"))
	r.GET("/admin", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	r.ServeHTTP(w, req)

	if w.Code != 403 {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestPassFilterAdmin_WrongPassword(t *testing.T) {
	r := newRouter(passFilterAdmin("secret"))
	r.GET("/admin", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	req.Header.Set("admin_password", "wrong")
	r.ServeHTTP(w, req)

	if w.Code != 403 {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestPassFilterAdmin_CorrectHeader(t *testing.T) {
	r := newRouter(passFilterAdmin("secret"))
	r.GET("/admin", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	req.Header.Set("admin_password", "secret")
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestPassFilterAdmin_CorrectQuery(t *testing.T) {
	r := newRouter(passFilterAdmin("secret"))
	r.GET("/admin", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin?admin_password=secret", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestPassFilterAdmin_EmptyPassword_AlwaysPass(t *testing.T) {
	// When configured password is empty, admin filter still checks
	r := newRouter(passFilterAdmin(""))
	r.GET("/admin", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	r.ServeHTTP(w, req)

	// empty password matches empty header → should pass
	if w.Code != 200 {
		t.Fatalf("expected 200 when password is empty, got %d", w.Code)
	}
}

func TestIdentifierFilter_DebugMode(t *testing.T) {
	r := newRouter(identifierFilterDebug())
	r.GET("/user", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200 in debug mode, got %d", w.Code)
	}
}

// ── inline middleware implementations for testing ─────────────────────────────

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, admin_password, parse_password")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func passFilterAdmin(configuredPassword string) gin.HandlerFunc {
	return func(c *gin.Context) {
		password := c.GetHeader("admin_password")
		if password == "" {
			password = c.Query("admin_password")
		}
		if password != configuredPassword {
			c.JSON(403, gin.H{"code": 20001, "message": "admin password error"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func identifierFilterDebug() gin.HandlerFunc {
	return func(c *gin.Context) {
		// debug=true → always pass
		c.Next()
	}
}
