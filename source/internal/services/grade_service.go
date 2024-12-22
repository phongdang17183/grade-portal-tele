package services

import (
	config "Grade_Portal_TelegramBot/config"
	"Grade_Portal_TelegramBot/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func GetStudentInfo(chatID int64, cfg *config.Config) (*struct {
	Email   string `json:"Email"`
	Name    string `json:"Name"`
	Ms      string `json:"Ms"`
	Faculty string `json:"Faculty"`
}, error) {
	endpoint := "/info"
	url := cfg.APIURL + endpoint

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	token, err := GetTokenByChatID(chatID, config.MongoClient)
	if err != nil {
		return nil, fmt.Errorf("error getting token: %w", err)
	} else {
		req.Header.Add("Authorization", "Bearer "+token.Token)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var info models.InfoSV
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	resInfo := struct {
		Email   string `json:"Email"`
		Name    string `json:"Name"`
		Ms      string `json:"Ms"`
		Faculty string `json:"Faculty"`
	}{
		Email:   info.InfoSv.Email,
		Name:    info.InfoSv.Name,
		Ms:      info.InfoSv.Ms,
		Faculty: info.InfoSv.Faculty,
	}
	return &resInfo, nil
}

func GetGrades(chatID int64, semesterOrCourseID string, cfg *config.Config) (*models.Grade, error) {

	if semesterOrCourseID == "" {
		return nil, fmt.Errorf("bạn chưa nhập mã môn. Định dạng đúng: /grade mã môn-học kỳ, ví dụ /grade CO3103-HK233")
	}

	if !strings.Contains(semesterOrCourseID, "-") {
		return nil, fmt.Errorf("mã môn sai (thiếu học kỳ). Định dạng đúng: /grade mã môn-học kỳ, ví dụ /grade CO3103-HK233")
	}

	endpoint := `/resultScore/getmark/` + semesterOrCourseID

	url := cfg.APIURL + endpoint

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	token, err := GetTokenByChatID(chatID, config.MongoClient)
	// Thêm Authorization header với biến token
	if err != nil {
		return nil, fmt.Errorf("error getting token: %w", err)
	} else {
		req.Header.Set("Authorization", "Bearer "+token.Token)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, fmt.Errorf("mã môn %s không tồn tại. Vui lòng kiểm tra lại", semesterOrCourseID)
	case http.StatusBadRequest:
		return nil, fmt.Errorf("mã môn %s không hợp lệ hoặc đã hết hạn. Vui lòng kiểm tra lại", semesterOrCourseID)
	case http.StatusOK:
		// Tiếp tục xử lý nếu trạng thái là 200 OK
		break
	default:
		return nil, fmt.Errorf("lỗi không mong muốn từ API: mã trạng thái %d", resp.StatusCode)
	}

	var grades models.Grade
	if err := json.NewDecoder(resp.Body).Decode(&grades); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	course := models.Course{
		CourseID:   semesterOrCourseID,
		CourseName: grades.Name,
		Score: models.Score{
			BT:  grades.Score.BT,
			TN:  grades.Score.TN,
			BTL: grades.Score.BTL,
			GK:  grades.Score.GK,
			CK:  grades.Score.CK,
		},
	}

	res := AddCourseToHistory(chatID, semesterOrCourseID, course)
	if res != nil {
		log.Fatalf("Lỗi khi thêm khóa học: %v", res)
	}
	return &grades, nil
}

func GetAllGrades(chatID int64, cfg *config.Config) (*models.AllGrades, error) {

	endpoint := "/resultScore/getmark"
	url := cfg.APIURL + endpoint

	// Tạo HTTP request với phương thức GET
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Thêm header Authorization với token
	token, err := GetTokenByChatID(chatID, config.MongoClient)
	if err != nil {
		return nil, fmt.Errorf("error getting token: %v", err)
	} else {
		req.Header.Set("Authorization", "Bearer "+token.Token)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Gửi request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Kiểm tra status code của response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var allGrades models.AllGrades
	if err := json.NewDecoder(resp.Body).Decode(&allGrades); err != nil {
		log.Fatalf("Something wrong in get all grade %v", err)
	}

	for _, grade := range allGrades.AllGrades {
		// Lấy điểm của khóa học
		score := models.Score{
			BT:  grade.Score.BT,
			TN:  grade.Score.TN,
			BTL: grade.Score.BTL,
			GK:  grade.Score.GK,
			CK:  grade.Score.CK,
		}

		// Lưu vào history cho từng khóa học
		err := AddAllCourseToHistory(chatID, grade, score)
		if err != nil {
			log.Fatalf("Lỗi khi thêm khóa học: %v", err)
		}
	}
	return &allGrades, nil
}
