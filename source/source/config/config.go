// config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APIURL string
	DBUser string
}

// config/config.go
func LoadConfig() *Config {
	// Nếu bạn muốn tải dev.env, bạn cần chỉ định rõ ràng tên tệp
	err := godotenv.Load("dev.env") // Thêm tên tệp ở đây
	if err != nil {
	    log.Fatal("Error loading .env file")
	}
 
	return &Config{
	    APIURL: os.Getenv("API_URL"),
	    DBUser: os.Getenv("DB_USER"),
	}
 }
 