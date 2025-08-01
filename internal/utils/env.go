package utils

import (
	"log"

	"github.com/joho/godotenv"
)

func InitEnvVars() {
	err := godotenv.Load()

	if err != nil {
		log.Printf("Ошибка при загрузке переменных из .env: %v", err)
	}
}
