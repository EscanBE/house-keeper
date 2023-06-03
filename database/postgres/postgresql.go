package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/EscanBE/go-app-name/database"
	dbtx "github.com/EscanBE/go-app-name/database/postgres_tx"
	libdbtypes "github.com/EscanBE/go-lib/database/types"
	"github.com/EscanBE/go-lib/logging"
)

// type check to ensure interface is properly implemented
var _ database.Database = &Database{}

// Database defines a wrapper around a Postgres SQL database and implements functionality
type Database struct {
	Config libdbtypes.PostgresDatabaseConfig
	Sql    *sql.DB
	Logger logging.Logger
}

// NewPostgresDatabase creates a database connection with the given database connection info
// from config. It returns a database connection handle or an error if the
// connection fails.
func NewPostgresDatabase(dbCfg libdbtypes.PostgresDatabaseConfig, logger logging.Logger) (*Database, error) {
	sslMode := "disable"
	if dbCfg.EnableSsl {
		sslMode = "enable"
	}

	schema := "public"
	if dbCfg.Schema != "" {
		schema = dbCfg.Schema
	}

	connStr := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s search_path=%s",
		dbCfg.Host, dbCfg.Port, dbCfg.Name, dbCfg.Username, dbCfg.Password, sslMode, schema,
	)

	postgresDb, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Set max open connections
	postgresDb.SetMaxOpenConns(int(dbCfg.MaxOpenConnectionCount))
	postgresDb.SetMaxIdleConns(int(dbCfg.MaxIdleConnectionCount))

	return &Database{
		Sql:    postgresDb,
		Logger: logger,
	}, nil
}

// BeginDatabaseTransaction implements Database
func (db *Database) BeginDatabaseTransaction(ctx context.Context) (database.DbTransaction, error) {
	tx, err := db.Sql.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return dbtx.BeginDatabaseTransaction(ctx, tx), nil
}

// CommitTransaction implements Database
func (db *Database) CommitTransaction(ptx database.DbTransaction) error {
	return ptx.CommitTransaction()
}

// RollbackTransaction implements Database
func (db *Database) RollbackTransaction(ptx database.DbTransaction) error {
	err := ptx.RollbackTransaction()
	if err != nil {
		db.Logger.Error("problem when rollback pipe tx", "err", err.Error())
		return err
	}
	return nil
}

// Close implements Database
func (db *Database) Close() {
	err := db.Sql.Close()
	if err != nil {
		db.Logger.Error("error while closing connection", "err", err)
	}
}
