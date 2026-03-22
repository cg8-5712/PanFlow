package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"panflow/internal/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// TestCors_AllHeaders 测试 CORS 所有头部
func TestCors_AllHeaders(t *testing.T) {
	r := gin.New()
	r.Use(middleware.Cors())
	r.GET("/test", func(c *gin.Context) {
		c.Status(200)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	headers := []string{
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Headers",
		"Access-Control-Max-Age",
	}

	for _, h := range headers {
		if w.Header().Get(h) == "" {
			t.Errorf("missing CORS header: %s", h)
		}
	}
}

// TestCors_OptionsRequest 测试 OPTIONS 预检请求
func TestCors_OptionsRequest(t *testing.T) {
	r := gin.New()
	r.Use(middleware.Cors())
	r.OPTIONS("/test", func(c *gin.Context) {})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != 204 {
		t.Fatalf("expected 204 for OPTIONS, got %d", w.Code)
	}
}

// TestCors_AllowOrigin 测试允许的源
func TestCors_AllowOrigin(t *testing.T) {
	r := gin.New()
	r.Use(middleware.Cors())
	r.GET("/test", func(c *gin.Context) {
		c.Status(200)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	r.ServeHTTP(w, req)

	origin := w.Header().Get("Access-Control-Allow-Origin")
	if origin != "*" {
		t.Fatalf("expected origin *, got %s", origin)
	}
}

// TestCors_AllowMethods 测试允许的方法
func TestCors_AllowMethods(t *testing.T) {
	r := gin.New()
	r.Use(middleware.Cors())
	r.GET("/test", func(c *gin.Context) {
		c.Status(200)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	methods := w.Header().Get("Access-Control-Allow-Methods")
	if methods == "" {
		t.Fatal("allow-methods should not be empty")
	}
}

// TestCors_AllowHeaders 测试允许的请求头
func TestCors_AllowHeaders(t *testing.T) {
	r := gin.New()
	r.Use(middleware.Cors())
	r.GET("/test", func(c *gin.Context) {
		c.Status(200)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	allowHeaders := w.Header().Get("Access-Control-Allow-Headers")
	if allowHeaders == "" {
		t.Fatal("allow-headers should not be empty")
	}
}

// TestCors_Credentials 测试凭据头部不干扰通配符源
func TestCors_Credentials(t *testing.T) {
	r := gin.New()
	r.Use(middleware.Cors())
	r.GET("/test", func(c *gin.Context) {
		c.Status(200)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	// 通配符 Origin 不设置 Allow-Credentials（浏览器安全限制）
	origin := w.Header().Get("Access-Control-Allow-Origin")
	if origin != "*" {
		t.Fatalf("expected wildcard origin, got %s", origin)
	}
}

// TestCors_MaxAge 测试最大缓存时间
func TestCors_MaxAge(t *testing.T) {
	r := gin.New()
	r.Use(middleware.Cors())
	r.GET("/test", func(c *gin.Context) {
		c.Status(200)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	maxAge := w.Header().Get("Access-Control-Max-Age")
	if maxAge == "" {
		t.Fatal("max-age should not be empty")
	}
}

// TestCors_MultipleRequests 测试多次请求
func TestCors_MultipleRequests(t *testing.T) {
	r := gin.New()
	r.Use(middleware.Cors())
	r.GET("/test", func(c *gin.Context) {
		c.Status(200)
	})

	for i := range 3 {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatalf("request %d: expected 200, got %d", i, w.Code)
		}
	}
}
