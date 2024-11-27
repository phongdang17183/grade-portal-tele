package models

type DBToken struct {
	Mssv   string `json:"mssv"`
	IDTele int64  `json:"id_tele"`
	Token  string `json:"token"`
}
type DBHistory struct {
	Mssv      string   `json:"mssv"`
	IDTele    int64    `json:"id_tele"`
	HisCourse []string `json:"his_course"`
}
