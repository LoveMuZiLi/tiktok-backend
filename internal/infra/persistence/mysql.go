package persistence

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/LoveMuZiLi/tiktok-backend/internal/config"
	"github.com/LoveMuZiLi/tiktok-backend/internal/entity"
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

	if err := db.AutoMigrate(
		&entity.User{},
		&entity.Video{},
		&entity.Follow{},
		&entity.Conversation{},
		&entity.ConversationMember{},
		&entity.Message{},
		&entity.Notification{},
	); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	if err := migrateLegacyVideos(db); err != nil {
		return nil, fmt.Errorf("legacy migrate: %w", err)
	}

	if err := Seed(db); err != nil {
		return nil, fmt.Errorf("seed: %w", err)
	}

	log.Println("mysql connected and schema ready")
	return db, nil
}

// migrateLegacyVideos 将旧表（含 username/avatar 列）迁移到 user_id 关联
func migrateLegacyVideos(db *gorm.DB) error {
	if !db.Migrator().HasColumn(&entity.Video{}, "username") {
		return nil
	}
	var userCount int64
	if err := db.Model(&entity.User{}).Count(&userCount).Error; err != nil {
		return err
	}
	if userCount > 0 {
		return nil
	}

	type legacyVideo struct {
		ID          int64
		Username    string
		Avatar      string
		Description string
		Music       string
		Image       string
		Likes       int64
		Comments    int64
		Shares      int64
	}
	var rows []legacyVideo
	if err := db.Table("videos").Find(&rows).Error; err != nil {
		return err
	}
	for _, row := range rows {
		u := entity.User{
			Username: row.Username,
			Nickname: row.Username,
			Avatar:   row.Avatar,
			DouyinID: fmt.Sprintf("%d", 10000000+row.ID*1234567),
			Bio:      "",
		}
		if err := db.Where("username = ?", u.Username).FirstOrCreate(&u).Error; err != nil {
			return err
		}
		if err := db.Table("videos").Where("id = ?", row.ID).Update("user_id", u.ID).Error; err != nil {
			return err
		}
	}
	return nil
}
