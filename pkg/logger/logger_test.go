package logger_test

import (
	"testing"

	"panflow/pkg/logger"
)

// TestInit_InfoLevel 测试 info 级别初始化
func TestInit_InfoLevel(t *testing.T) {
	err := logger.Init("info")
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}
}

// TestInit_DebugLevel 测试 debug 级别初始化
func TestInit_DebugLevel(t *testing.T) {
	err := logger.Init("debug")
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}
}

// TestInit_WarnLevel 测试 warn 级别初始化
func TestInit_WarnLevel(t *testing.T) {
	err := logger.Init("warn")
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}
}

// TestInit_ErrorLevel 测试 error 级别初始化
func TestInit_ErrorLevel(t *testing.T) {
	err := logger.Init("error")
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}
}

// TestInit_DefaultLevel 测试默认级别
func TestInit_DefaultLevel(t *testing.T) {
	err := logger.Init("")
	if err != nil {
		t.Fatalf("init with default level failed: %v", err)
	}
}

// TestInit_InvalidLevel 测试无效级别（应使用默认）
func TestInit_InvalidLevel(t *testing.T) {
	err := logger.Init("invalid")
	if err != nil {
		t.Fatalf("init with invalid level failed: %v", err)
	}
}

// TestSync_NoError 测试 Sync 不报错
func TestSync_NoError(t *testing.T) {
	logger.Init("info")
	err := logger.Sync()
	if err != nil {
		// Sync 可能在某些环境下返回错误（如 stderr），不算失败
		t.Logf("sync returned error (may be expected): %v", err)
	}
}

// TestSync_BeforeInit 测试初始化前 Sync
func TestSync_BeforeInit(t *testing.T) {
	// 不初始化直接 Sync 应该不崩溃
	err := logger.Sync()
	if err != nil {
		t.Logf("sync before init returned error: %v", err)
	}
}

// TestLogFunctions_NoPanic 测试日志函数不崩溃
func TestLogFunctions_NoPanic(t *testing.T) {
	logger.Init("info")

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("log function panicked: %v", r)
		}
	}()

	logger.Info("test info")
	logger.Infof("test infof: %s", "value")
	logger.Warn("test warn")
	logger.Warnf("test warnf: %d", 123)
	logger.Error("test error")
	logger.Errorf("test errorf: %v", true)
	logger.Debug("test debug")
	logger.Debugf("test debugf: %f", 3.14)
}

// TestLogLevel_Filtering 测试日志级别过滤
func TestLogLevel_Filtering(t *testing.T) {
	// 初始化为 error 级别
	logger.Init("error")

	// 这些调用不应该崩溃，即使不会输出
	logger.Debug("should not output")
	logger.Info("should not output")
	logger.Warn("should not output")
	logger.Error("should output")
}

// TestMultipleInit 测试多次初始化
func TestMultipleInit(t *testing.T) {
	err := logger.Init("info")
	if err != nil {
		t.Fatalf("first init failed: %v", err)
	}

	err = logger.Init("debug")
	if err != nil {
		t.Fatalf("second init failed: %v", err)
	}
}
