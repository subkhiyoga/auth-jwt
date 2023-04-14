package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func DotEnv(key string) string {
	// load .env file
	if err := godotenv.Load("database.env"); err != nil {
		log.Fatalln("error saat load .env file")
	}

	return os.Getenv(key)
}
