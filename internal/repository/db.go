package repository

import (
	"fmt"

	"panflow/internal/config"
	"panflow/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// NewDB opens a MySQL connection and runs AutoMigrate
func NewDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	gormCfg := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	}

	db, err := gorm.Open(mysql.Open(dsn), gormCfg)
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

	if err := seedGuestToken(db); err != nil {
		return nil, fmt.Errorf("seed guest token failed: %w", err)
	}

	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Account{},
		&model.Token{},
		&model.FileList{},
		&model.Record{},
		&model.BlackList{},
		&model.Proxy{},
	)
}

// seedGuestToken ensures a guest token (id=1) always exists
func seedGuestToken(db *gorm.DB) error {
	var t model.Token
	err := db.First(&t, 1).Error
	if err == nil {
		return nil // already exists
	}
	guest := model.Token{
		Token:         "guest",
		TokenType:     "daily",
		Count:         10,
		Size:          10 * 1024 * 1024 * 1024, // 10 GB
		Day:           1,
		CanUseIPCount: 99999,
		IP:            model.JSONSlice{},
		Switch:        true,
		Reason:        "",
	}
	return db.Create(&guest).Error
}
