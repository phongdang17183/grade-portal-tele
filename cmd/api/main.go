package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Lấy BOT_TOKEN từ biến môi trường
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN is not set")
	}

	// Khởi tạo bot
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err) // In lỗi chi tiết
	}

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
		if update.Message == nil { // nếu không có tin nhắn, bỏ qua
			continue
		}

		log.Printf("[%s] %s\n", update.Message.From.UserName, update.Message.Text)

		// Phản hồi với người dùng
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Bạn đã gửi: "+update.Message.Text + "\nGửi cái đéo gì nữa, HP không biết làm cái gì cả :(")
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
	}
}