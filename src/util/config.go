package util

import (
	"fmt"
	"os"
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
