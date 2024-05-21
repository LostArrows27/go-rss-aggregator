package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(name string) string {
	portString := os.Getenv(name)

	if portString == "" {
		log.Fatal("Not found port")
	}

	return portString
}

func LoadEnv(path string) {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatal("Error in loading .env")
	}
}
