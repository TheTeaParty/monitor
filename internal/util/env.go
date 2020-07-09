package util

import "os"

// GetEnv shortcut to get env value with default value
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		value = defaultValue
	}

	return value
}
