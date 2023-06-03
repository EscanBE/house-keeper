package pg_db_tx

import (
	"context"
	"database/sql"
	"github.com/EscanBE/go-app-name/database"
)

// type check to ensure interface is properly implemented
var _ database.DbTransaction = &dbTxImpl{}

// dbTxImpl defines a wrapper around a Postgres SQL database transaction and implements functionality
type dbTxImpl struct {
	ctx         context.Context
	Tx          *sql.Tx
	PartitionID int64
}

// BeginDatabaseTransaction init and return db transaction implementation
func BeginDatabaseTransaction(ctx context.Context, tx *sql.Tx) *dbTxImpl {
	return &dbTxImpl{
		ctx: ctx,
		Tx:  tx,
	}
}

// ExecWithContext receives query and it's arguments, execute it as a part of the db transaction
func (c *dbTxImpl) ExecWithContext(query string, args ...any) (sql.Result, error) {
	return c.Tx.ExecContext(c.ctx, query, args...)
}

// QueryWithContext receives query and it's arguments, execute query and return results from db transaction context
func (c *dbTxImpl) QueryWithContext(query string, args ...any) (*sql.Rows, error) {
	return c.Tx.QueryContext(c.ctx, query, args...)
}
