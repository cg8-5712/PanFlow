package main

import (
	"fmt"
	"log"
	"os"

	"panflow/internal/config"
	"panflow/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. 初始化日志
	if err := logger.Init(cfg.Log.Level); err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("PanFlow starting...")
	logger.Infof("Server mode: %s", cfg.Server.Mode)

	// 3. 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 4. 创建 Gin 路由
	r := gin.Default()

	// 5. 健康检查路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 6. 启动服务器
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	logger.Infof("Server listening on %s", addr)

	if err := r.Run(addr); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
		os.Exit(1)
	}
}
