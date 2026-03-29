package repository

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"panflow/internal/config"
	"panflow/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// NewDB opens a database connection and runs AutoMigrate.
func NewDB(cfg config.DatabaseConfig, dev bool) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch cfg.Driver {
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
		dialector = postgres.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(cfg.Name)
	default: // mysql
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
		dialector = mysql.Open(dsn)
	}

	gormCfg := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	}

	db, err := gorm.Open(dialector, gormCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	if err := seedDefaultConfigs(db); err != nil {
		return nil, fmt.Errorf("seed default configs failed: %w", err)
	}

	if dev {
		if err := seedDevData(db); err != nil {
			return nil, fmt.Errorf("seed dev data failed: %w", err)
		}
	} else {
		if err := seedAdminUser(db, ""); err != nil {
			return nil, fmt.Errorf("seed admin user failed: %w", err)
		}
	}

	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Account{},
		&model.User{},
		&model.Config{},
		&model.FileList{},
		&model.Record{},
		&model.BlackList{},
	)
}

// seedAdminUser 确保 admin 用户存在，password 为空时跳过创建（生产环境由运维手动创建）
func seedAdminUser(db *gorm.DB, password string) error {
	var u model.User
	if err := db.Where("username = ?", "admin").First(&u).Error; err == nil {
		return nil // already exists
	}
	if password == "" {
		return nil // 生产环境不自动创建
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return db.Create(&model.User{
		Username: "admin",
		Password: string(hash),
		UserType: "admin",
	}).Error
}

// seedDevData 开发模式：随机 admin 密码 + 测试用户数据
func seedDevData(db *gorm.DB) error {
	// 随机生成 admin 密码
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	adminPwd := hex.EncodeToString(b)

	hash, err := bcrypt.GenerateFromPassword([]byte(adminPwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var adminUser model.User
	if db.Where("username = ?", "admin").First(&adminUser).Error != nil {
		if err := db.Create(&model.User{
			Username: "admin",
			Password: string(hash),
			UserType: "admin",
		}).Error; err != nil {
			return err
		}
		fmt.Printf("\n====================================\n")
		fmt.Printf(" DEV MODE — Admin credentials:\n")
		fmt.Printf(" username: admin\n")
		fmt.Printf(" password: %s\n", adminPwd)
		fmt.Printf("====================================\n\n")
	}

	// Seed 测试用户
	testUsers := []struct {
		username string
		password string
		userType string
	}{
		{"test_guest", "guest123", "guest"},
		{"test_vip", "vip123", "vip"},
	}
	for _, u := range testUsers {
		var existing model.User
		if db.Where("username = ?", u.username).First(&existing).Error == nil {
			continue
		}
		h, _ := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
		db.Create(&model.User{
			Username: u.username,
			Password: string(h),
			UserType: u.userType,
		})
	}

	return seedDefaultConfigs(db)
}

// seedDefaultConfigs ensures default configuration entries exist
func seedDefaultConfigs(db *gorm.DB) error {
	defaults := []model.Config{
		{Key: "guest_daily_limit", Value: "5", Type: "int", Description: "普通用户每日次数"},
		{Key: "vip_count_based", Value: "true", Type: "bool", Description: "VIP按次数计费"},
		{Key: "svip_daily_limit", Value: "100", Type: "int", Description: "SVIP用户每日次数"},
		{Key: "admin_unlimited", Value: "true", Type: "bool", Description: "Admin无限制"},
	}

	for _, cfg := range defaults {
		var existing model.Config
		if db.Where("`key` = ?", cfg.Key).First(&existing).Error == nil {
			continue
		}
		if err := db.Create(&cfg).Error; err != nil {
			return err
		}
	}
	return nil
}
