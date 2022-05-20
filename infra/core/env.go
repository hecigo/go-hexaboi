package core

import (
	"os"
	"strconv"
	"time"
)

func Getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Get env as a bool value
func GetBoolEnv(key string, fallback bool) bool {
	strVal := Getenv(key, "")
	if strVal == "" {
		return fallback
	}
	val, err := strconv.ParseBool(strVal)
	if err != nil {
		return fallback
	}
	return val
}

// Get env as an integer number, if is impossible then return `fallback`
func GetIntEnv(key string, fallback int) int {
	strVal := Getenv(key, "")
	val, err := strconv.Atoi(strVal)
	if err != nil {
		return fallback
	}
	return val
}

// Ge env as a time duration.
func GetDurationEnv(key string, fallback time.Duration) time.Duration {
	strVal := Getenv(key, "")
	val, err := time.ParseDuration(strVal)
	if err != nil {
		return fallback
	}
	return val
}
