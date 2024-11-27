package models

import "time"

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
