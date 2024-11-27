package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"Grade_Portal_TelegramBot/internal/services"
	"fmt"
)


func HandleStart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.From.ID
	response := fmt.Sprintf("Chào mừng %d tôi là hệ thống tra cứu điểm - một bot-chat hỗ trợ tra cứu điểm nhanh chóng!\n\n"+
		"Hướng dẫn sử dụng: Đăng nhập qua lệnh /login + [MSSV] + [password]\n" +
		"/grade + [Mã học phần] - tra cứu điểm \n" +
		"/allGrade - xem tất cả điểm của bạn \n"+
		"/history - xem lịch sử điểm\n"+
		"/clear - xóa lịch sử điểm\n"+
		"/info - xem thông tin tài khoản\n"+
		"/getOTP - lấy OTP để đăng nhập\n"+
		"/register [MSSV] [password] [OTP] - đăng ký tài khoản\n"+
		"/resetPassWord [MSSV] [password] [OTP] - đổi mật khẩu\n"+
		"/help - để biết thêm các lệnh khác." ,
		userID)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}

func HandleHelp(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	response := fmt.Sprintf(
		"Hướng dẫn sử dụng: Đăng nhập qua lệnh /login + [MSSV] + [password]\n" +
		"/grade - tra cứu điểm \n" +
		"/allGrade - xem tất cả điểm của bạn \n"+
		"/history - xem lịch sử điểm\n"+
		"/clear - xóa lịch sử điểm\n"+
		"/info - xem thông tin tài khoản\n"+
		"/getOTP - lấy OTP để đăng nhập\n"+
		"/register - đăng ký tài khoản\n"+
		"/resetPassWord - đổi mật khẩu\n"+
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
	_, err := services.GetHistory(update.Message.Chat.ID)
	var response string
	if err != nil {
		response = "Không có lịch sử tra cứu nào."
	} else {
		response = "res"
	}
	//fmt.Println(res)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}