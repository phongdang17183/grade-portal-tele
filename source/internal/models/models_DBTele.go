package models

type DBToken struct {
	Mssv   string `bson:"mssv"`
	ChatID int64  `bson:"chat_id"`
	Token  string `bson:"token"`
}
type DBHistory struct {
	ChatID     int64    `bson:"chat_id"`
	ListCourse []Course `bson:"list_course"`
}
type Course struct {
	CourseName string `bson:"course_name"`
	Score      Score  `bson:"score"`
}
