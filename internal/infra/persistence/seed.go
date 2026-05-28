package persistence

import (
	"gorm.io/gorm"

	"github.com/LoveMuZiLi/tiktok-backend/internal/entity"
)

func Seed(db *gorm.DB) error {
	var userCount int64
	if err := db.Model(&entity.User{}).Count(&userCount).Error; err != nil {
		return err
	}
	if userCount > 0 {
		return nil
	}

	users := []entity.User{
		{Username: "my_account", Nickname: "我的账号", Avatar: "https://images.unsplash.com/photo-1728046666898-7e42ed206c9f?w=200&h=200&fit=crop", DouyinID: "12345678", Bio: "记录生活，分享美好 ✨"},
		{Username: "fashion_daily", Nickname: "时尚日记", Avatar: "https://images.unsplash.com/photo-1728046666898-7e42ed206c9f?w=100&h=100&fit=crop", DouyinID: "11111111", Bio: "每日穿搭分享"},
		{Username: "beauty_tips", Nickname: "美妆技巧", Avatar: "https://images.unsplash.com/photo-1583318605147-8e52610d9c75?w=100&h=100&fit=crop", DouyinID: "22222222", Bio: "新手也能学会的妆"},
		{Username: "photo_artist", Nickname: "摄影艺术家", Avatar: "https://images.unsplash.com/photo-1676288785587-0d4398fbf38e?w=100&h=100&fit=crop", DouyinID: "33333333", Bio: "记录美好瞬间"},
		{Username: "urban_life", Nickname: "城市生活", Avatar: "https://images.unsplash.com/photo-1733473571611-2cf5460d91fc?w=100&h=100&fit=crop", DouyinID: "44444444", Bio: "城市夜景爱好者"},
		{Username: "nature_lover", Nickname: "自然爱好者", Avatar: "https://images.unsplash.com/photo-1728046666871-7ff531542fb1?w=100&h=100&fit=crop", DouyinID: "55555555", Bio: "大自然的治愈力量"},
		{Username: "kkkkkk_y", Nickname: "kkkkkk_y", Avatar: "https://images.unsplash.com/photo-1583318605147-8e52610d9c75?w=100&h=100&fit=crop", DouyinID: "66666666", Bio: ""},
		{Username: "DramaMindFactory", Nickname: "DramaMindFactory", Avatar: "https://images.unsplash.com/photo-1676288785587-0d4398fbf38e?w=100&h=100&fit=crop", DouyinID: "77777777", Bio: ""},
	}
	if err := db.Create(&users).Error; err != nil {
		return err
	}

	videos := []entity.Video{
		{UserID: users[1].ID, Description: "今天的穿搭分享 🌸 喜欢记得点赞哦~", Music: "原创音乐 - Fashion Daily", Image: "https://images.unsplash.com/photo-1728046666898-7e42ed206c9f?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&q=80&w=1080", Likes: 128500, Comments: 2341, Shares: 892},
		{UserID: users[2].ID, Description: "超简单的妆容教程 💄 新手也能学会！", Music: "热门音乐 - Beauty Vibes", Image: "https://images.unsplash.com/photo-1583318605147-8e52610d9c75?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&q=80&w=1080", Likes: 256700, Comments: 4523, Shares: 1634},
		{UserID: users[3].ID, Description: "记录生活的美好瞬间 📷✨", Music: "温柔吉他 - Chill Beats", Image: "https://images.unsplash.com/photo-1676288785587-0d4398fbf38e?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&q=80&w=1080", Likes: 89300, Comments: 1876, Shares: 567},
		{UserID: users[4].ID, Description: "城市夜景太美了 🌃 你最喜欢哪座城市？", Music: "City Lights - Urban Sounds", Image: "https://images.unsplash.com/photo-1733473571611-2cf5460d91fc?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&q=80&w=1080", Likes: 445600, Comments: 7890, Shares: 2341},
		{UserID: users[5].ID, Description: "大自然的治愈力量 🌿 放慢脚步，享受当下", Music: "自然之声 - Peaceful Moments", Image: "https://images.unsplash.com/photo-1728046666871-7ff531542fb1?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&q=80&w=1080", Likes: 167800, Comments: 3245, Shares: 1123},
		{UserID: users[0].ID, Description: "我的第一条作品 🎬", Music: "原创音乐", Image: "https://images.unsplash.com/photo-1728046666898-7e42ed206c9f?w=400&h=600&fit=crop", Likes: 12800, Comments: 320, Shares: 88},
	}
	if err := db.Create(&videos).Error; err != nil {
		return err
	}

	follows := []entity.Follow{
		{FollowerID: users[0].ID, FollowingID: users[1].ID},
		{FollowerID: users[0].ID, FollowingID: users[2].ID},
		{FollowerID: users[0].ID, FollowingID: users[3].ID},
		{FollowerID: users[1].ID, FollowingID: users[0].ID},
		{FollowerID: users[2].ID, FollowingID: users[0].ID},
		{FollowerID: users[6].ID, FollowingID: users[0].ID},
		{FollowerID: users[7].ID, FollowingID: users[0].ID},
	}
	if err := db.Create(&follows).Error; err != nil {
		return err
	}

	conv := entity.Conversation{}
	if err := db.Create(&conv).Error; err != nil {
		return err
	}
	members := []entity.ConversationMember{
		{ConversationID: conv.ID, UserID: users[0].ID, UnreadCount: 2},
		{ConversationID: conv.ID, UserID: users[6].ID},
	}
	if err := db.Create(&members).Error; err != nil {
		return err
	}
	messages := []entity.Message{
		{ConversationID: conv.ID, SenderID: users[6].ID, Text: "你好！看到你的视频了，拍得真棒！"},
		{ConversationID: conv.ID, SenderID: users[0].ID, Text: "谢谢！😊"},
		{ConversationID: conv.ID, SenderID: users[6].ID, Text: "能分享一下拍摄技巧吗？"},
		{ConversationID: conv.ID, SenderID: users[0].ID, Text: "主要是光线和角度，找自然光好的地方很重要"},
	}
	if err := db.Create(&messages).Error; err != nil {
		return err
	}

	notifications := []entity.Notification{
		{UserID: users[0].ID, Type: entity.NotificationTypeFollow, Title: "新粉丝", Body: "DramaMindFactory 关注了你。", IsRead: false},
		{UserID: users[0].ID, Type: entity.NotificationTypeActivity, Title: "活动", Body: "Im andreiaaa!! 查看了你的主页。", IsRead: false},
		{UserID: users[0].ID, Type: entity.NotificationTypeSystem, Title: "系统通知", Body: "直播: 观众希望看到你的更多内容 · 3天", IsRead: false},
	}
	return db.Create(&notifications).Error
}
