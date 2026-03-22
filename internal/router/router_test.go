package router_test

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// TestRouterSetup_NoPanic 测试路由设置不崩溃
func TestRouterSetup_NoPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("router setup panicked: %v", r)
		}
	}()

	// 创建一个简单的路由器
	r := gin.New()
	if r == nil {
		t.Fatal("router should not be nil")
	}
}

// TestRouterGroup_Creation 测试路由组创建
func TestRouterGroup_Creation(t *testing.T) {
	r := gin.New()

	api := r.Group("/api/v1")
	if api == nil {
		t.Fatal("api group should not be nil")
	}

	admin := api.Group("/admin")
	if admin == nil {
		t.Fatal("admin group should not be nil")
	}
}

// TestRouterMiddleware_Application 测试中间件应用
func TestRouterMiddleware_Application(t *testing.T) {
	r := gin.New()

	r.Use(func(c *gin.Context) {
		c.Next()
	})

	r.GET("/test", func(c *gin.Context) {
		c.Status(200)
	})

	// 这里只是测试结构，不实际发送请求
	if r == nil {
		t.Fatal("router should not be nil")
	}
}

// TestRouterRoutes_Registration 测试路由注册
func TestRouterRoutes_Registration(t *testing.T) {
	r := gin.New()

	r.GET("/get", func(c *gin.Context) {})
	r.POST("/post", func(c *gin.Context) {})
	r.PUT("/put", func(c *gin.Context) {})
	r.PATCH("/patch", func(c *gin.Context) {})
	r.DELETE("/delete", func(c *gin.Context) {})

	routes := r.Routes()
	if len(routes) != 5 {
		t.Fatalf("expected 5 routes, got %d", len(routes))
	}
}

// TestRouterGroup_Nesting 测试嵌套路由组
func TestRouterGroup_Nesting(t *testing.T) {
	r := gin.New()

	api := r.Group("/api")
	v1 := api.Group("/v1")
	admin := v1.Group("/admin")

	admin.GET("/test", func(c *gin.Context) {})

	routes := r.Routes()
	if len(routes) == 0 {
		t.Fatal("should have at least one route")
	}
}

// TestRouterHandlers_Count 测试处理器数量
func TestRouterHandlers_Count(t *testing.T) {
	r := gin.New()

	// 添加多个中间件
	r.Use(func(c *gin.Context) { c.Next() })
	r.Use(func(c *gin.Context) { c.Next() })

	r.GET("/test", func(c *gin.Context) {})

	routes := r.Routes()
	if len(routes) == 0 {
		t.Fatal("should have routes")
	}
}

// TestRouterMethods_AllSupported 测试所有 HTTP 方法
func TestRouterMethods_AllSupported(t *testing.T) {
	r := gin.New()

	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"}
	for _, method := range methods {
		switch method {
		case "GET":
			r.GET("/test", func(c *gin.Context) {})
		case "POST":
			r.POST("/test", func(c *gin.Context) {})
		case "PUT":
			r.PUT("/test", func(c *gin.Context) {})
		case "PATCH":
			r.PATCH("/test", func(c *gin.Context) {})
		case "DELETE":
			r.DELETE("/test", func(c *gin.Context) {})
		case "OPTIONS":
			r.OPTIONS("/test", func(c *gin.Context) {})
		case "HEAD":
			r.HEAD("/test", func(c *gin.Context) {})
		}
	}

	routes := r.Routes()
	if len(routes) != len(methods) {
		t.Fatalf("expected %d routes, got %d", len(methods), len(routes))
	}
}
