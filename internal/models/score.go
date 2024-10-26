package models

type Score struct {
	MSSV      string   `json:"mssv"`
	Grades    []float64 `json:"grades"`
	GPA       float64   `json:"gpa"`
}
