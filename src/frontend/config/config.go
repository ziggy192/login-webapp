package config

import (
	"bitbucket.org/ziggy192/ng_lu/src/util"
)

const defaultApiRoot = "http://localhost:8080"

type Config struct {
	APIRoot        string
	Port           string
	LoginURI       string
	GoogleClientID string
}

func New() *Config {
	return &Config{
		APIRoot:        util.GetEnvString("API_ROOT", defaultApiRoot),
		Port:           util.GetEnvString("PORT", "9090"),
		LoginURI:       util.GetEnvString("LOGIN_URI", "http://localhost:9090/auth"),
		GoogleClientID: util.GetEnvString("GOOGLE_CLIENT_ID", "588338350106-u3e7ddin0njjervl05577fioq678nbi5.apps.googleusercontent.com"),
	}
}
