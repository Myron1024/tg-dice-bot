package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"log"
	"os"
	"tg-dice-bot/internal/database"
	"tg-dice-bot/internal/model"

	"github.com/go-redis/redis/v8"
)

const (
	TelegramAPIToken = "TELEGRAM_API_TOKEN"
)

var (
	db      *gorm.DB
	redisDB *redis.Client
)

func StartBot() {
	initDB()

	bot := initTelegramBot()

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			go handleMessage(bot, update.Message)
		} else if update.CallbackQuery != nil {
			go handleCallbackQuery(bot, update.CallbackQuery)
		}
	}
}

func initDB() {
	var err error
	db, err = database.InitDB(os.Getenv(database.DBConnectionString))
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	err = db.AutoMigrate(&model.LotteryRecord{})
	if err != nil {
		log.Fatal("自动迁移表结构失败:", err)
	}

	err = db.AutoMigrate(&model.TgUser{})
	if err != nil {
		log.Fatal("自动迁移表结构失败:", err)
	}

	err = db.AutoMigrate(&model.BetRecord{})
	if err != nil {
		log.Fatal("自动迁移表结构失败:", err)
	}

	err = db.AutoMigrate(&model.ChatDiceConfig{})
	if err != nil {
		log.Fatal("自动迁移表结构失败:", err)
	}

	redisDB, err = database.InitRedisDB(os.Getenv(database.RedisDBConnectionString))
	if err != nil {
		log.Fatal("连接Redis数据库失败:", err)
	}

}
func initTelegramBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(os.Getenv(TelegramAPIToken))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("已授权帐户 %s", bot.Self.UserName)
	return bot
}
