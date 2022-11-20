package config

import "bitbucket.org/ziggy192/ng_lu/src/util"

type Config struct {
	AuthSecret string
	Port       string
}

func New() *Config {
	return &Config{
		AuthSecret: util.GetEnvString("AUTH_SECRET", "test_secret"),
		Port:       util.GetEnvString("PORT", "8080"),
	}
}
