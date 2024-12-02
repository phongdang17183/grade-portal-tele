package handlers

import (
	"Grade_Portal_TelegramBot/internal/services"
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleStart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.From.ID
	response := fmt.Sprintf("Chào mừng %d tôi là hệ thống tra cứu điểm - một bot-chat hỗ trợ tra cứu điểm nhanh chóng!\n\n"+
		"Hướng dẫn sử dụng: Đăng nhập qua lệnh /login + [MSSV] + [password]\n"+
		"/grade + [Mã học phần] - tra cứu điểm \n"+
		"/allgrade - xem tất cả điểm của bạn \n"+
		"/history - xem lịch sử điểm\n"+
		"/clear - xóa lịch sử điểm\n"+
		"/info - xem thông tin tài khoản\n"+
		"/getotp + [MSSV] - lấy OTP để đăng ký hoặc đổi mật khẩu\n"+
		"/register + [MSSV] + [password] + [OTP] - đăng ký tài khoản\n"+
		"/resetpassword + [MSSV] + [password] + [OTP] - đổi mật khẩu\n"+
		"/help - để biết thêm các lệnh khác.",
		userID)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}

func HandleHelp(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	response := fmt.Sprintf(
		"Hướng dẫn sử dụng: Đăng nhập qua lệnh /login + [MSSV] + [password]\n" +
			"/grade - tra cứu điểm \n" +
			"/allgrade - xem tất cả điểm của bạn \n" +
			"/history - xem lịch sử điểm\n" +
			"/clear - xóa lịch sử điểm\n" +
			"/info - xem thông tin tài khoản\n" +
			"/getotp + [MSSV] - lấy OTP để đăng nhập\n" +
			"/register + [MSSV] + [password] + [OTP] - đăng ký tài khoản\n" +
			"/resetpassword + [MSSV] + [password] + [OTP] - đổi mật khẩu\n" +
			"/help - để biết thêm các lệnh khác.")
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}

func HandleClear(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	res := services.ClearHistory(update.Message.Chat.ID)
	var response string
	if res {
		response = "Lịch sử tra cứu đã được xóa."
	} else {
		response = "Error"
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
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
