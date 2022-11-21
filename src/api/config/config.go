package config

import "bitbucket.org/ziggy192/ng_lu/src/util"

type Config struct {
	AuthSecret             string
	Port                   string
	MySQL                  *MySQLConfig
	Redis                  *RedisConfig
	JWTExpiresAfterMinutes int
}

func New() *Config {
	return &Config{
		AuthSecret:             util.GetEnvString("AUTH_SECRET", "test_secret"),
		Port:                   util.GetEnvString("PORT", "8080"),
		MySQL:                  NewMySQLConfig(),
		Redis:                  NewRedisConfig(),
		JWTExpiresAfterMinutes: util.GetEnvInt("JWT_EXPIRES_AFTER_MINUTES", 300),
	}
}

// MySQLConfig contains config data to connect to MySQL database
type MySQLConfig struct {
	Server                    string
	Schema                    string
	User                      string
	Password                  string
	Option                    string
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
		Option:                    util.GetEnvString("DB_OPTION", ""),
		ConnectionLifetimeSeconds: util.GetEnvInt("DB_CONNECTION_LIFETIME_SECONDS", 300),
		MaxIdleConnections:        util.GetEnvInt("DB_MAX_IDLE_CONNECTIONS", 10),
		MaxOpenConnections:        util.GetEnvInt("DB_MAX_OPEN_CONNECTIONS", 20),
	}
}

// RedisConfig contains configuration to connect to redis
type RedisConfig struct {
	Addr      string
	Password  string
	DB        int
	PoolSize  int
	EnableTLS bool
}

func NewRedisConfig() *RedisConfig {
	return &RedisConfig{
		Addr:      util.GetEnvString("REDIS_ADDRESS", "localhost:6379"),
		Password:  util.GetEnvString("REDIS_PASSWORD", "password"),
		DB:        util.GetEnvInt("REDIS_DB", 0),
		PoolSize:  util.GetEnvInt("REDIS_POOL_SIZE", 32),
		EnableTLS: util.GetEnvBool("REDIS_ENABLE_TLS", false),
	}
}
