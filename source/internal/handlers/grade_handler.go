package handlers

// import (
// 	"gopkg.in/telebot.v3"
// 	"net/http"
// 	"bytes"
// 	"encoding/json"
// )

// type GradeRequest struct {
// 	StudentID string `json:"ms"`
// 	Semester  string `json:"semester"`
// }

// func GetGrades(c *telebot.Context) {
// 	// Gọi API backend để lấy điểm
// 	reqBody := GradeRequest{StudentID: c.Args()[0], Semester: c.Args()[1]}
// 	body, _ := json.Marshal(reqBody)
// 	resp, err := http.Post("http://your-backend.com/api/grades", "application/json", bytes.NewBuffer(body))
// 	if err != nil {
// 		c.Send("Có lỗi xảy ra!")
// 		return
// 	}
// 	defer resp.Body.Close()

// 	// Xử lý dữ liệu trả về
// }
