// cmd/api/other.go
package bot

import (
	"Grade_Portal_TelegramBot/config"
	"Grade_Portal_TelegramBot/internal/handlers"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func Start() {
	// Dùng hàm LoadConfig() định nghĩa trong file config.go - package config.
	cfg := config.LoadConfig()
	// Khởi tạo bot
	bot, err := tgbotapi.NewBotAPI(cfg.BOT_TOKEN)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err) // In lỗi chi tiết
	}
	// connet DBMongo
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DBURL))
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal("Failed to disconnect: %v ", err)
		}
	}()

	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	// Cấu hình để nhận cập nhật
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10

	// Nhận cập nhật từ bot
	fmt.Println("Starting bot...")
	updates := bot.GetUpdatesChan(u)
	fmt.Println("Listening for updates...")

	// Vòng lặp để xử lý các bản cập nhật
	for update := range updates {
		if update.Message != nil {
			handlers.HandleUpdate(bot, update)
		}
	}

}
