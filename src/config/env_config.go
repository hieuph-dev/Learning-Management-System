package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// Kiểm tra nếu file .env tồn tại
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load(".env")
		if err != nil {
			log.Println("Error loading .env file:", err)
		} else {
			log.Println("Loaded .env file successfully")
		}
	} else {
		log.Println("No .env file found, using environment variables from system/docker")
	}
}
