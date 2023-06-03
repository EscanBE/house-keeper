package pg_db_tx

// CommitTransaction implements DbTransaction
func (c *dbTxImpl) CommitTransaction() error {
	return c.Tx.Commit()
}

// RollbackTransaction implements DbTransaction
func (c *dbTxImpl) RollbackTransaction() error {
	return c.Tx.Rollback()
}
