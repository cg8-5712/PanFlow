package config_test

import (
	"testing"

	"panflow/internal/config"
)

// TestServerConfig_Structure 测试服务器配置结构
func TestServerConfig_Structure(t *testing.T) {
	cfg := config.ServerConfig{
		Host: "0.0.0.0",
		Port: "8080",
		Mode: "release",
	}

	if cfg.Host == "" {
		t.Fatal("host should not be empty")
	}
	if cfg.Port == "" {
		t.Fatal("port should not be empty")
	}
	if cfg.Mode == "" {
		t.Fatal("mode should not be empty")
	}
}

// TestDatabaseConfig_Structure 测试数据库配置结构
func TestDatabaseConfig_Structure(t *testing.T) {
	cfg := config.DatabaseConfig{
		Driver:   "mysql",
		Host:     "127.0.0.1",
		Port:     "3306",
		User:     "root",
		Password: "password",
		Name:     "panflow",
	}

	if cfg.Driver == "" {
		t.Fatal("driver should not be empty")
	}
	if cfg.Name == "" {
		t.Fatal("database name should not be empty")
	}
}

// TestRedisConfig_Structure 测试 Redis 配置结构
func TestRedisConfig_Structure(t *testing.T) {
	cfg := config.RedisConfig{
		Host:     "127.0.0.1",
		Port:     "6379",
		Password: "",
		DB:       0,
	}

	if cfg.Host == "" {
		t.Fatal("host should not be empty")
	}
	if cfg.Port == "" {
		t.Fatal("port should not be empty")
	}
}

// TestLogConfig_Structure 测试日志配置结构
func TestLogConfig_Structure(t *testing.T) {
	cfg := config.LogConfig{
		Level: "info",
	}

	if cfg.Level == "" {
		t.Fatal("log level should not be empty")
	}
}

// TestPanflowConfig_Structure 测试 Panflow 配置结构
func TestPanflowConfig_Structure(t *testing.T) {
	cfg := config.PanflowConfig{
		AdminPassword:  "admin123",
		JWTSecret:      "secret",
		JWTExpireHours: 2,
		JWTRefreshDays: 7,
		GuestUserAgent: "netdisk;P2SP;3.0.20.138",
		ProxyHTTP:      "",
	}

	if cfg.JWTSecret == "" {
		t.Fatal("jwt_secret should not be empty")
	}
	if cfg.JWTExpireHours <= 0 {
		t.Fatal("jwt_expire_hours should be positive")
	}
}

// TestConfig_Complete 测试完整配置
func TestConfig_Complete(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "0.0.0.0",
			Port: "8080",
			Mode: "release",
		},
		Database: config.DatabaseConfig{
			Driver: "mysql",
			Host:   "127.0.0.1",
			Port:   "3306",
			Name:   "panflow",
		},
		Redis: config.RedisConfig{
			Host: "127.0.0.1",
			Port: "6379",
			DB:   0,
		},
		Log: config.LogConfig{
			Level: "info",
		},
		Panflow: config.PanflowConfig{
			JWTSecret:      "secret",
			JWTExpireHours: 2,
		},
	}

	if cfg.Server.Port == "" {
		t.Fatal("server port should not be empty")
	}
	if cfg.Database.Name == "" {
		t.Fatal("database name should not be empty")
	}
	if cfg.Panflow.JWTSecret == "" {
		t.Fatal("jwt secret should not be empty")
	}
}

// TestDatabaseDriver_Validation 测试数据库驱动验证
func TestDatabaseDriver_Validation(t *testing.T) {
	validDrivers := []string{"mysql", "postgres", "sqlite"}
	for _, driver := range validDrivers {
		cfg := config.DatabaseConfig{
			Driver: driver,
		}
		if cfg.Driver == "" {
			t.Fatalf("driver %s should not be empty", driver)
		}
	}
}

// TestServerMode_Validation 测试服务器模式验证
func TestServerMode_Validation(t *testing.T) {
	validModes := []string{"debug", "release", "test"}
	for _, mode := range validModes {
		cfg := config.ServerConfig{
			Mode: mode,
		}
		if cfg.Mode == "" {
			t.Fatalf("mode %s should not be empty", mode)
		}
	}
}

// TestLogLevel_Validation 测试日志级别验证
func TestLogLevel_Validation(t *testing.T) {
	validLevels := []string{"debug", "info", "warn", "error"}
	for _, level := range validLevels {
		cfg := config.LogConfig{
			Level: level,
		}
		if cfg.Level == "" {
			t.Fatalf("level %s should not be empty", level)
		}
	}
}

// TestPanflowConfig_Defaults 测试默认值
func TestPanflowConfig_Defaults(t *testing.T) {
	cfg := config.PanflowConfig{
		JWTExpireHours: 2,
		JWTRefreshDays: 7,
	}

	if cfg.JWTExpireHours != 2 {
		t.Fatal("default jwt_expire_hours should be 2")
	}
	if cfg.JWTRefreshDays != 7 {
		t.Fatal("default jwt_refresh_days should be 7")
	}
}
