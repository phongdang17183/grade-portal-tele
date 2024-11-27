package models

type DBToken struct {
	Mssv   string `json:"mssv"`
	IDTele int64  `json:"id_tele"`
	Token  string `json:"token"`
}
type DBHistory struct {
	ChatID     int64    `json:"chat_id"`
	ListCourse []string `json:"list_course"`
}
