package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jihedmastouri/game-integration-api-demo/repository"
	"github.com/jihedmastouri/game-integration-api-demo/service"
	"github.com/jihedmastouri/game-integration-api-demo/transport"
	"github.com/joho/godotenv"
)

func main() {
	loadDotenv()

	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		getDefaultEnv("PG_USER", "postgres"),
		getDefaultEnv("PG_PASS", "postgres"),
		getDefaultEnv("PG_URL", "localhost"),
		getDefaultEnv("PG_PORT", "5432"),
		getDefaultEnv("PG_DB", "postgres"),
	)

	_ = repository.Connect(databaseUrl)
	srv := service.Service{}
	transport.Web(srv)
}

func getDefaultEnv(name, defaultValue string) string {
	if envValue := os.Getenv(name); envValue != "" {
		return envValue
	}
	return defaultValue
}

func loadDotenv() {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		log.Println("Warning: .env file not found")
	} else {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading the .env")
		}
	}
}
