package handlers

import (
	"Grade_Portal_TelegramBot/config"
	"Grade_Portal_TelegramBot/internal/services"
	"fmt"
	"log"
	"regexp"
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
	// Thêm cơ chế bắt lỗi toàn cục để ngăn bot bị crash (Panic + Runtime Error)
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic caught: %v", r)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Đã xảy ra lỗi không mong muốn. Vui lòng thử lại sau. :(")
			bot.Send(msg)
		}
	}()

	// Tách input và kiểm tra đủ tham số
	parts := strings.Fields(input) // Xử lý cả chuỗi có nhiều khoảng trắng
	if len(parts) < 2 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Có vẻ như bạn chưa nhập MSSV hoặc mật khẩu. Vui lòng nhập đúng cú pháp: /login [MSSV] [mật khẩu].")
		bot.Send(msg)
		return
	}

	// Xác thực dữ liệu nhập (MSSV là chuỗi số, mật khẩu không chứa khoảng trắng).
	mssv, pw := parts[0], strings.Join(parts[1:], " ")
	isValidMSSV := regexp.MustCompile(`^\d{7}$`).MatchString(mssv)
	if !isValidMSSV {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "MSSV không hợp lệ. Vui lòng nhập MSSV gồm 7 chữ số.")
		bot.Send(msg)
		return
	}
	if strings.Contains(pw, " ") {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Mật khẩu không được chứa khoảng trắng. Vui lòng nhập lại.")
		bot.Send(msg)
		return
	}

	// Gửi yêu cầu đăng nhập qua API
	resp, err := services.Login(update.Message.Chat.ID, mssv, pw, cfg)
	var response string
	if err != nil {
		// Log lỗi để hỗ trợ debug
		log.Printf("Login error for MSSV %s: %v", mssv, err)

		// Xử lý các loại lỗi
		if strings.Contains(err.Error(), "Timeout") {
			response = "API không phản hồi, vui lòng thử lại sau."
		} else if strings.Contains(err.Error(), "HTTP 400") {
			response = "Sai MSSV hoặc mật khẩu. Vui lòng kiểm tra và thử lại."
		} else if strings.Contains(err.Error(), "unexpected status code") {
			response = "Hệ thống đang gặp sự cố, vui lòng thử lại sau."
		} else {
			response = fmt.Sprintf("Đăng nhập thất bại. Chi tiết lỗi: %s", err.Error())
		}
	} else {
		// Đăng nhập thành công
		response = "Đăng nhập thành công! Các khóa học bạn đang có là: " + strings.Join(resp.ListCourse, ", ")
	}

	// Gửi phản hồi đến người dùng
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}
