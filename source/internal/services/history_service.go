package services

import (
	"Grade_Portal_TelegramBot/config"
	"Grade_Portal_TelegramBot/internal/models"
	"context"
	"fmt"
	"log"
	"time"

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

func GetHistory(chatID int64) (*[]models.Course, error) {
	// Lấy lịch sử từ chatID
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.MongoClient.Database("Do_an").Collection("HISTORY")

	filter := map[string]interface{}{"chat_id": chatID}

	var result models.DBHistory

	err := collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no history found for chatID %d", chatID)
		}
		return nil, fmt.Errorf("error finding history: %w", err)
	}

	return &result.ListCourse, nil
}


func AddCourseToHistory(chatID int64, courseName string, course models.Course) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.MongoClient.Database("Do_an").Collection("HISTORY")

	filter := map[string]interface{}{"chat_id": chatID}
	var history models.DBHistory
	err := collection.FindOne(ctx, filter).Decode(&history)
	// Nếu không tìm thấy tài liệu, tạo tài liệu mới
	if err == mongo.ErrNoDocuments {
		newHistory := map[string]interface{}{
			"chat_id": chatID,
			"list_course": []models.Course{
				{
					CourseName: courseName,
					Score:      course.Score,
				},
			},
		}

		_, insertErr := collection.InsertOne(ctx, newHistory)
		if insertErr != nil {
			return fmt.Errorf("error inserting new history: %w", insertErr)
		}
		return nil
	}

	if len(history.ListCourse) == 1 && history.ListCourse[0].CourseName == "" {
		history.ListCourse = []models.Course{}
	}

	for _, c := range history.ListCourse {
		if c.CourseName == course.CourseName {
			return nil
		}
	}

	update := map[string]interface{}{
		"$push": map[string]interface{}{
			"list_course": []models.Course{
				{
					CourseName: courseName,
					Score:      course.Score,
				},
			},
		},
	}

	_, updateErr := collection.UpdateOne(ctx, filter, update)
	if updateErr != nil {
		return fmt.Errorf("error updating history: %w", updateErr)
	}
	return nil
}

func AddAllCourseToHistory(chatID int64, Ms string, score models.Score) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.MongoClient.Database("Do_an").Collection("HISTORY")

	// Tìm kiếm lịch sử của người dùng
	filter := map[string]interface{}{"chat_id": chatID}
	var history models.DBHistory
	err := collection.FindOne(ctx, filter).Decode(&history)

	// Nếu không tìm thấy tài liệu, tạo tài liệu mới
	if err == mongo.ErrNoDocuments {
		newHistory := map[string]interface{}{
			"chat_id": chatID,
			"list_course": []models.Course{
				{
					CourseName: Ms,
					Score:      score,
				},
			},
		}
		_, insertErr := collection.InsertOne(ctx, newHistory)
		if insertErr != nil {
			return fmt.Errorf("error inserting new history: %w", insertErr)
		}
		return nil
	}
	if len(history.ListCourse) == 1 && history.ListCourse[0].CourseName == "" {
		history.ListCourse = []models.Course{}
	}
	// Kiểm tra nếu khóa học đã tồn tại
	for _, c := range history.ListCourse {
		if c.CourseName == Ms {
			return nil
		}
	}

	// Nếu chưa có, cập nhật vào danh sách khóa học
	update := map[string]interface{}{
		"$push": map[string]interface{}{
			"list_course": models.Course{
				CourseName: Ms,
				Score:      score,
			},
		},
	}

	_, updateErr := collection.UpdateOne(ctx, filter, update)
	if updateErr != nil {
		return fmt.Errorf("error updating history: %w", updateErr)
	}
	return nil
}
