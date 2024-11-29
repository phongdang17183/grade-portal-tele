package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APIURL    string
	DBUser    string
	BOT_TOKEN string
}

func LoadConfig() *Config {
	log.Println("Loading .env file...")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	config := &Config{
		APIURL:    os.Getenv("API_URL"),
		DBUser:    os.Getenv("DB_USER"),
		BOT_TOKEN: os.Getenv("BOT_TOKEN"),
	}

	if config.BOT_TOKEN == "" {
		log.Fatal("BOT_TOKEN is not set in the environment")
	}

	log.Println(".env file loaded successfully.")
	return config
}
