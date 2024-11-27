package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Info struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Ms        string    `json:"ms"`
	Faculty   string    `json:"faculty"`
	Role      string    `json:"role"`
	CreatedBy string    `json:"created_by"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Score struct {
	BT  *float64 `json:"BT"`
	TN  *float64 `json:"TN"`
	BTL *float64 `json:"BTL"`
	GK  *float64 `json:"GK"`
	CK  *float64 `json:"CK"`
}

type Grade struct {
	Name  string `json:"name"`
	Score Score  `json:"score"`
}
type Grades struct {
	Ms    string `json:"ms"`
	Name  string `json:"name"`
	Score Score  `json:"score"`
}
type AllGrades struct {
	AllGrades []Grades `json:"all_grades"`
}

// RegisterStudent đăng ký tài khoản sinh viên cho người dùng Telegram
func RegisterStudent(chatID int64, mssv string, pw string, otp string) bool {

	base_url := "https://api.example.com"
	endpoint := "/register"

	url := base_url + endpoint

	data := struct {
		MSSV string `json:"mssv"`
		PW   string `json:"password"`
		OTP  string `json:"otp"`
	}{
		MSSV: mssv,
		PW:   pw,
		OTP:  otp,
	}

	// Chuyển dữ liệu sang JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
		return false
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")

	// Gửi request bằng HTTP client
	client := &http.Client{
		Timeout: 10 * time.Second, // Timeout sau 10 giây
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	defer resp.Body.Close()

	return true
}

// Hàm lấy OTP xác thực
func GetOTP(chatID int64, mssv string) bool {

	base_url := "https://api.example.com"
	endpoint := "/otp"

	url := base_url + endpoint

	data := struct {
		MSSV string `json:"mssv"`
	}{
		MSSV: mssv,
	}

	// Chuyển dữ liệu sang JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
		return false
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")

	// Gửi request bằng HTTP client
	client := &http.Client{
		Timeout: 10 * time.Second, // Timeout sau 10 giây
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	return true
}

// Login đăng nhập để bắt đầu sử dụng hệ thống
func Login(mssv string, pw string) bool {

	base_url := "https://api.example.com"
	endpoint := "/login"

	url := base_url + endpoint

	data := struct {
		Ms string `json:"ms"`
		PW string `json:"password"`
	}{
		Ms: mssv,
		PW: pw,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
		return false
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Gửi request bằng HTTP client
	client := &http.Client{
		Timeout: 10 * time.Second, // Timeout sau 10 giây
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Đọc dữ liệu từ response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// In ra kết quả
	fmt.Println(string(body))

	return true
}

// GetStudentInfo lấy thông tin sinh viên dựa trên chatID
func GetStudentInfo(chatID int64) (Student, error) {
	student, exists := users[chatID]
	if !exists {
		return Student{}, errors.New("student not found")
	}

	return student, nil
}

// GetGrades lấy danh sách điểm của sinh viên dựa trên chatID và học kỳ hoặc mã môn học
func GetGrades(chatID int64, semesterOrCourseID string) ([]Grade, error) {
	studentGrades, exists := grades[chatID]
	if !exists {
		return nil, errors.New("grades not found")
	}

	base_url := "https://api.example.com"
	endpoint := "/register"

	url := base_url + endpoint

	// Tạo HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// token := GetToken(chatID) // can hien thuc TODO
	token := token[chatID]

	// Thêm Authorization header với biến token
	req.Header.Set("Authorization", "Bearer "+token)

	// Gửi request
	client := &http.Client{
		Timeout: 10 * time.Second, // Timeout sau 10 giây
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Kiểm tra mã trạng thái
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Request failed with status: %d", resp.StatusCode)
	}

	// Đọc dữ liệu từ response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// In ra kết quả
	fmt.Println("Response:", string(body))

	return studentGrades, nil
}

func GetAllGrades(apiURL string, token string) ([]byte, error) {
	// Tạo một HTTP client
	client := &http.Client{}

	// Tạo HTTP request với phương thức GET
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Thêm header Authorization với token
	req.Header.Add("Authorization", "Bearer "+token)

	// Gửi request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Kiểm tra status code của response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return body, nil
}

// ClearHistory xóa lịch sử tra cứu của sinh viên dựa trên chatID
func ClearHistory(chatID int64) {
	delete(history, chatID)
}

// GetHistory lấy lịch sử tra cứu của sinh viên dựa trên chatID
func GetHistory(chatID int64) ([]Grade, error) {
	studentHistory, exists := history[chatID]
	if !exists {
		return nil, errors.New("history not found")
	}
	return studentHistory, nil
}
