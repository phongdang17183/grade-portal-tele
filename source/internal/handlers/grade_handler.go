package handlers

import (
	"Grade_Portal_TelegramBot/internal/services"
	"encoding/json"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleUpdate chính: xác định và chuyển lệnh đến các hàm xử lý riêng
func HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		HandleStart(bot, update)
	case "register":
		HandleRegister(bot, update, update.Message.CommandArguments())
	case "help":
		HandleHelp(bot, update)
	case "info":
		HandleInfo(bot, update)
	case "grade":
		HandleGrade(bot, update, update.Message.CommandArguments())
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
	response := fmt.Sprintf("Chào mừng %d tôi là hệ thống tra cứu điểm - một bot-chat hỗ trợ tra cứu điểm nhanh chóng!\n\n"+
		"Hướng dẫn sử dụng: Đăng nhập qua lệnh /register + [MSSV]. Một số lệnh bạn có thể dùng:\n"+
		"/grade - tra cứu điểm\n/history- xem lịch sử điểm\n/help - để biết thêm các lệnh khác.",
		userID)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}

func HandleRegister(bot *tgbotapi.BotAPI, update tgbotapi.Update, studentID string) {
	success := services.RegisterStudent(update.Message.Chat.ID, studentID)
	var response string
	if success {
		response = "Chào mừng đến với hệ thống."
	} else {
		response = "Tài khoản Telegram này đã đăng ký với MSSV khác. Đăng ký không thành công.jyhht"
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}

func HandleHelp(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	helpMessage := ` Danh sách lệnh hỗ trợ: 
	/start - Bắt đầu sử dụng bot và nhận hướng dẫn 
	/register <MSSV> - Đăng ký tài khoản với mã số sinh viên (MSSV) của bạn 
	/info - Xem thông tin tài khoản của bạn 
	/grade <semester or course ID> - Tra cứu điểm theo học kỳ hoặc mã môn học 
	/clear - Xóa lịch sử tra cứu điểm 
	/history - Xem lịch sử tra cứu điểm 
	/help - Xem danh sách lệnh hỗ trợ này `
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpMessage)
	bot.Send(msg)
}

func HandleInfo(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	resp, err := services.GetStudentInfo(update.Message.Chat.ID)
	errorDetails := "Lỗi: " + err.Error()

	if err != nil {
		var response string
		if strings.Contains(err.Error(), "token not found") {
			response = "Không tìm thấy thông tin đăng nhập. Hãy đăng nhập trước khi sử dụng dịch vụ"
		} else if strings.Contains(err.Error(), "database error") {
			response = "Không kết nối được với cơ sở dữ liệu. Hãy thử lại vào lần sau."
		} else if strings.Contains(err.Error(), "error getting token") {
			response = "Không tìm thấy thông tin đăng nhập. Hãy đăng nhập trước khi sử dụng dịch vụ"
		} else if strings.Contains(err.Error(), "error creating request") {
			response = "Không  kết nối được với hệ thống. Hãy thử lại vào lần sau"
		} else if strings.Contains(err.Error(), "error sending request") {
			response = "Hệ thống không phản hồi.  Hãy thử lại vào lần sau."
		} else if strings.Contains(err.Error(), "unexpected status code") {
			response = "Hệ thống gặp lỗi khi truy xuất thông tin. Mã lỗi API không hợp lệ."
		} else if strings.Contains(err.Error(), "error decoding response") {
			response = "Dữ liệu nhận được không hợp lệ. Hãy thử lại vào lần sau."
		} else if strings.Contains(err.Error(), "access forbidden") {
			response = "Hệ thống từ chối yêu cầu. Hãy liên hệ với dịch vụ hỗ trợ."
		} else if strings.Contains(err.Error(), "internal server error") {
			response = "Lỗi hệ thống. Hãy thử lại vào lần sau."
		} else if strings.Contains(err.Error(), "unauthorized access") {
			response = "Không có quyền truy cập. Hãy kiểm tra thông tin đăng nhập."
		} else if strings.Contains(err.Error(), "timeout") {
			response = "Kết nối bị gián đoạn. Hãy thử lại vào lần sau."
		} else {
			response = "Đã xảy ra lỗi không xác định. Hãy thử lại vào lần sau."
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response+"\n\n"+errorDetails)
		bot.Send(msg)
		return
	}

	response, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		fmt.Println(err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Lỗi xử lý dữ liệu.")
		bot.Send(msg)
		return
	}
	msgText := fmt.Sprintf("```json\n%s\n```", response)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	msg.ParseMode = "MarkdownV2"
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
