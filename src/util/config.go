package util

import (
	"fmt"
	"os"
	"strconv"
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
	apiRoot := os.Getenv(key)
	if len(apiRoot) == 0 {
		if IsDevelopmentEnvironment() {
			return fallback
		}
		panic(fmt.Errorf("undefined environment variable %s", key))
	}
	return apiRoot
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
