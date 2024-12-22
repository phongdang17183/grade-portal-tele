// cmd/api/other.go
package services

import (
	"Grade_Portal_TelegramBot/config"
	"Grade_Portal_TelegramBot/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func RegisterStudent(ms string, pw string, otp string, cfg *config.Config) (*models.MsgResp, error) {

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

func GetOTP(mssv string, cfg *config.Config) (*models.MsgResp, error) {

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

	var msgResp models.MsgResp
	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode != http.StatusOK {
		json.NewDecoder(resp.Body).Decode(&msgResp)
		// return nil, fmt.Errorf("unexpected status code: %d ", &msgResp)
		return nil, fmt.Errorf(msgResp.Msg)
	}

	if err := json.NewDecoder(resp.Body).Decode(&msgResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &msgResp, nil
}

func Login(chatID int64, mssv string, pw string, cfg *config.Config) (*models.ResLogin, error) {
	endpoint := "/loginTele"
	url := cfg.APIURL + endpoint

	// Dữ liệu gửi lên API nhóm BE
	data := struct {
		Ms string `json:"ms"`
		PW string `json:"password"`
	}{
		Ms: mssv,
		PW: pw,
	}

	// Mã hóa JSONN
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error encoding JSON: %w", err)
	}

	// Tạo request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Gửi request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("timeout: api không phản hồi")
		}
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Xử lý mã trạng thái HTTP
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	// Giải mã JSON phản hồi lại
	var resLogin models.ResLogin
	if err := json.NewDecoder(resp.Body).Decode(&resLogin); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	// Lưu token vào MONGODB
	// token := models.DBToken{
	// 	Mssv:   mssv,
	// 	ChatID: chatID,
	// 	Token:  resLogin.Token,
	// }
	err = saveTokenToDB(chatID, mssv, resLogin.Token)
	if err != nil {
		return nil, fmt.Errorf("error saving token: %w", err)
	}
	return &resLogin, nil
}

// Hàm lưu token vào MongoDB (được tách ra để tái sử dụng hoặc kiểm thử)
func saveTokenToDB(chatID int64, mssv, token string) error {
	collection := config.MongoClient.Database("Do_an").Collection("TOKEN")

	filter := map[string]interface{}{"chat_id": chatID} // Kiểm tra dựa trên MSSV
	update := map[string]interface{}{
		"$set": map[string]interface{}{
			"mssv": mssv,
			"token":   token,
		},
	}
	opts := options.Update().SetUpsert(true)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("MongoDB UpdateOne failed: %w", err)
	}

	// Log kết quả lưu trữ
	if result.MatchedCount > 0 {
		log.Println("Token đã được cập nhật.")
	} else if result.UpsertedCount > 0 {
		log.Printf("Thêm mới token thành công, ID: %v\n", result.UpsertedID)
	}
	return nil
}

func GetTokenByChatID(chatID int64, client *mongo.Client) (*models.DBToken, error) {

	var token models.DBToken

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database("Do_an").Collection("TOKEN")
	// Bộ lọc tìm kiếm
	filter := map[string]interface{}{"chat_id": chatID}

	err := collection.FindOne(ctx, filter).Decode(&token)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no token found for chatID %d", chatID)
		}
		return nil, fmt.Errorf("error finding token: %w", err)
	}
	return &token, nil
}
