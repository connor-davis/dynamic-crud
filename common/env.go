package common

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func EnvString(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		return defaultValue
	}

	return value
}
