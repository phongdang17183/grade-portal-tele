package handlers

import (
	"Grade_Portal_TelegramBot/config"
	"Grade_Portal_TelegramBot/internal/services"
	"encoding/json"
	"fmt"
	"strings"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleInfo(bot *tgbotapi.BotAPI, update tgbotapi.Update, cfg *config.Config) {
	resp, err := services.GetStudentInfo(update.Message.Chat.ID, cfg)
	
	if err != nil {
		errorDetails := "Lỗi: " + err.Error()
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

func HandleGrade(bot *tgbotapi.BotAPI, update tgbotapi.Update, semesterOrCourseID string, cfg *config.Config) {
	resp, err := services.GetGrades(update.Message.Chat.ID, semesterOrCourseID, cfg)
	var response string

	if err != nil {
		response = "Không thể lấy dữ liệu điểm: " + err.Error()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		bot.Send(msg)
		return
	}

	result := map[string]interface{}{
		"course_id":   semesterOrCourseID,
		"course_name": resp.Name,
		"scores": map[string]interface{}{
			"BT":  resp.Score.BT,
			"TN":  resp.Score.TN,
			"BTL": resp.Score.BTL,
			"GK":  resp.Score.GK,
			"CK":  resp.Score.CK,
		},
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		response = "Lỗi khi tạo điểm vui lòng thử lại sau: " + err.Error()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		bot.Send(msg)
		return
	}
	msgText := fmt.Sprintf("```json\n%s\n```", string(jsonData))
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	msg.ParseMode = "MarkdownV2"
	bot.Send(msg)
}

func HandleAllGrade(bot *tgbotapi.BotAPI, update tgbotapi.Update, cfg *config.Config) {
	resp, err := services.GetAllGrades(update.Message.Chat.ID, cfg)
	var response interface{}
	if err != nil {
		response = map[string]string{
			"error": "Không thể lấy dữ liệu điểm: " + err.Error(),
		}
	} else {
		var grades []map[string]interface{}
		for _, grade := range resp.AllGrades {
			gradeData := map[string]interface{}{
				"course_id":   grade.Ms,
				"course_name": grade.Name,
				"scores": map[string]interface{}{
					"BT":  grade.Score.BT,
					"TN":  grade.Score.TN,
					"BTL": grade.Score.BTL,
					"GK":  grade.Score.GK,
					"CK":  grade.Score.CK,
				},
			}
			grades = append(grades, gradeData)
		}
		response = map[string]interface{}{
			"grades": grades,
		}
	}
	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		// Xử lý lỗi khi mã hóa JSON
		fmt.Println("Lỗi khi mã hóa JSON:", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Không thể xử lý dữ liệu JSON.")
		bot.Send(msg)
		return
	}
	fmt.Println("Dữ liệu JSON trả về:", string(responseJSON))
	// Gửi phản hồi dưới dạng JSON thô
	msgText := fmt.Sprintf("```json\n%s\n```", string(responseJSON))
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(msgText))
	msg.ParseMode = "MarkdownV2" // Nếu bạn muốn hiển thị trong markdown
	bot.Send(msg)
}

