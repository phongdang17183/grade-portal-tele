package services

import (
	"Grade_Portal_TelegramBot/internal/models"
	"bytes"
	"encoding/json"
	// "errors"
	"fmt"
    config "Grade_Portal_TelegramBot/config"
	//"io"
	"log"
	"net/http"
	"time"
    "context"
)

type ResLogin struct {
	ListCourse []string `json:"listCourse"`
	Token      string   `json:"token"`
}

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
type History struct {
	CourseID []string `json:"courseId"`
}
type Grades struct {
	Ms    string `json:"ms"`
	Name  string `json:"name"`
	Score Score  `json:"score"`
}
type AllGrades struct {
	AllGrades []Grades `json:"all_grades"`
}
type MsgResp struct {
	Msg string `json:"msg"`
}

func RegisterStudent(mssv string, pw string, otp string) (*MsgResp, error) {

	base_url := "https://api.example.com"
	endpoint := "/resetpassword"

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
		return nil, fmt.Errorf("error encoding JSON: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var msgResp MsgResp
	if err := json.NewDecoder(resp.Body).Decode(&msgResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &msgResp, nil
}

func GetOTP(mssv string) (*MsgResp, error) {

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
		return nil, fmt.Errorf("error encoding JSON: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second, // Timeout sau 10 giây
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var msgResp MsgResp
	if err := json.NewDecoder(resp.Body).Decode(&msgResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &msgResp, nil
}

func Login(chatID int64, mssv string, pw string) (*ResLogin, error) {

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
		return nil, fmt.Errorf("error encoding JSON: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
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
	var resLogin ResLogin
	if err := json.NewDecoder(resp.Body).Decode(&resLogin); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	token := models.DBToken{
		Mssv:   mssv,
		IDTele: chatID,
		Token:  resLogin.Token,
	}

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

    collection := config.MongoClient.Database("Do_an").Collection("TOKEN")
    result, err := collection.InsertOne(ctx, token)
    if err != nil {
		log.Fatalf("Lỗi khi lưu dữ liệu: %v", err)
	}
    fmt.Printf("Lưu thành công, ID tài liệu: %v\n", result.InsertedID)

	fmt.Println(token)
	return &resLogin, nil
}

func GetStudentInfo(chatID int64) (*Info, error) {
	// baseURL := "https://api.example.com"
	// endpoint := "/info"
	// url := baseURL + endpoint

	// client := &http.Client{
	// 	Timeout: 10 * time.Second,
	// }

	// req, err := http.NewRequest("GET", url, nil)
	// if err != nil {
	// 	return nil, fmt.Errorf("error creating request: %w", err)
	// }

	// token := token[chatID]
	// req.Header.Set("Authorization", `Bearer `+token)
	// req.Header.Set("Content-Type", "application/json")

	// resp, err := client.Do(req)
	// if err != nil {
	// 	return nil, fmt.Errorf("error sending request: %w", err)
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	// }

	// var info Info
	// if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
	// 	return nil, fmt.Errorf("error decoding response: %w", err)
	// }

	// return &info, nil
    return nil, nil
}

func GetGrades(chatID int64, semesterOrCourseID string) (*Grade, error) {

	// baseURL := "https://api.example.com"
	// endpoint := `/resultScore/getmark/` + semesterOrCourseID

	// url := baseURL + endpoint

	// req, err := http.NewRequest("GET", url, nil)
	// if err != nil {
	// 	return nil, fmt.Errorf("error creating request: %w", err)
	// }

	// // token := GetToken(chatID) // can hien thuc TODO
	// token := token[chatID]

	// // Thêm Authorization header với biến token
	// req.Header.Set("Authorization", "Bearer "+token)
	// req.Header.Set("Content-Type", "application/json")

	// client := &http.Client{
	// 	Timeout: 10 * time.Second,
	// }
	// resp, err := client.Do(req)
	// if err != nil {
	// 	return nil, fmt.Errorf("error sending request: %w", err)
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	// }

	// var grades Grade
	// if err := json.NewDecoder(resp.Body).Decode(&grades); err != nil {
	// 	return nil, fmt.Errorf("error decoding response: %w", err)
	// }

	// return &grades, nil
    return nil, nil
}

func GetAllGrades() (*AllGrades, error) {

	// baseURL := "https://api.example.com"
	// endpoint := "/resultScore/getmark"
	// url := baseURL + endpoint

	// // Tạo HTTP request với phương thức GET
	// req, err := http.NewRequest("GET", url, nil)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create request: %v", err)
	// }

	// // Thêm header Authorization với token
	// token := token[chatID]
	// req.Header.Add("Authorization", "Bearer "+token)
	// req.Header.Set("Content-Type", "application/json")

	// client := &http.Client{
	// 	Timeout: 10 * time.Second,
	// }

	// // Gửi request
	// resp, err := client.Do(req)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to send request: %v", err)
	// }
	// defer resp.Body.Close()

	// // Kiểm tra status code của response
	// if resp.StatusCode != http.StatusOK {
	// 	return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	// }

	// var allGrades AllGrades
	// if err := json.NewDecoder(resp.Body).Decode(&allGrades); err != nil {
	// }

	// return &allGrades, nil
    return nil, nil
}

func ClearHistory(chatID int64) {
	// delete(history, chatID)

}

func GetHistory(chatID int64) (*History, error) {
	// studentHistory, exists := history[chatID]
	// if !exists {
	// 	return nil, errors.New("history not found")
	// }
	// return studentHistory, nil
    return nil, nil
}
