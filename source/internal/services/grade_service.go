package services

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	// "io/ioutil"
	"log"
	"net/http"
	"time"
    "encoding/json"
)

// Dữ liệu cứng cho mục đích kiểm tra
var users = map[int64]Student{
    12345: {Name: "Nguyễn Văn A", StudentID: "12345"},
}

var grades = map[int64][]Grade{
    12345: {
        {CourseName: "Toán học", Score: 8.5},
        {CourseName: "Vật lý", Score: 7.0},
    },
}

var history = map[int64][]Grade{
    12345: {
        {CourseName: "Toán học", Score: 8.5},
        {CourseName: "Vật lý", Score: 7.0},
    },
}

var token = map[int64]string{
    123 : "123",
}

type Student struct {
    Name      string
    StudentID string
}

type Grade struct {
    CourseName string
    Score      float64
}

// RegisterStudent đăng ký tài khoản sinh viên cho người dùng Telegram
func RegisterStudent(chatID int64, email string) bool {
    
    base_url := "https://api.example.com"
	endpoint := "/register"
    
	url := base_url + endpoint
    
	data := struct {
		Email string `json:"email"`
	}{
		Email: email,
	}

	// Chuyển dữ liệu sang JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
        return false
	}
        
    req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(jsonData)))
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
func Login(chayID int64, otp string) bool{

	base_url := "https://api.example.com"
	endpoint := "/login"
    
	url := base_url + endpoint
    
	data := struct {
		Otp string `json:"otp"`
	}{
		Otp: otp,
	}

	// Chuyển dữ liệu sang JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
        return false
	}
        
    req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(jsonData)))
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
	req.Header.Set("Authorization", "Bearer "+ token)

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
