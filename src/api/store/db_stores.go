package store

import (
	"bitbucket.org/ziggy192/ng_lu/src/api/config"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

const (
	mysqlDriver = "mysql"
	mysqlOption = "charset=utf8&parseTime=True&loc=Local&multiStatements=True&maxAllowedPacket=0"
)

// DBStores defines struct for stores
type DBStores struct {
	DB     *sql.DB
	Config *config.MySQLConfig

	Account *AccountStore
}

// NewDBStores returns a new instance of DBStores
func NewDBStores(ctx context.Context, config *config.MySQLConfig) (*DBStores, error) {
	database, err := ConnectMySQL(ctx, config)
	if err != nil {
		logger.Err(ctx, err)
		return nil, err
	}

	return &DBStores{
		DB:      database,
		Config:  config,
		Account: NewAccountStore(database),
	}, nil
}

// ConnectMySQL setups connections to MySQL database
func ConnectMySQL(ctx context.Context, config *config.MySQLConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", config.User, config.Password,
		config.Server, config.Schema, mysqlOption)

	db, err := sql.Open(mysqlDriver, dsn)
	if err != nil {
		logger.Err(ctx, "cannot connect to database", config.Schema)
		return nil, err
	}

	db.SetConnMaxLifetime(time.Duration(config.ConnectionLifetimeSeconds) * time.Second)
	db.SetMaxIdleConns(config.MaxIdleConnections)
	db.SetMaxOpenConns(config.MaxOpenConnections)

	logger.Info(ctx, "connected to database", config.Schema)
	return db, nil
}

func (d *DBStores) Ping(ctx context.Context) error {
	return d.DB.PingContext(ctx)
}

// Close current db connection
func (s *DBStores) Close() error {
	return s.DB.Close()
}
