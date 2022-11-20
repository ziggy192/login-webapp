package config

import (
	"bitbucket.org/ziggy192/ng_lu/src/util"
)

const defaultApiRoot = "http://localhost:8080"

type Config struct {
	APIRoot string
	Port    string
}

func New() *Config {
	return &Config{
		APIRoot: util.GetEnvString("API_ROOT", defaultApiRoot),
		Port:    util.GetEnvString("PORT", "9090"),
	}
}
