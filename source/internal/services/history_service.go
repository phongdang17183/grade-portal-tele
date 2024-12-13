package services

import (
	"Grade_Portal_TelegramBot/config"
	"Grade_Portal_TelegramBot/internal/models"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func AddCourseToHistory(chatID int64, courseID string, course models.Course) error {
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
					CourseName: course.CourseName,
					CourseID:   courseID,
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

	if err != nil {
		return fmt.Errorf("error finding history: %w", err)
	}

	// Kiểm tra nếu list_course trống và cần khởi tạo
	if len(history.ListCourse) == 1 && history.ListCourse[0].CourseName == "" {
		history.ListCourse = []models.Course{}
	}

	// Kiểm tra nếu course đã tồn tại, cập nhật thông tin
	for _, c := range history.ListCourse {
		if c.CourseID == course.CourseID {
			update := map[string]interface{}{
				"$set": map[string]interface{}{
					"list_course.$[elem].Score": course.Score,
				},
			}
			arrayFilters := []interface{}{
				map[string]interface{}{"elem.CourseName": course.CourseName},
			}
			updateOptions := options.Update().SetArrayFilters(options.ArrayFilters{Filters: arrayFilters})

			_, updateErr := collection.UpdateOne(ctx, filter, update, updateOptions)
			if updateErr != nil {
				return fmt.Errorf("error updating course score: %w", updateErr)
			}
			return nil
		}
	}

	// Nếu course chưa tồn tại, thêm mới vào list_course
	update := map[string]interface{}{
		"$push": map[string]interface{}{
			"list_course": models.Course{
				CourseName: course.CourseName,
				CourseID:   courseID,
				Score:      course.Score,
			},
		},
	}

	_, updateErr := collection.UpdateOne(ctx, filter, update)
	if updateErr != nil {
		return fmt.Errorf("error adding new course: %w", updateErr)
	}

	return nil
}

func AddAllCourseToHistory(chatID int64, grade models.Grades, score models.Score) error {
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
					CourseName: grade.Name,
					CourseID:   grade.Ms,
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
	// Kiểm tra nếu list_course trống và cần khởi tạo
	if len(history.ListCourse) == 1 && history.ListCourse[0].CourseName == "" {
		history.ListCourse = []models.Course{}
	}

	// Kiểm tra nếu khóa học đã tồn tại
	for _, c := range history.ListCourse {
		if c.CourseID == grade.Ms {
			// Cập nhật thông tin điểm nếu đã tồn tại
			update := map[string]interface{}{
				"$set": map[string]interface{}{
					"list_course.$[elem].Score": score,
				},
			}
			arrayFilters := options.ArrayFilters{
				Filters: []interface{}{
					map[string]interface{}{"elem.CourseID": grade.Ms},
				},
			}
			updateOptions := options.Update().SetArrayFilters(arrayFilters)

			_, updateErr := collection.UpdateOne(ctx, filter, update, updateOptions)
			if updateErr != nil {
				return fmt.Errorf("error updating course score: %w", updateErr)
			}
			return nil
		}
	}

	// Nếu chưa có, cập nhật vào danh sách khóa học
	update := map[string]interface{}{
		"$push": map[string]interface{}{
			"list_course": models.Course{
				CourseName: grade.Name,
				CourseID:   grade.Ms,
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
