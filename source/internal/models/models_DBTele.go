package models

type DBToken struct {
	Mssv   string `json:"mssv"`
	ChatID int64  `json:"chat_id"`
	Token  string `json:"token"`
}
type DBHistory struct {
	ChatID     int64    `json:"chat_id"`
	ListCourse []string `json:"list_course"`
}
