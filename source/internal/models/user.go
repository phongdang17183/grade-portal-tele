package models

type User struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	MSSV    string `json:"ms"`
	Faculty string `json:"faculty"`
	Role    string `json:"role"`
}
