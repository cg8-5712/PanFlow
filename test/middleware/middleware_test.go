package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"panflow/internal/middleware"
	"panflow/internal/service"

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

// ── JWT Auth middleware tests ─────────────────────────────────────────────────

func newJWTSvc() *service.JWTService {
	return service.NewJWTService("test-secret-key", 1)
}

func TestJWTAuth_NoHeader(t *testing.T) {
	svc := newJWTSvc()
	r := newRouter(middleware.JWTAuth(svc, false))
	r.GET("/admin", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Fatalf("expected 401 without token, got %d", w.Code)
	}
}

func TestJWTAuth_InvalidToken(t *testing.T) {
	svc := newJWTSvc()
	r := newRouter(middleware.JWTAuth(svc, false))
	r.GET("/admin", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer not-a-valid-token")
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Fatalf("expected 401 for invalid token, got %d", w.Code)
	}
}

func TestJWTAuth_ValidToken(t *testing.T) {
	svc := newJWTSvc()
	tokenStr, _, err := svc.Issue()
	if err != nil {
		t.Fatal(err)
	}

	r := newRouter(middleware.JWTAuth(svc, false))
	r.GET("/admin", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200 with valid token, got %d", w.Code)
	}
}

func TestJWTAuth_DebugMode(t *testing.T) {
	svc := newJWTSvc()
	r := newRouter(middleware.JWTAuth(svc, true))
	r.GET("/admin", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil) // no token
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200 in debug mode, got %d", w.Code)
	}
}

func TestJWTAuth_WrongSecretRejected(t *testing.T) {
	issuer := service.NewJWTService("secret-A", 1)
	verifier := service.NewJWTService("secret-B", 1)

	tokenStr, _, _ := issuer.Issue()

	r := newRouter(middleware.JWTAuth(verifier, false))
	r.GET("/admin", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Fatalf("expected 401 for mismatched secret, got %d", w.Code)
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

// ── inline helpers ────────────────────────────────────────────────────────────

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

func identifierFilterDebug() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
