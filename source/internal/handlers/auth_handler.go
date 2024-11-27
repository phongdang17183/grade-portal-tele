package handlers

import (
	"Grade_Portal_TelegramBot/internal/services"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


func HandleRegister(bot *tgbotapi.BotAPI, update tgbotapi.Update, input string) {
	parts := strings.Split(input, " ")
	var mssv, pw, otp string
	mssv, pw, otp = parts[0], parts[1], parts[2]
	resp, err := services.RegisterStudent(mssv, pw, otp)
	var response string
	if err == nil {
		response = resp.Msg + ", vui lòng login bằng cú pháp /login_mssv_password để sử dụng dịch vụ."
	} else {
		response = "Error fetching student info:" + err.Error()
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}

func HandleOTP(bot *tgbotapi.BotAPI, update tgbotapi.Update, mssv string) {
	_, err := services.GetOTP(mssv)
	var response string
	if err == nil {
		response = "OTP đã được gửi về email của bạn, vui kiểm tra email."
	} else {
		response = err.Error()
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}

func HanldeLogin(bot *tgbotapi.BotAPI, update tgbotapi.Update, input string) {
	parts := strings.Split(input, " ")
	var mssv, pw string
	mssv, pw = parts[0], parts[1]
	resp, err := services.Login(update.Message.Chat.ID, mssv, pw)
	var response string
	if err == nil {
		response = "Đăng nhập thành công, các khóa học bạn đang có là: " + strings.Join(resp.ListCourse, ", ")
	} else {
		response = "Có lỗi trong việc xác thực hãy thử lại sau: " + err.Error()
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}

