package main

import (
	"fmt"
	"log"
	"os"

	"panflow/internal/config"
	"panflow/internal/repository"
	"panflow/internal/router"
	"panflow/pkg/cache"
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

	// 3. 初始化 L1 缓存（Otter）
	if err := cache.InitOtter(10000); err != nil {
		logger.Fatalf("Failed to init otter cache: %v", err)
	}
	defer cache.OtterClose()
	logger.Info("L1 cache (Otter) initialized")

	// 4. 初始化 L2 缓存（Redis）
	if err := cache.InitRedis(cfg.Redis); err != nil {
		logger.Warnf("Redis unavailable, L2 cache disabled: %v", err)
	} else {
		defer cache.RedisClose()
		logger.Info("L2 cache (Redis) initialized")
	}

	// 5. 初始化数据库（AutoMigrate + seed）
	db, err := repository.NewDB(cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to init database: %v", err)
	}
	logger.Info("Database initialized")

	// 6. 初始化各 Repository
	accountRepo := repository.NewAccountRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	userRepo := repository.NewUserRepository(db)
	configRepo := repository.NewConfigRepository(db)
	recordRepo := repository.NewRecordRepository(db)
	fileListRepo := repository.NewFileListRepository(db)
	blackListRepo := repository.NewBlackListRepository(db)

	logger.Info("Repositories initialized")

	// 7. 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 8. 创建 Gin 路由
	r := gin.Default()

	// 9. 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// 10. 注册路由
	router.Setup(r, cfg, accountRepo, tokenRepo, userRepo, configRepo, recordRepo, fileListRepo, blackListRepo)

	// 11. 启动服务器
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	logger.Infof("Server listening on %s", addr)

	if err := r.Run(addr); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
		os.Exit(1)
	}
}
