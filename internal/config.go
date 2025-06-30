package internal

import (
	"fmt"
	"os"
)

func init() {
	Config.JWT_SECRET = getDefaultEnv("JWT_SECRET", "naUsB1EQS9U")

	Config.WALLET_API_KEY = getDefaultEnv("WALLET_API_KEY", "naUsB1EQS9U")
	Config.WALLET_API_URL = getDefaultEnv("WALLET_API_URL", "http://locahost:8000")

	Config.DATABASE_URL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		getDefaultEnv("PG_USER", "postgres"),
		getDefaultEnv("PG_PASS", "postgres"),
		getDefaultEnv("PG_URL", "localhost"),
		getDefaultEnv("PG_PORT", "5432"),
		getDefaultEnv("PG_DB", "postgres"),
	)

	Config.APP_URL = fmt.Sprintf("%s:%s",
		getDefaultEnv("PG_USER", "postgres"),
		getDefaultEnv("PG_PASS", "postgres"),
	)
}

var Config struct {
	APP_URL        string
	DATABASE_URL   string
	WALLET_API_URL string
	WALLET_API_KEY string
	JWT_SECRET     string
}

func getDefaultEnv(name, defaultValue string) string {
	if envValue := os.Getenv(name); envValue != "" {
		return envValue
	}
	return defaultValue
}
