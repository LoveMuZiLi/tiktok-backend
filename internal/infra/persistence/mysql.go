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
		Logger:                                   logger.Default.LogMode(logger.Warn),
		DisableForeignKeyConstraintWhenMigrating: true,
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

	// 先迁移旧数据，再 AutoMigrate，避免 user_id 外键约束失败
	if err := migrateLegacyVideos(db); err != nil {
		return nil, fmt.Errorf("legacy migrate: %w", err)
	}
	if err := backfillOrphanVideoUserIDs(db); err != nil {
		return nil, fmt.Errorf("backfill videos: %w", err)
	}

	if err := db.AutoMigrate(
		&entity.User{},
		&entity.Video{},
		&entity.VideoLike{},
		&entity.VideoFavorite{},
		&entity.VideoShare{},
		&entity.VideoComment{},
		&entity.Follow{},
		&entity.Conversation{},
		&entity.ConversationMember{},
		&entity.Message{},
		&entity.Notification{},
	); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
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
		return backfillFromLegacyUsername(db)
	}

	type legacyVideo struct {
		ID       int64
		Username string
		Avatar   string
	}
	var rows []legacyVideo
	if err := db.Table("videos").Select("id, username, avatar").Find(&rows).Error; err != nil {
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

func backfillFromLegacyUsername(db *gorm.DB) error {
	if !db.Migrator().HasColumn(&entity.Video{}, "username") {
		return nil
	}
	type legacyVideo struct {
		ID       int64
		Username string
		Avatar   string
	}
	var rows []legacyVideo
	if err := db.Table("videos").Where("user_id = 0 OR user_id IS NULL").Select("id, username, avatar").Find(&rows).Error; err != nil {
		return err
	}
	for _, row := range rows {
		var u entity.User
		err := db.Where("username = ?", row.Username).First(&u).Error
		if err != nil {
			u = entity.User{
				Username: row.Username,
				Nickname: row.Username,
				Avatar:   row.Avatar,
				DouyinID: fmt.Sprintf("%d", 10000000+row.ID*1234567),
			}
			if err := db.Create(&u).Error; err != nil {
				return err
			}
		}
		if err := db.Table("videos").Where("id = ?", row.ID).Update("user_id", u.ID).Error; err != nil {
			return err
		}
	}
	return nil
}

func backfillOrphanVideoUserIDs(db *gorm.DB) error {
	if !db.Migrator().HasTable(&entity.Video{}) {
		return nil
	}
	var orphan int64
	if err := db.Table("videos").Where("user_id = 0 OR user_id IS NULL").Count(&orphan).Error; err != nil {
		return err
	}
	if orphan == 0 {
		return nil
	}

	var defaultUser entity.User
	err := db.Order("id ASC").First(&defaultUser).Error
	if err != nil {
		defaultUser = entity.User{
			Username: "my_account",
			Nickname: "我的账号",
			Avatar:   "https://images.unsplash.com/photo-1728046666898-7e42ed206c9f?w=200&h=200&fit=crop",
			DouyinID: "12345678",
			Bio:      "记录生活，分享美好 ✨",
		}
		if err := db.Create(&defaultUser).Error; err != nil {
			return err
		}
	}
	return db.Table("videos").Where("user_id = 0 OR user_id IS NULL").Update("user_id", defaultUser.ID).Error
}
