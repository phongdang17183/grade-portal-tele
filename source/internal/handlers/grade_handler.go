package handlers

import (
	"Grade_Portal_TelegramBot/internal/services"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleUpdate chính: xác định và chuyển lệnh đến các hàm xử lý riêng
func HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		HandleStart(bot, update)
	case "getOTP":
		HandleOTP(bot, update, update.Message.CommandArguments())
	case "register":
		HandleRegister(bot, update, update.Message.CommandArguments(), update.Message.CommandArguments(), update.Message.CommandArguments())
	case "resetPassWord":
		HandleRegister(bot, update, update.Message.CommandArguments(), update.Message.CommandArguments(), update.Message.CommandArguments())
	case "login":
		HanldeLogin(bot, update, update.Message.CommandArguments(), update.Message.CommandArguments())
	case "help":
		HandleHelp(bot, update)
	case "info":
		HandleInfo(bot, update)
	case "grade":
		HandleGrade(bot, update, update.Message.CommandArguments())
	case "allGrade":
		HandleAllGrade(bot, update)
	case "clear":
		HandleClear(bot, update)
	case "history":
		HandleHistory(bot, update)
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Lệnh không hợp lệ. Dùng /help để xem danh sách lệnh.")
		bot.Send(msg)
	}
}

func HandleStart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.From.ID
	print(userID)
	response := fmt.Sprintf("Chào mừng %d tôi là hệ thống tra cứu điểm - một bot-chat hỗ trợ tra cứu điểm nhanh chóng!\n\n"+
		"Hướng dẫn sử dụng: Đăng nhập qua lệnh /register + [MSSV]. Một số lệnh bạn có thể dùng:\n"+
		"/grade - tra cứu điểm\n/history- xem lịch sử điểm\n/help - để biết thêm các lệnh khác.",
		userID)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}

func HandleRegister(bot *tgbotapi.BotAPI, update tgbotapi.Update, mssv string, pw string, otp string) {
	success := services.RegisterStudent(update.Message.Chat.ID, mssv, pw, otp)
	var response string
	if success {
		response = "Bạn đã đăng ký thành công, vui lòng login bằng cú pháp /login_mssv_password để sử dụng dịch vụ."
	} else {
		response = "Có lỗi trong việc đăng ký hãy thử lại sau"
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}
func HandleOTP(bot *tgbotapi.BotAPI, update tgbotapi.Update, mssv string) {
	success := services.GetOTP(update.Message.Chat.ID, mssv)
	var response string
	if success {
		response = "OTP đã được gửi về email của bạn, vui kiểm tra email."
	} else {
		response = "Có lỗi trong việc xác thực hãy thử lại sau"
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}

func HanldeLogin(bot *tgbotapi.BotAPI, update tgbotapi.Update, mssv string, pw string) {
	success := services.Login(update.Message.Chat.ID, mssv, pw)
	var response string
	if success {
		response = "Đăng nhập thành công, bạn có thể sử dụng các chức năng của bot"
	} else {
		response = "Có lỗi trong việc xác thực hãy thử lại sau"
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}

func HandleHelp(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	helpMessage := ` Danh sách lệnh hỗ trợ: 
	/start - Bắt đầu sử dụng bot và nhận hướng dẫn 
	/register <email> - Đăng ký tài khoản với email của bạn hệ thống sẽ gửi bạn otp
	/login <otp> - Đăng nhập với otp đã gửi 
	/info - Xem thông tin tài khoản của bạn 
	/grade <semester or course ID> - Tra cứu điểm theo học kỳ hoặc mã môn học 
	/clear - Xóa lịch sử tra cứu điểm 
	/history - Xem lịch sử tra cứu điểm 
	/help - Xem danh sách lệnh hỗ trợ này `
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpMessage)
	bot.Send(msg)
}

func HandleInfo(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userInfo, err := services.GetStudentInfo(update.Message.Chat.ID)
	if err != nil {
		response := "Không tìm thấy thông tin đăng nhập."
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		bot.Send(msg)
		return
	}
	response := fmt.Sprintf("Thông tin đăng nhập\n________\nHọ và tên: %s\nMSSV: %s", userInfo.Name, userInfo.StudentID)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}

func HandleGrade(bot *tgbotapi.BotAPI, update tgbotapi.Update, semesterOrCourseID string) {
	grades, err := services.GetGrades(update.Message.Chat.ID, semesterOrCourseID)
	var response string
	if err != nil {
		response = "Không thể lấy dữ liệu điểm."
	} else {
		response = fmt.Sprintf("Kết quả điểm cho %s:\n________\n", semesterOrCourseID)
		for _, grade := range grades {
			response += fmt.Sprintf("%s: %.1f\n", grade.CourseName, grade.Score)
		}
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}

func HandleAllGrade(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	grades, err := services.GetAllGrades()
	var response string
	if err != nil {
		response = "Không thể lấy dữ liệu điểm."
	} else {
		response = fmt.Sprintf("Kết quả điểm cho %s:\n________\n", semesterOrCourseID)
		for _, grade := range grades {
			response += fmt.Sprintf("%s: %.1f\n", grade.CourseName, grade.Score)
		}
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}

func HandleClear(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	services.ClearHistory(update.Message.Chat.ID)
	response := "Lịch sử tra cứu đã được xóa."
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}

func HandleHistory(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	history, err := services.GetHistory(update.Message.Chat.ID)
	var response string
	if err != nil {
		response = "Không có lịch sử tra cứu nào."
	} else {
		response = "Lịch sử tra cứu:\n"
		for _, entry := range history {
			response += fmt.Sprintf("%s: %.1f\n", entry.CourseName, entry.Score)
		}
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}
