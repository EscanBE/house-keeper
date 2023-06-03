package database

import "context"

// Database represents an abstract database that can be used to CRUD
type Database interface {
	// BeginDatabaseTransaction will start a database transaction, then you can perform CRUD on it, finally can do commit or rollback on purpose
	BeginDatabaseTransaction(ctx context.Context) (DbTransaction, error)

	// CommitTransaction performs committing the db transaction into store
	CommitTransaction(ptx DbTransaction) error

	// RollbackTransaction performs rolling back changes made by the db transaction
	RollbackTransaction(ptx DbTransaction) error

	// Close closes the underlying database connection
	Close()
}
