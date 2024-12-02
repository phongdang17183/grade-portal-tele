package handlers

import (
	"Grade_Portal_TelegramBot/config"
	"Grade_Portal_TelegramBot/internal/services"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleRegister(bot *tgbotapi.BotAPI, update tgbotapi.Update, input string, cfg *config.Config) {
	parts := strings.Split(input, " ")
	var mssv, pw, otp string
	var response string
	if len(parts) < 3 {
		response = "Thiếu MSSV, password hoặc OTP."
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		bot.Send(msg)
		return
	}
	mssv, pw, otp = parts[0], parts[1], parts[2]

	resp, err := services.RegisterStudent(mssv, pw, otp, cfg)
	if err != nil {
		if strings.Contains(err.Error(), "error encoding JSON") {
			response = "Hệ thống gặp sự cố. Hãy thử lại vào lần sau."
		} else if strings.Contains(err.Error(), "error creating request") {
			response = "Không kết nối được với hệ thống. Hãy thử lại vào lần sau."
		} else if strings.Contains(err.Error(), "error sending request") {
			response = "Hệ thống không phản hồi. Hãy thử lại vào lần sau."
		} else if strings.Contains(err.Error(), "unexpected status code") {
			response = "Hệ thống gặp lỗi khi truy xuất thông tin."
		} else if strings.Contains(err.Error(), "error decoding response") {
			response = "Dữ liệu nhận được không hợp lệ. Hãy thử lại vào lần sau."
		} else {
			response = "Đã xảy ra lỗi không xác định. Hãy thử lại vào lần sau."
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response+"Lỗi: "+err.Error())
		bot.Send(msg)
		return
	}

	response = resp.Msg + ", vui lòng login bằng cú pháp /login_mssv_password để sử dụng dịch vụ."
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}

func HandleOTP(bot *tgbotapi.BotAPI, update tgbotapi.Update, mssv string, cfg *config.Config) {
	var response string
	if mssv == "" {
		response = "Thiếu MSSV."
	} else {
		_, err := services.GetOTP(mssv, cfg)
		if err == nil {
			response = "OTP đã được gửi về email của bạn, vui kiểm tra email."
		} else {
			response = "Có lỗi trong việc lấy OTP, vui lòng thử lại sau: " + err.Error() + "\n"
			// response = err.Error()
		}
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}

func HanldeLogin(bot *tgbotapi.BotAPI, update tgbotapi.Update, input string, cfg *config.Config) {
	parts := strings.Split(input, " ")
	var mssv, pw string
	var response string
	if len(parts) < 2 {
		response = "Thiếu MSSV hoặc password."
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		bot.Send(msg)
		return
	}
	mssv, pw = parts[0], parts[1]
	resp, err := services.Login(update.Message.Chat.ID, mssv, pw, cfg)
	if err == nil {
		response = "Đăng nhập thành công, các khóa học bạn đang có là: " + strings.Join(resp.ListCourse, ", ")
	} else {
		response = "Có lỗi trong việc xác thực hãy thử lại sau: "
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}
