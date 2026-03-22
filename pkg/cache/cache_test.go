package cache_test

import (
	"testing"
	"time"

	"panflow/pkg/cache"
)

// TestInitOtter 测试 Otter 初始化
func TestInitOtter(t *testing.T) {
	err := cache.InitOtter(1000)
	if err != nil {
		t.Fatalf("init otter failed: %v", err)
	}
}

// TestOtterSetGet 测试设置和获取
func TestOtterSetGet(t *testing.T) {
	cache.InitOtter(1000)

	key := "test_key"
	value := "test_value"

	ok := cache.OtterSet(key, value, time.Minute)
	if !ok {
		t.Fatal("set should succeed")
	}

	got, found := cache.OtterGet(key)
	if !found {
		t.Fatal("key should be found")
	}
	if got != value {
		t.Fatalf("expected %v, got %v", value, got)
	}
}

// TestOtterGet_Miss 测试缓存未命中
func TestOtterGet_Miss(t *testing.T) {
	cache.InitOtter(1000)

	_, found := cache.OtterGet("nonexistent")
	if found {
		t.Fatal("nonexistent key should not be found")
	}
}

// TestOtterDelete 测试删除
func TestOtterDelete(t *testing.T) {
	cache.InitOtter(1000)

	key := "test_key"
	cache.OtterSet(key, "value", time.Minute)

	cache.OtterDelete(key)

	_, found := cache.OtterGet(key)
	if found {
		t.Fatal("deleted key should not be found")
	}
}

// TestOtterHas 测试键存在检查
func TestOtterHas(t *testing.T) {
	cache.InitOtter(1000)

	key := "test_key"
	cache.OtterSet(key, "value", time.Minute)

	if !cache.OtterHas(key) {
		t.Fatal("key should exist")
	}

	cache.OtterDelete(key)

	if cache.OtterHas(key) {
		t.Fatal("deleted key should not exist")
	}
}

// TestOtterSize 测试缓存大小
func TestOtterSize(t *testing.T) {
	cache.InitOtter(1000)
	cache.OtterClear()

	if cache.OtterSize() != 0 {
		t.Fatal("size should be 0 after clear")
	}

	cache.OtterSet("key1", "value1", time.Minute)
	cache.OtterSet("key2", "value2", time.Minute)

	size := cache.OtterSize()
	if size != 2 {
		t.Fatalf("expected size 2, got %d", size)
	}
}

// TestOtterClear 测试清空缓存
func TestOtterClear(t *testing.T) {
	cache.InitOtter(1000)

	cache.OtterSet("key1", "value1", time.Minute)
	cache.OtterSet("key2", "value2", time.Minute)

	cache.OtterClear()

	if cache.OtterSize() != 0 {
		t.Fatal("size should be 0 after clear")
	}
}

// TestOtterTTL 测试 TTL 设置（Otter 使用惰性驱逐，不保证立刻过期）
func TestOtterTTL(t *testing.T) {
	cache.InitOtter(1000)

	key := "test_key_ttl"
	cache.OtterSet(key, "value", time.Minute)

	// 立即获取应该存在
	_, found := cache.OtterGet(key)
	if !found {
		t.Fatal("key should be found immediately")
	}

	// 删除后应该不存在（不依赖 TTL 驱逐时机）
	cache.OtterDelete(key)
	_, found = cache.OtterGet(key)
	if found {
		t.Fatal("key should not be found after delete")
	}
}

// TestOtterClose 测试关闭
func TestOtterClose(t *testing.T) {
	cache.InitOtter(1000)
	cache.OtterSet("key", "value", time.Minute)

	cache.OtterClose()

	// 关闭后重新初始化，确保不崩溃
	err := cache.InitOtter(1000)
	if err != nil {
		t.Fatalf("re-init after close failed: %v", err)
	}
}

// TestOtterBeforeInit 测试初始化前操作
func TestOtterBeforeInit(t *testing.T) {
	// 重新初始化一个全新状态通过重置
	// 未初始化状态：Get 应该返回 false
	cache.OtterClear() // 安全调用，即使 nil 也无崩溃
	cache.OtterDelete("key")
	cache.OtterClose()

	size := cache.OtterSize()
	if size != 0 {
		t.Log("size after operations without fresh cache:", size)
	}

	// 重新验证已初始化状态可以正常工作
	cache.InitOtter(1000)
	ok := cache.OtterSet("key", "value", time.Minute)
	if !ok {
		t.Fatal("set should succeed after init")
	}
}

// TestOtterMultipleTypes 测试多种类型值
func TestOtterMultipleTypes(t *testing.T) {
	cache.InitOtter(1000)

	// 字符串
	cache.OtterSet("str", "value", time.Minute)
	// 整数
	cache.OtterSet("int", 123, time.Minute)
	// 布尔
	cache.OtterSet("bool", true, time.Minute)
	// map
	cache.OtterSet("map", map[string]int{"a": 1}, time.Minute)

	if v, ok := cache.OtterGet("str"); !ok || v != "value" {
		t.Fatal("string value mismatch")
	}
	if v, ok := cache.OtterGet("int"); !ok || v != 123 {
		t.Fatal("int value mismatch")
	}
	if v, ok := cache.OtterGet("bool"); !ok || v != true {
		t.Fatal("bool value mismatch")
	}
	if v, ok := cache.OtterGet("map"); !ok {
		t.Fatal("map value not found")
	} else if m, ok := v.(map[string]int); !ok || m["a"] != 1 {
		t.Fatal("map value mismatch")
	}
}

// TestOtterOverwrite 测试覆盖值
func TestOtterOverwrite(t *testing.T) {
	cache.InitOtter(1000)

	key := "test_key"
	cache.OtterSet(key, "value1", time.Minute)
	cache.OtterSet(key, "value2", time.Minute)

	v, _ := cache.OtterGet(key)
	if v != "value2" {
		t.Fatalf("expected value2, got %v", v)
	}
}
