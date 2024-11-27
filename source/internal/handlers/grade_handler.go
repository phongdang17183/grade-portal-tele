package handlers

import (
	"Grade_Portal_TelegramBot/internal/services"
	"fmt"
	"strings"

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

func HanldeLogin(bot *tgbotapi.BotAPI, update tgbotapi.Update, mssv string, pw string) {
	resp, err := services.Login(update.Message.Chat.ID, mssv, pw)
	var response string
	if err == nil {
		response = "Đăng nhập thành công, các khóa học bạn đang có là: " + strings.Join(resp.ListCourse, ", ")
	} else {
		response = "Có lỗi trong việc xác thực hãy thử lại sau" + err.Error()
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
	resp, err := services.GetStudentInfo(update.Message.Chat.ID)
	if err != nil {
		response := "Không tìm thấy thông tin đăng nhập."
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		bot.Send(msg)
		return
	}
	response := fmt.Sprintf("Thông tin đăng nhập\n________\nHọ và tên: %s\nMSSV: %s %s", resp.Name, resp.ID, resp.Faculty)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}

func HandleGrade(bot *tgbotapi.BotAPI, update tgbotapi.Update, semesterOrCourseID string) {
	resp, err := services.GetGrades(update.Message.Chat.ID, semesterOrCourseID)
	var response string
	if err != nil {
		response = "Không thể lấy dữ liệu điểm." + err.Error()
	} else {
		response = fmt.Sprintf("Kết quả điểm cho %s:\n________\n%s:\n", semesterOrCourseID, resp.Name)

		if resp.Score.BT != nil {
			response += fmt.Sprintf("  - BT: %.1f\n", *resp.Score.BT)
		} else {
			response += "  - BT: null\n"
		}

		if resp.Score.TN != nil {
			response += fmt.Sprintf("  - TN: %.1f\n", *resp.Score.TN)
		} else {
			response += "  - TN: null\n"
		}

		if resp.Score.BTL != nil {
			response += fmt.Sprintf("  - BTL: %.1f\n", *resp.Score.BTL)
		} else {
			response += "  - BTL: null\n"
		}

		if resp.Score.GK != nil {
			response += fmt.Sprintf("  - Giữa kỳ: %.1f\n", *resp.Score.GK)
		} else {
			response += "  - GK: null\n"
		}

		if resp.Score.CK != nil {
			response += fmt.Sprintf("  - CK: %.1f\n", *resp.Score.CK)
		} else {
			response += "  - CK: null\n"
		}
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}

func HandleAllGrade(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	resp, err := services.GetAllGrades()
	var response string
	if err != nil {
		response = "Không thể lấy dữ liệu điểm." + err.Error()
	} else {
		response = "Kết quả điểm:\n________\n"
		for _, grade := range resp.AllGrades {
			response += fmt.Sprintf("Sinh viên: %s\nMSSV: %s\n", grade.Ms, grade.Name)

			if grade.Score.BT != nil {
				response += fmt.Sprintf("  - BT: %.1f\n", *grade.Score.BT)
			} else {
				response += "  - BT: null\n"
			}

			if grade.Score.TN != nil {
				response += fmt.Sprintf("  - TN: %.1f\n", *grade.Score.TN)
			} else {
				response += "  - TN: null\n"
			}

			if grade.Score.BTL != nil {
				response += fmt.Sprintf("  - BTL: %.1f\n", *grade.Score.BTL)
			} else {
				response += "  - BTL: null\n"
			}

			if grade.Score.GK != nil {
				response += fmt.Sprintf("  - Giữa kỳ: %.1f\n", *grade.Score.GK)
			} else {
				response += "  - GK: null\n"
			}

			if grade.Score.CK != nil {
				response += fmt.Sprintf("  - CK: %.1f\n", *grade.Score.CK)
			} else {
				response += "  - CK: null\n"
			}
			response += "________\n"
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
	// history, err := services.GetHistory(update.Message.Chat.ID)
	// var response string
	// if err != nil {
	// 	response = "Không có lịch sử tra cứu nào."
	// } else {
	// 	response = "Lịch sử tra cứu:\n"
	// 	for _, entry := range history {
	// 		response += fmt.Sprintf("%s: %.1f\n", entry.CourseName, entry.Score)
	// 	}
	// }
	// msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	// bot.Send(msg)
	fmt.Print("History")
}
