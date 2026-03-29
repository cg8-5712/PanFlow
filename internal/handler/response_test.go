package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"panflow/internal/handler"
	"panflow/pkg/logger"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
	_ = logger.Init("debug")
}

// TestResponse_Structure 测试统一响应结构
func TestResponse_Structure(t *testing.T) {
	resp := handler.Response{
		Code:    0,
		Message: "success",
		Data:    map[string]string{"key": "value"},
	}

	if resp.Code != 0 {
		t.Fatal("code should be 0")
	}
	if resp.Message != "success" {
		t.Fatal("message should be success")
	}
}

// TestSuccess 测试成功响应
func TestSuccess(t *testing.T) {
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		handler.Success(c, gin.H{"result": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp handler.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 0 {
		t.Fatalf("expected code 0, got %d", resp.Code)
	}
}

// TestSuccessMsg 测试消息响应
func TestSuccessMsg(t *testing.T) {
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		handler.SuccessMsg(c, "operation completed")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	var resp handler.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Message != "operation completed" {
		t.Fatalf("expected custom message, got %s", resp.Message)
	}
}

// TestFailBadRequest 测试 400 错误
func TestFailBadRequest(t *testing.T) {
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		handler.FailBadRequest(c, 40000, "invalid input")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	var resp handler.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 40000 {
		t.Fatalf("expected code 40000, got %d", resp.Code)
	}
}

// TestFailUnauthorized 测试 401 错误
func TestFailUnauthorized(t *testing.T) {
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		handler.FailUnauthorized(c, 20001, "unauthorized")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

// TestFailForbidden 测试 403 错误
func TestFailForbidden(t *testing.T) {
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		handler.FailForbidden(c, 20014, "forbidden")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != 403 {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

// TestFailNotFound 测试 404 错误
func TestFailNotFound(t *testing.T) {
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		handler.FailNotFound(c, 40400, "not found")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// TestFailInternal 测试 500 错误
func TestFailInternal(t *testing.T) {
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		handler.FailInternal(c, "internal error")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Fatalf("expected 500, got %d", w.Code)
	}

	var resp handler.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 50000 {
		t.Fatalf("expected code 50000, got %d", resp.Code)
	}
}

// TestFail_CustomStatus 测试自定义状态码
func TestFail_CustomStatus(t *testing.T) {
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		handler.Fail(c, 418, 41800, "I'm a teapot")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != 418 {
		t.Fatalf("expected 418, got %d", w.Code)
	}
}

// TestResponse_JSONFormat 测试 JSON 格式
func TestResponse_JSONFormat(t *testing.T) {
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		handler.Success(c, gin.H{
			"items": []string{"a", "b", "c"},
			"count": 3,
		})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	var resp handler.Response
	json.Unmarshal(w.Body.Bytes(), &resp)

	data, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatal("data should be a map")
	}
	if data["count"].(float64) != 3 {
		t.Fatal("count should be 3")
	}
}

// TestResponse_NilData 测试空数据
func TestResponse_NilData(t *testing.T) {
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		handler.SuccessMsg(c, "no data")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	var resp handler.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Data != nil {
		t.Fatal("data should be nil")
	}
}

// TestResponse_WithPOST 测试 POST 请求响应
func TestResponse_WithPOST(t *testing.T) {
	r := gin.New()
	r.POST("/test", func(c *gin.Context) {
		var body map[string]string
		c.ShouldBindJSON(&body)
		handler.Success(c, body)
	})

	reqBody := bytes.NewBufferString(`{"name":"test"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", reqBody)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
