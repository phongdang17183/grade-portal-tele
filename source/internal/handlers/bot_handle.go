package handlers

import (
	"Grade_Portal_TelegramBot/internal/services"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var commands = []struct {
	Command     string
	Description string
}{
	{"/login + [MSSV] + [password]", "Đăng nhập vào hệ thống"},
	{"/grade + [Mã học phần]", "Tra cứu điểm"},
	{"/allgrade", "Xem tất cả điểm"},
	{"/history", "Xem lịch sử điểm"},
	{"/clear", "Xóa lịch sử điểm"},
	{"/info", "Xem thông tin tài khoản"},
	{"/getotp + [MSSV]", "Lấy OTP"},
	{"/register + [MSSV] + [password] + [OTP]", "Đăng ký tài khoản"},
	{"/resetpassword + [MSSV] + [password] + [OTP]", "Đổi mật khẩu"},
	{"/help", "Để biết thêm các lệnh khác"},
}

func HandleStart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.From.ID
	escapedUserID := fmt.Sprintf("\\%d", userID)
	var response strings.Builder

	response.WriteString(fmt.Sprintf("Chào mừng *%s* đến với hệ thống tra cứu điểm, tôi là một bot-chat hỗ trợ tra cứu điểm nhanh chóng!🎉\n\n", escapedUserID))
	response.WriteString("*Hướng dẫn:*\n\n")

	for i, cmd := range commands {
		response.WriteString(fmt.Sprintf("%d\\. `%s` \\- %s\n", i+1, cmd.Command, cmd.Description))
	}

	imagePath := "img/Hello.png"
	photo := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FilePath(imagePath))
	if _, err := bot.Send(photo); err != nil {
		log.Println("Lỗi gửi ảnh:", err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response.String())
	msg.ParseMode = "Markdown"
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Lỗi khi gửi tin nhắn trong: %v", err)
	}
}
func HandleHelp(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var response strings.Builder
	response.WriteString("*Hướng dẫn sử dụng:*\n\n")
	for i, cmd := range commands {
		response.WriteString(fmt.Sprintf("%d\\. `%s` \\- %s\n", i+1, cmd.Command, cmd.Description))
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response.String())
	msg.ParseMode = "MarkdownV2"
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Lỗi khi gửi tin nhắn trong: %v", err)
	}
}

func HandleClear(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	res := services.ClearHistory(update.Message.Chat.ID)
	var response string
	if res {
		response = "Lịch sử tra cứu đã được xóa ✅\\."
	} else {
		response = "Error ❌\\."
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	msg.ParseMode = "MarkdownV2"
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Lỗi khi gửi tin nhắn: %v", err)
	}
}

func HandleHistory(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	response, err := services.GetHistory(update.Message.Chat.ID)
	var msg tgbotapi.MessageConfig
	if err != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Không có lịch sử tra cứu nào.")
		bot.Send(msg)
	} else {
		for _, course := range *response {
			jsonStr, _ := json.Marshal(course)
			fmt.Println(string(jsonStr))
			msgText := fmt.Sprintf("```json\n%s\n```", string(jsonStr))
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
			msg.ParseMode = "MarkdownV2"
			bot.Send(msg)
		}
	}
}
