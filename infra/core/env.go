package core

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var API_CLIENT_SECRETS map[string]string

func Getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(strings.ToUpper(key)); ok {
		return value
	}
	return fallback
}

// Get env as a bool value
func GetBoolEnv(key string, fallback bool) bool {
	strVal := Getenv(strings.ToUpper(key), "")
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
	strVal := Getenv(strings.ToUpper(key), "")
	val, err := strconv.Atoi(strVal)
	if err != nil {
		return fallback
	}
	return val
}

// Ge env as a time duration.
func GetDurationEnv(key string, fallback time.Duration) time.Duration {
	strVal := Getenv(strings.ToUpper(key), "")
	val, err := time.ParseDuration(strVal)
	if err != nil {
		return fallback
	}
	return val
}

// Get application name
func AppName() string {
	return Getenv("APP_NAME", "gohexaboi")
}

// Get application version
func AppVersion() string {
	return Getenv("APP_VERSION", "v0.0.0")
}

// Load API keys from env
func LoadClientSecretKeys() {
	apiClients := strings.Split(Getenv("API_CLIENT_IDS", ""), ",")
	API_CLIENT_SECRETS = make(map[string]string)
	for _, client := range apiClients {
		clientName := strings.ToUpper(client)
		API_CLIENT_SECRETS[clientName] = Getenv(fmt.Sprintf("API_CLIENT_SECRET_%s", strings.ToUpper(clientName)), "")
	}
}
