package store

import (
	"bitbucket.org/ziggy192/ng_lu/src/api/config"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

const (
	mysqlDriver = "mysql"
	mysqlOption = "charset=utf8&parseTime=True&loc=Local&multiStatements=True&maxAllowedPacket=0"

	dropDBPath = "schema/drop_db.sql"
	initDBPath = "schema/init_db.sql"
)

// DBStores defines struct for stores
type DBStores struct {
	DB     *sql.DB
	Config *config.MySQLConfig

	Account *AccountStore
	Profile *ProfileStore
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
		Profile: NewProfileStore(database),
	}, nil
}

// ConnectMySQL setups connections to MySQL database
func ConnectMySQL(ctx context.Context, config *config.MySQLConfig) (*sql.DB, error) {
	option := config.Option
	if len(option) == 0 {
		option = mysqlOption
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", config.User, config.Password,
		config.Server, config.Schema, option)

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

func (s *DBStores) Ping(ctx context.Context) error {
	return s.DB.PingContext(ctx)
}

// Close current db connection
func (s *DBStores) Close() error {
	return s.DB.Close()
}

// Reset drops the databases and recreates it, only for testing purposes
// Must not be used in production
func (s *DBStores) Reset(ctx context.Context) error {
	err := ExecuteSQLFile(ctx, s.DB, dropDBPath)
	if err != nil {
		return err
	}

	return ExecuteSQLFile(ctx, s.DB, initDBPath)
}

// ExecuteSQLFile executes a specific file on a MySQL database
func ExecuteSQLFile(ctx context.Context, db *sql.DB, filePath string) error {
	migrationSQL, err := os.ReadFile(filePath)
	if err != nil {
		logger.Err(ctx, err)
		return err
	}

	_, err = db.ExecContext(ctx, string(migrationSQL))
	if err != nil {
		logger.Err(ctx, err)
		return err
	}

	logger.Info(ctx, "migrated", filePath)
	return nil
}
