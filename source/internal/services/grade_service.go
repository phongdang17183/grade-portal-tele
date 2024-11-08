package services

import (
    "errors"
    "fmt"
)

// Dữ liệu cứng cho mục đích kiểm tra
var users = map[int64]Student{
    // 12345: {Name: "Nguyễn Văn A", StudentID: "12345"},
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

type Student struct {
    Name      string
    StudentID string
}

type Grade struct {
    CourseName string
    Score      float64
}

// RegisterStudent đăng ký tài khoản sinh viên cho người dùng Telegram
func RegisterStudent(chatID int64, studentID string) bool {
    if _, exists := users[chatID]; exists {
        return false // Tài khoản Telegram đã đăng ký với MSSV khác
    }
    users[chatID] = Student{Name: fmt.Sprintf("Sinh viên %s", studentID), StudentID: studentID}
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
