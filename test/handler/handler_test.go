package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// responseBody is a helper to decode the unified response
type responseBody struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func newTestRouter() *gin.Engine {
	r := gin.New()
	return r
}

func doRequest(r *gin.Engine, method, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, nil)
	r.ServeHTTP(w, req)
	return w
}

func TestPingRoute(t *testing.T) {
	r := newTestRouter()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	w := doRequest(r, "GET", "/ping")
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse body: %v", err)
	}
	if body["message"] != "pong" {
		t.Fatalf("expected pong, got %s", body["message"])
	}
}

func TestNotFound(t *testing.T) {
	r := newTestRouter()
	w := doRequest(r, "GET", "/nonexistent")
	if w.Code != 404 {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestResponseFormat(t *testing.T) {
	r := newTestRouter()
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, responseBody{
			Code:    0,
			Message: "success",
			Data:    json.RawMessage(`{"key":"value"}`),
		})
	})

	w := doRequest(r, "GET", "/test")
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body responseBody
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse body: %v", err)
	}
	if body.Code != 0 {
		t.Fatalf("expected code 0, got %d", body.Code)
	}
	if body.Message != "success" {
		t.Fatalf("expected success, got %s", body.Message)
	}
}

func TestErrorResponseFormat(t *testing.T) {
	r := newTestRouter()
	r.GET("/error", func(c *gin.Context) {
		c.JSON(400, responseBody{
			Code:    40000,
			Message: "bad request",
		})
	})

	w := doRequest(r, "GET", "/error")
	if w.Code != 400 {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	var body responseBody
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse body: %v", err)
	}
	if body.Code != 40000 {
		t.Fatalf("expected code 40000, got %d", body.Code)
	}
}

func TestCORSHeaders(t *testing.T) {
	r := newTestRouter()
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Next()
	})
	r.GET("/cors", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := doRequest(r, "GET", "/cors")
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("expected CORS header")
	}
}

func TestCORSPreflight(t *testing.T) {
	r := newTestRouter()
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	r.OPTIONS("/preflight", func(c *gin.Context) {})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/preflight", nil)
	r.ServeHTTP(w, req)

	if w.Code != 204 {
		t.Fatalf("expected 204 for preflight, got %d", w.Code)
	}
}

func TestAdminPasswordMissing(t *testing.T) {
	r := newTestRouter()
	r.GET("/admin/test", func(c *gin.Context) {
		password := c.GetHeader("admin_password")
		if password == "" {
			c.JSON(403, responseBody{Code: 20001, Message: "admin password error"})
			return
		}
		c.JSON(200, responseBody{Code: 0, Message: "ok"})
	})

	w := doRequest(r, "GET", "/admin/test")
	if w.Code != 403 {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestAdminPasswordCorrect(t *testing.T) {
	r := newTestRouter()
	r.GET("/admin/test", func(c *gin.Context) {
		password := c.GetHeader("admin_password")
		if password != "secret" {
			c.JSON(403, responseBody{Code: 20001, Message: "admin password error"})
			return
		}
		c.JSON(200, responseBody{Code: 0, Message: "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin/test", nil)
	req.Header.Set("admin_password", "secret")
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestPaginationDefaults(t *testing.T) {
	type pageQuery struct {
		Page  int `form:"page"`
		Limit int `form:"limit"`
	}

	r := newTestRouter()
	r.GET("/list", func(c *gin.Context) {
		var q pageQuery
		_ = c.ShouldBindQuery(&q)
		if q.Page < 1 {
			q.Page = 1
		}
		if q.Limit < 1 || q.Limit > 100 {
			q.Limit = 20
		}
		c.JSON(200, gin.H{"page": q.Page, "limit": q.Limit})
	})

	w := doRequest(r, "GET", "/list")
	var body map[string]int
	_ = json.Unmarshal(w.Body.Bytes(), &body)

	if body["page"] != 1 {
		t.Fatalf("expected page=1, got %d", body["page"])
	}
	if body["limit"] != 20 {
		t.Fatalf("expected limit=20, got %d", body["limit"])
	}
}
