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

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file : %v", err)
	}

	config := &Config{
		APIURL:    os.Getenv("API_URL"),
		BOT_TOKEN: os.Getenv("BOT_TOKEN"),
		DBURL:     os.Getenv("DBURL"),
	}

	if config.APIURL == "" {
		log.Fatal("API_URL is not set in the environment")
	}
	if config.BOT_TOKEN == "" {
		log.Fatal("BOT_TOKEN is not set in the environment")
	}
	if config.DBURL == "" {
		log.Fatal("DBURL is not set in the environment")
	}
	return config
}
