package persistence

import (
	"gorm.io/gorm"

	"github.com/LoveMuZiLi/tiktok-backend/internal/entity"
)

func Seed(db *gorm.DB) error {
	var count int64
	if err := db.Model(&entity.Video{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	rows := []entity.Video{
		{
			Username: "fashion_daily",
			Avatar:   "https://images.unsplash.com/photo-1728046666898-7e42ed206c9f?w=100&h=100&fit=crop",
			Description: "今天的穿搭分享 🌸 喜欢记得点赞哦~",
			Music:    "原创音乐 - Fashion Daily",
			Image:    "https://images.unsplash.com/photo-1728046666898-7e42ed206c9f?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&q=80&w=1080",
			Likes: 128500, Comments: 2341, Shares: 892,
		},
		{
			Username: "beauty_tips",
			Avatar:   "https://images.unsplash.com/photo-1583318605147-8e52610d9c75?w=100&h=100&fit=crop",
			Description: "超简单的妆容教程 💄 新手也能学会！",
			Music:    "热门音乐 - Beauty Vibes",
			Image:    "https://images.unsplash.com/photo-1583318605147-8e52610d9c75?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&q=80&w=1080",
			Likes: 256700, Comments: 4523, Shares: 1634,
		},
		{
			Username: "photo_artist",
			Avatar:   "https://images.unsplash.com/photo-1676288785587-0d4398fbf38e?w=100&h=100&fit=crop",
			Description: "记录生活的美好瞬间 📷✨",
			Music:    "温柔吉他 - Chill Beats",
			Image:    "https://images.unsplash.com/photo-1676288785587-0d4398fbf38e?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&q=80&w=1080",
			Likes: 89300, Comments: 1876, Shares: 567,
		},
		{
			Username: "urban_life",
			Avatar:   "https://images.unsplash.com/photo-1733473571611-2cf5460d91fc?w=100&h=100&fit=crop",
			Description: "城市夜景太美了 🌃 你最喜欢哪座城市？",
			Music:    "City Lights - Urban Sounds",
			Image:    "https://images.unsplash.com/photo-1733473571611-2cf5460d91fc?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&q=80&w=1080",
			Likes: 445600, Comments: 7890, Shares: 2341,
		},
		{
			Username: "nature_lover",
			Avatar:   "https://images.unsplash.com/photo-1728046666871-7ff531542fb1?w=100&h=100&fit=crop",
			Description: "大自然的治愈力量 🌿 放慢脚步，享受当下",
			Music:    "自然之声 - Peaceful Moments",
			Image:    "https://images.unsplash.com/photo-1728046666871-7ff531542fb1?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&q=80&w=1080",
			Likes: 167800, Comments: 3245, Shares: 1123,
		},
	}
	return db.Create(&rows).Error
}
