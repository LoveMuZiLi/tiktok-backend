package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"tiktok-app/backend/internal/config"
	"tiktok-app/backend/internal/model"
)

func Connect(cfg config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("connect mysql: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)

	if err := db.AutoMigrate(&model.Video{}); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	if err := Seed(db); err != nil {
		return nil, fmt.Errorf("seed: %w", err)
	}

	log.Println("mysql connected and schema ready")
	return db, nil
}
