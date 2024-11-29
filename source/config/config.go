package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APIURL    string
	BOT_TOKEN string
	DBURL     string
}

func LoadConfig() *Config {
	log.Println("Loading .env file...")
	err1 := godotenv.Load(".env")
	if err1 != nil {
		// log.Fatalf("Error loading .env file: %v", err)
		println("Error loading .env file aaaaa: %v", err1)

	}

	config := &Config{
		APIURL:    os.Getenv("API_URL"),
		BOT_TOKEN: os.Getenv("BOT_TOKEN"),
		DBURL:     os.Getenv("DBURL"),
	}

	if config.BOT_TOKEN == "" {
		log.Fatal("BOT_TOKEN is not set in the environment")
	}

	log.Println(".env file loaded successfully.")
	return config
}
