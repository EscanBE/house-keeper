package database

// DbTransaction represents an abstract database transaction that can be used to CRUD, then commit or rollback on purpose
type DbTransaction interface {
	// CommitTransaction commits the current transaction instance, returns error if any problem happened using commit progress
	CommitTransaction() error
	// RollbackTransaction rollbacks the current transaction instance, returns error if any problem happened using rollback progress
	RollbackTransaction() error
}
