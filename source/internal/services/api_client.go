package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// User struct để lưu dữ liệu người dùng từ API
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// FetchUsers function để gọi API từ JSONPlaceholder
func FetchUsers() ([]User, error) {
	// Gửi yêu cầu GET tới API
	resp, err := http.Get("https://jsonplaceholder.typicode.com/users")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Đọc dữ liệu từ phản hồi
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Phân tích JSON
	var users []User
	err = json.Unmarshal(body, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}
