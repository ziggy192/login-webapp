package config

import "bitbucket.org/ziggy192/ng_lu/src/util"

type Config struct {
	AuthSecret string
	Port       string
	MySQL      *MySQLConfig
}

func New() *Config {
	return &Config{
		AuthSecret: util.GetEnvString("AUTH_SECRET", "test_secret"),
		Port:       util.GetEnvString("PORT", "8080"),
		MySQL:      NewMySQLConfig(),
	}
}

// MySQLConfig contains config data to connect to MySQL database
type MySQLConfig struct {
	Server                    string
	Schema                    string
	User                      string
	Password                  string
	ConnectionLifetimeSeconds int
	MaxIdleConnections        int
	MaxOpenConnections        int
}

func NewMySQLConfig() *MySQLConfig {
	return &MySQLConfig{
		Server:                    util.GetEnvString("DB_SERVER", "0.0.0.0:3306"),
		Schema:                    util.GetEnvString("DB_SCHEMA", "ng_lu"),
		User:                      util.GetEnvString("DB_USER", "ng_lu"),
		Password:                  util.GetEnvString("DB_PASSWORD", "password"),
		ConnectionLifetimeSeconds: util.GetEnvInt("DB_CONNECTION_LIFETIME_SECONDS", 300),
		MaxIdleConnections:        util.GetEnvInt("DB_MAX_IDLE_CONNECTIONS", 10),
		MaxOpenConnections:        util.GetEnvInt("DB_MAX_OPEN_CONNECTIONS", 20),
	}
}
