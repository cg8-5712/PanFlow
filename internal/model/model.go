package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// LocalTime wraps time.Time to accept both RFC3339 and datetime-local (no timezone) from JSON
type LocalTime struct{ time.Time }

func (t *LocalTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" || s == "\"\"" {
		return nil
	}
	// strip quotes
	if len(s) >= 2 && s[0] == '"' {
		s = s[1 : len(s)-1]
	}
	for _, layout := range []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
	} {
		if parsed, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			t.Time = parsed
			return nil
		}
	}
	return fmt.Errorf("cannot parse time: %q", s)
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + t.Format("2006-01-02T15:04:05") + `"`), nil
}

func (t LocalTime) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *LocalTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		t.Time = v
	}
	return nil
}

// JSONSlice is a []string that serialises to/from JSON in MySQL
type JSONSlice []string

func (j JSONSlice) Value() (driver.Value, error) {
	b, _ := json.Marshal(j)
	return string(b), nil
}
func (j *JSONSlice) Scan(value interface{}) error {
	var b []byte
	switch v := value.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("unsupported type for JSONSlice: %T", value)
	}
	return json.Unmarshal(b, j)
}

// JSONMap is a map[string]interface{} for account_data
type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	b, _ := json.Marshal(j)
	return string(b), nil
}
func (j *JSONMap) Scan(value interface{}) error {
	var b []byte
	switch v := value.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("unsupported type for JSONMap: %T", value)
	}
	return json.Unmarshal(b, j)
}

// Account represents a Baidu SVIP account
type Account struct {
	ID                 uint           `gorm:"primaryKey"                      json:"id"`
	BaiduName          string         `gorm:"column:baidu_name"               json:"baidu_name"`
	UK                 string         `gorm:"column:uk"                       json:"uk"`
	AccountType        string         `gorm:"column:account_type"             json:"account_type"`
	AccountData        JSONMap        `gorm:"column:account_data;type:text"   json:"account_data"`
	Switch             bool           `gorm:"column:switch;default:true"      json:"switch"`
	Reason             string         `gorm:"column:reason"                   json:"reason"`
	Prov               *string        `gorm:"column:prov"                     json:"prov"`
	ProviderUserID     *uint          `gorm:"column:provider_user_id;index:idx_accounts_provider_user" json:"provider_user_id"`
	UsedCount          int64          `gorm:"column:used_count;default:0"     json:"used_count"`
	UsedSize           int64          `gorm:"column:used_size;default:0"      json:"used_size"`
	TotalSize          int64          `gorm:"column:total_size;default:0"     json:"total_size"`
	TotalSizeUpdatedAt *time.Time     `gorm:"column:total_size_updated_at"    json:"total_size_updated_at"`
	LastUseAt          *time.Time     `gorm:"column:last_use_at"              json:"last_use_at"`
	CreatedAt          time.Time      `                                        json:"created_at"`
	UpdatedAt          time.Time      `                                        json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index"                           json:"deleted_at"`
}

// FileList caches Baidu Pan file metadata from share links
type FileList struct {
	ID        uint      `gorm:"primaryKey"  json:"id"`
	Surl      string    `gorm:"column:surl" json:"surl"`
	Pwd       string    `gorm:"column:pwd"  json:"pwd"`
	FsID      string    `gorm:"column:fs_id;uniqueIndex" json:"fs_id"`
	Size      int64     `gorm:"column:size" json:"size"`
	Filename  string    `gorm:"column:filename" json:"filename"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Record stores one download history entry
type Record struct {
	ID          uint      `gorm:"primaryKey"           json:"id"`
	IP          string    `gorm:"column:ip"            json:"ip"`
	Fingerprint string    `gorm:"column:fingerprint"   json:"fingerprint"`
	FsID        uint      `gorm:"column:fs_id"         json:"fs_id"`
	URLs        JSONSlice `gorm:"column:urls;type:text" json:"urls"`
	UA          string    `gorm:"column:ua"            json:"ua"`
	AccountID   uint      `gorm:"column:account_id;index:idx_account" json:"account_id"`
	UserID      *uint     `gorm:"column:user_id;index:idx_user" json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	// Preload associations
	Account *Account  `gorm:"foreignKey:AccountID"                 json:"account,omitempty"`
	User    *User     `gorm:"foreignKey:UserID"                    json:"user,omitempty"`
	File    *FileList `gorm:"foreignKey:FsID;references:ID"        json:"file,omitempty"`
}

// BlackList stores blocked IPs and browser fingerprints
type BlackList struct {
	ID         uint       `gorm:"primaryKey"        json:"id"`
	Type       string     `gorm:"column:type"       json:"type"` // ip | fingerprint
	Identifier string     `gorm:"column:identifier" json:"identifier"`
	Reason     string     `gorm:"column:reason"     json:"reason"`
	ExpiresAt  *LocalTime `gorm:"column:expires_at" json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// User represents a user account with different privilege levels
type User struct {
	ID             uint           `gorm:"primaryKey"                                json:"id"`
	Username       string         `gorm:"column:username;uniqueIndex;not null"      json:"username"`
	Password       string         `gorm:"column:password"                           json:"-"`
	Email          string         `gorm:"column:email"                              json:"email"`
	UserType       string         `gorm:"column:user_type;default:guest;index:idx_users_user_type" json:"user_type"` // guest | vip | svip | admin
	VipBalance     int64          `gorm:"column:vip_balance;default:0"              json:"vip_balance"`
	DailyUsedCount int64          `gorm:"column:daily_used_count;default:0"         json:"daily_used_count"`
	DailyLimit     int            `gorm:"column:daily_limit;default:5"              json:"daily_limit"`
	BaiduAccountID *uint          `gorm:"column:baidu_account_id"                   json:"baidu_account_id"`
	CreatedAt      time.Time      `                                                  json:"created_at"`
	UpdatedAt      time.Time      `                                                  json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index"                                     json:"deleted_at"`
	// Preload associations
	BaiduAccount *Account `gorm:"foreignKey:BaiduAccountID" json:"baidu_account,omitempty"`
}

// Config represents a database-stored configuration entry
type Config struct {
	ID          uint      `gorm:"primaryKey"                   json:"id"`
	Key         string    `gorm:"column:key;uniqueIndex;not null" json:"key"`
	Value       string    `gorm:"column:value;type:text"       json:"value"`
	Type        string    `gorm:"column:type;default:string"   json:"type"` // string | int | bool | json
	Description string    `gorm:"column:description"           json:"description"`
	CreatedAt   time.Time `                                     json:"created_at"`
	UpdatedAt   time.Time `                                     json:"updated_at"`
}
