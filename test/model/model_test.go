package model_test

import (
	"encoding/json"
	"testing"

	"panflow/internal/model"
)

// TestJSONSlice_Marshal 测试 JSONSlice 序列化
func TestJSONSlice_Marshal(t *testing.T) {
	slice := model.JSONSlice{"a", "b", "c"}
	data, err := json.Marshal(slice)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	expected := `["a","b","c"]`
	if string(data) != expected {
		t.Fatalf("expected %s, got %s", expected, string(data))
	}
}

// TestJSONSlice_Unmarshal 测试 JSONSlice 反序列化
func TestJSONSlice_Unmarshal(t *testing.T) {
	data := []byte(`["x","y","z"]`)
	var slice model.JSONSlice
	err := json.Unmarshal(data, &slice)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if len(slice) != 3 {
		t.Fatalf("expected 3 items, got %d", len(slice))
	}
	if slice[0] != "x" || slice[1] != "y" || slice[2] != "z" {
		t.Fatal("slice content mismatch")
	}
}

// TestJSONMap_Marshal 测试 JSONMap 序列化
func TestJSONMap_Marshal(t *testing.T) {
	m := model.JSONMap{
		"key1": "value1",
		"key2": 123,
		"key3": true,
	}
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if result["key1"] != "value1" {
		t.Fatal("key1 mismatch")
	}
}

// TestJSONMap_Unmarshal 测试 JSONMap 反序列化
func TestJSONMap_Unmarshal(t *testing.T) {
	data := []byte(`{"name":"test","age":25,"active":true}`)
	var m model.JSONMap
	err := json.Unmarshal(data, &m)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if m["name"] != "test" {
		t.Fatal("name mismatch")
	}
	if m["age"].(float64) != 25 {
		t.Fatal("age mismatch")
	}
	if m["active"] != true {
		t.Fatal("active mismatch")
	}
}

// TestJSONMap_NestedStructure 测试 JSONMap 嵌套结构
func TestJSONMap_NestedStructure(t *testing.T) {
	m := model.JSONMap{
		"cookie": "BDUSS=xxx",
		"nested": map[string]interface{}{
			"key": "value",
		},
	}

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var result model.JSONMap
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if result["cookie"] != "BDUSS=xxx" {
		t.Fatal("cookie mismatch")
	}
}

// TestRecord_Structure 测试 Record 模型
func TestRecord_Structure(t *testing.T) {
	userID := uint(1)
	record := &model.Record{
		IP:          "1.2.3.4",
		Fingerprint: "abc123",
		FsID:        100,
		URLs:        model.JSONSlice{"https://cdn1.example.com/file"},
		UA:          "Mozilla/5.0",
		TokenID:     1,
		AccountID:   1,
		UserID:      &userID,
	}

	if record.IP == "" {
		t.Fatal("IP should not be empty")
	}
	if len(record.URLs) == 0 {
		t.Fatal("URLs should not be empty")
	}
	if record.UserID == nil {
		t.Fatal("UserID should not be nil")
	}
}

// TestBlackList_Structure 测试 BlackList 模型
func TestBlackList_Structure(t *testing.T) {
	blacklist := &model.BlackList{
		Type:       "ip",
		Identifier: "1.2.3.4",
		Reason:     "恶意请求",
	}

	if blacklist.Type != "ip" {
		t.Fatal("type should be ip")
	}
	if blacklist.Identifier == "" {
		t.Fatal("identifier should not be empty")
	}
}

// TestBlackList_Types 测试黑名单类型
func TestBlackList_Types(t *testing.T) {
	types := []string{"ip", "fingerprint"}
	for _, typ := range types {
		bl := &model.BlackList{Type: typ}
		if bl.Type == "" {
			t.Fatalf("type %s should not be empty", typ)
		}
	}
}

// TestFileList_Structure 测试 FileList 模型
func TestFileList_Structure(t *testing.T) {
	file := &model.FileList{
		Surl:     "test_surl",
		Pwd:      "test_pwd",
		FsID:     "123456789",
		Size:     1024000,
		Filename: "test.mp4",
	}

	if file.FsID == "" {
		t.Fatal("FsID should not be empty")
	}
	if file.Size <= 0 {
		t.Fatal("Size should be positive")
	}
	if file.Filename == "" {
		t.Fatal("Filename should not be empty")
	}
}

// TestConfig_Structure 测试 Config 模型
func TestConfig_Structure(t *testing.T) {
	config := &model.Config{
		Key:         "guest_daily_limit",
		Value:       "5",
		Type:        "int",
		Description: "游客每日次数限制",
	}

	if config.Key == "" {
		t.Fatal("Key should not be empty")
	}
	if config.Value == "" {
		t.Fatal("Value should not be empty")
	}
}

// TestConfig_Types 测试配置类型
func TestConfig_Types(t *testing.T) {
	types := []string{"string", "int", "bool", "json"}
	for _, typ := range types {
		cfg := &model.Config{Type: typ}
		if cfg.Type == "" {
			t.Fatalf("type %s should not be empty", typ)
		}
	}
}
