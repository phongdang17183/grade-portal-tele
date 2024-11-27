package services

import (
	config "Grade_Portal_TelegramBot/config"
	"Grade_Portal_TelegramBot/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	// "errors"
	"fmt"

	"net/http"
	"time"
)

var cfg = config.LoadConfig()

func GetTokenByChatID(chatID int64, client *mongo.Client) (*models.DBToken, error) {
	// Thiết lập bối cảnh
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Kết nối tới collection
	collection := client.Database("Do_an").Collection("TOKEN")

	// Bộ lọc tìm kiếm
	filter := bson.M{"id_tele": chatID}

	// Kết quả để lưu token tìm được
	var token models.DBToken

	// Truy vấn dữ liệu
	err := collection.FindOne(ctx, filter).Decode(&token)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no token found for chatID %d", chatID)
		}
		return nil, fmt.Errorf("error finding token: %w", err)
	}
	return &token, nil
}

func AddCourseToHistory(chatID int64, course string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.MongoClient.Database("Do_an").Collection("HISTORY")

	// Tìm tài liệu theo ChatID
	filter := bson.M{"chat_id": chatID}

	// Kiểm tra xem tài liệu có tồn tại không
	var history models.DBHistory
	err := collection.FindOne(ctx, filter).Decode(&history)
	if err != nil {
		// Nếu không tìm thấy tài liệu, tạo tài liệu mới
		if err == mongo.ErrNoDocuments {
			newHistory := bson.M{
				"chat_id":     chatID,
				"list_course": []string{course},
			}
			_, insertErr := collection.InsertOne(ctx, newHistory)
			if insertErr != nil {
				return fmt.Errorf("error inserting new history: %w", insertErr)
			}
			fmt.Println("Thêm lịch sử mới thành công!")
			return nil
		}
		return fmt.Errorf("error finding history: %w", err)
	}

	// Nếu tài liệu đã tồn tại, kiểm tra xem khóa học đã có chưa
	for _, c := range history.ListCourse {
		if c == course {
			fmt.Println("Khóa học đã tồn tại, không cần thêm lại!")
			return nil
		}
	}

	// Nếu khóa học chưa có, thêm vào danh sách
	update := bson.M{
		"$push": bson.M{
			"list_course": course,
		},
	}

	_, updateErr := collection.UpdateOne(ctx, filter, update)
	if updateErr != nil {
		return fmt.Errorf("error updating history: %w", updateErr)
	}

	fmt.Println("Thêm khóa học thành công!")
	return nil
}

func RegisterStudent(ms string, pw string, otp string) (*models.MsgResp, error) {

	endpoint := "/resetpassword"

	url := cfg.APIURL + endpoint

	data := struct {
		MS  string `json:"ms"`
		PW  string `json:"password"`
		OTP string `json:"otp"`
	}{
		MS:  ms,
		PW:  pw,
		OTP: otp,
	}
	fmt.Println(data)
	// Chuyển dữ liệu sang JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error encoding JSON: %w", err)
	}
	//fmt.Println(string(jsonData))
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
		fmt.Println(resp)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var msgResp models.MsgResp
	if err := json.NewDecoder(resp.Body).Decode(&msgResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &msgResp, nil
}

func GetOTP(mssv string) (*models.MsgResp, error) {

	endpoint := "/otp"

	url := cfg.APIURL + endpoint

	data := struct {
		MSSV string `json:"ms"`
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
		return nil, fmt.Errorf("unexpected status code: %d %w", resp.StatusCode, resp.Body)
	}

	var msgResp models.MsgResp
	if err := json.NewDecoder(resp.Body).Decode(&msgResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &msgResp, nil
}

func Login(chatID int64, mssv string, pw string) (*models.ResLogin, error) {

	endpoint := "/loginTele"
	url := cfg.APIURL + endpoint

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
	var resLogin models.ResLogin
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
	filter := bson.M{"mssv": token.Mssv} // Kiểm tra dựa trên MSSV
	update := bson.M{
		"$set": bson.M{
			"id_tele": token.IDTele,
			"token":   token.Token,
		},
	}
	opts := options.Update().SetUpsert(true)

	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, fmt.Errorf("error saving token to database: %w", err)
	}
	if result.MatchedCount > 0 {
		fmt.Println("Token đã được cập nhật.")
	} else if result.UpsertedCount > 0 {
		fmt.Printf("Thêm mới token thành công, ID: %v\n", result.UpsertedID)
	}
	return &resLogin, nil
}

func GetStudentInfo(chatID int64) (*models.InfoSV, error) {
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
	fmt.Println(info)
	return &info, nil
}

func GetGrades(chatID int64, semesterOrCourseID string) (*models.Grade, error) {

	endpoint := `/resultScore/getmark/` + semesterOrCourseID

	url := cfg.APIURL + endpoint

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// token := GetToken(chatID) // can hien thuc TODO
	token, err := GetTokenByChatID(chatID, config.MongoClient)
	fmt.Println(token)
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
	}
	for _, a := range allGrades.AllGrades {
		res := AddCourseToHistory(chatID, a.Ms)
		if res != nil {
			log.Fatalf("Lỗi khi thêm khóa học: %v", err)
		}
	}
	return &allGrades, nil
}

func ClearHistory(chatID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.MongoClient.Database("Do_an").Collection("history")

	filter := bson.M{"chatID": chatID}

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Error deleting document: %v", err)
		return false
	}
	return true
}

func GetHistory(chatID int64) (*models.DBHistory, error) {
	// Lấy lịch sử từ chatID
	history, err := GetHistoryByChatID(chatID)
	fmt.Println(history)
	if err != nil {
		return nil, err
	}

	// Trả về lịch sử
	return history, nil
}

func GetHistoryByChatID(chatID int64) (*models.DBHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.MongoClient.Database("Do_an").Collection("HISTORY")

	filter := bson.M{"chat_id": chatID}
	var result models.DBHistory

	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no history found for chatID %d", chatID)
		}
		return nil, fmt.Errorf("error finding history: %w", err)
	}
	fmt.Println(result.ListCourse)

	return &result, nil
}
