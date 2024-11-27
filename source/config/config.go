package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	APIURL    string
	DBUser    string
	BOT_TOKEN string
	DBURL     string
}

func LoadConfig() *Config {
	log.Println("Loading .env file from Grade_Portal_TelegramBot/config/.env...")
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	config := &Config{
		APIURL:    os.Getenv("API_URL"),
		DBUser:    os.Getenv("DB_USER"),
		BOT_TOKEN: os.Getenv("BOT_TOKEN"),
		DBURL:     os.Getenv("DBURL"),
	}

	if config.BOT_TOKEN == "" {
		log.Fatal("BOT_TOKEN is not set in the environment")
	}

	log.Println(".env file loaded successfully.")
	return config
}
