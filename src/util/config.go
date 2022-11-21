package util

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// EnvironmentKey key for setting the environment
const EnvironmentKey = "ENV"

// Environment types
const (
	Development = "development"
	Testing     = "testing"
	Staging     = "staging"
	Production  = "production"
)

func GetEnvString(key string, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		if IsDevelopmentEnvironment() {
			return fallback
		}
		panic(fmt.Errorf("undefined environment variable %s", key))
	}
	return value
}

func GetEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if len(value) == 0 {
		if IsDevelopmentEnvironment() {
			return fallback
		}
		panic(fmt.Errorf("undefined environment variable %s", key))
	}
	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err)
	}
	return int(n)
}

func GetEnvBool(key string, fallback bool) bool {
	value := strings.ToLower(os.Getenv(key))
	if (value == "true") || (value == "false") {
		return value == "true"
	}
	if IsDevelopmentEnvironment() {
		return fallback
	}
	panic(fmt.Errorf("undefined environment variable %s", key))
}

// Environment returns the current running environment
func Environment() string {
	environment := os.Getenv(EnvironmentKey)
	if environment == Production || environment == Staging || environment == Testing {
		return environment
	}
	return Development
}

// IsDevelopmentEnvironment checks if current running environment is development or not
func IsDevelopmentEnvironment() bool {
	return Environment() == Development
}
