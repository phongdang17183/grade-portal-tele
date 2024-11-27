package services

import (
	"Grade_Portal_TelegramBot/config"
	"context"
	"fmt"
	"log"
	"time"
	"Grade_Portal_TelegramBot/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

func ClearHistory(chatID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.MongoClient.Database("Do_an").Collection("HISTORY")

	filter := map[string]interface{}{"chat_id": chatID}

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

	filter := map[string]interface{}{
		"chat_id": chatID, 
	}
	var result models.DBHistory

	err := collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no history found for chatID %d", chatID)
		}
		return nil, fmt.Errorf("error finding history: %w", err)
	}

	return &result, nil
}

func AddCourseToHistory(chatID int64, course string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.MongoClient.Database("Do_an").Collection("HISTORY")

	// Tìm tài liệu theo ChatID
	filter := map[string]interface{}{"chat_id": chatID}

	// Kiểm tra xem tài liệu có tồn tại không
	var history models.DBHistory
	err := collection.FindOne(ctx, filter).Decode(&history)
	if err != nil {
		// Nếu không tìm thấy tài liệu, tạo tài liệu mới
		if err == mongo.ErrNoDocuments {
			newHistory := map[string]interface{}{
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
	update := map[string]interface{}{
		"$push": map[string]interface{}{
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