package services

import (
	config "Grade_Portal_TelegramBot/config"
	"Grade_Portal_TelegramBot/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var cfg = config.LoadConfig()

func GetStudentInfo(chatID int64) (*struct {
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

func GetGrades(chatID int64, semesterOrCourseID string) (*models.Grade, error) {

	endpoint := `/resultScore/getmark/` + semesterOrCourseID

	url := cfg.APIURL + endpoint

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	token, err := GetTokenByChatID(chatID, config.MongoClient)
	fmt.Println(token.ChatID)
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var grades models.Grade
	if err := json.NewDecoder(resp.Body).Decode(&grades); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	res := AddCourseToHistory(chatID, semesterOrCourseID)
	if res != nil {
		log.Fatalf("Lỗi khi thêm khóa học: %v", err)
	}
	return &grades, nil
}

func GetAllGrades(chatID int64) (*models.AllGrades, error) {

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
	fmt.Println(allGrades)
	for _, a := range allGrades.AllGrades {
		res := AddCourseToHistory(chatID, a.Ms)
		if res != nil {
			log.Fatalf("Lỗi khi thêm khóa học: %v", err)
		}
	}
	return &allGrades, nil
}
