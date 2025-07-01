package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	loadDotenv()

	Config.JWT_SECRET = getDefaultEnv("JWT_SECRET", "naUsB1EQS9U")

	Config.WALLET_API_KEY = getDefaultEnv("WALLET_API_KEY", "naUsB1EQS9U")
	Config.WALLET_API_URL = getDefaultEnv("WALLET_API_URL", "http://locahost:8000")

	mode := getDefaultEnv("MODE", "dev")
	if mode == "production" {
		Config.MODE = ModeProduction
	} else {
		Config.MODE = ModeDevelopment
	}

	Config.DATABASE_URL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s%s",
		getDefaultEnv("PG_USER", "postgres"),
		getDefaultEnv("PG_PASS", "postgres"),
		getDefaultEnv("PG_URL", "localhost"),
		getDefaultEnv("PG_PORT", "5432"),
		getDefaultEnv("PG_DB", "postgres"),
		getDefaultEnv("PG_MODE", ""),
	)

	Config.APP_URL = fmt.Sprintf("%s:%s",
		getDefaultEnv("APP_HOST", "0.0.0.0"),
		getDefaultEnv("APP_PORT", "3000"),
	)

	maxIdle, err := strconv.Atoi(getDefaultEnv("DB_MAX_IDLE", "3"))
	if err != nil {
		maxIdle = 3
	}
	Config.DB_MAX_IDLE = maxIdle

	maxOpen, err := strconv.Atoi(getDefaultEnv("DB_MAX_OPEN", "0"))
	if err != nil || maxOpen == 0 {
		maxOpen = 3
	}
	Config.DB_MAX_OPEN = maxOpen
}

type ModeType string

const (
	ModeProduction  ModeType = "prod"
	ModeDevelopment ModeType = "dev"
)

var Config struct {
	APP_URL        string
	DATABASE_URL   string
	WALLET_API_URL string
	WALLET_API_KEY string
	JWT_SECRET     string
	MODE           ModeType
	DB_MAX_OPEN    int
	DB_MAX_IDLE    int
}

func getDefaultEnv(name, defaultValue string) string {
	if envValue := os.Getenv(name); envValue != "" {
		return envValue
	}
	return defaultValue
}

func loadDotenv() {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory: ", err)
		return
	}

	filename := filepath.Join(path, ".env")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("Warning: .env file not found: ", filename)
		return
	}

	err = godotenv.Load(filename)
	if err != nil {
		fmt.Println("Error loading the .env file: ", err)
	} else {
		fmt.Println(".env file loaded successfully: ", filename)
	}
}
