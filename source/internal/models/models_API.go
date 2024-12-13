package models

import "time"

type ResLogin struct {
	ListCourse []string `json:"listCourse"`
	Token      string   `json:"token"`
}

type Info struct {
	ID        string    `json:"ID"`
	Email     string    `json:"Email"`
	Name      string    `json:"Name"`
	Ms        string    `json:"Ms"`
	Faculty   string    `json:"Faculty"`
	Role      string    `json:"Role"`
	CreatedBy string    `json:"CreatedBy"`
	ExpiresAt time.Time `json:"ExpiredAt"`
}

type InfoSV struct {
	InfoSv Info `json:"user"`
}

type Score struct {
	BT  []*float64 `json:"BT"`
	TN  []*float64 `json:"TN"`
	BTL []*float64 `json:"BTL"`
	GK  *float64 `json:"GK"`
	CK  *float64 `json:"CK"`
}

type Grade struct {
	Name  string `json:"name"`
	Score Score  `json:"score"`
}

type Grades struct {
	Ms    string `json:"ms"`
	Name  string `json:"name"`
	Score Score  `json:"data"`
}
type AllGrades struct {
	AllGrades []Grades `json:"scores"`
}
type MsgResp struct {
	Msg string `json:"msg"`
}
