package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/ESGI-M2/GO/orm/core/connection"
	"github.com/ESGI-M2/GO/orm/core/interfaces"
)

// TransactionDialect wraps a transaction to implement the Dialect interface
type TransactionDialect struct {
	tx interfaces.Transaction
}

// Connect is a no-op for transaction dialect
func (td *TransactionDialect) Connect(config interfaces.ConnectionConfig) error {
	return nil
}

// Close is a no-op for transaction dialect
func (td *TransactionDialect) Close() error {
	return nil
}

// Ping is a no-op for transaction dialect
func (td *TransactionDialect) Ping() error {
	return nil
}

// Exec executes a query on the transaction
func (td *TransactionDialect) Exec(query string, args ...interface{}) (sql.Result, error) {
	return td.tx.Exec(query, args...)
}

// Query executes a query on the transaction
func (td *TransactionDialect) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return td.tx.Query(query, args...)
}

// QueryRow executes a query on the transaction
func (td *TransactionDialect) QueryRow(query string, args ...interface{}) *sql.Row {
	return td.tx.QueryRow(query, args...)
}

// Begin is not supported for transaction dialect
func (td *TransactionDialect) Begin() (interfaces.Transaction, error) {
	return nil, fmt.Errorf("nested transactions not supported")
}

// BeginTx is not supported for transaction dialect
func (td *TransactionDialect) BeginTx(ctx context.Context, opts *sql.TxOptions) (interfaces.Transaction, error) {
	return nil, fmt.Errorf("nested transactions not supported")
}

// CreateTable is not supported for transaction dialect
func (td *TransactionDialect) CreateTable(tableName string, columns []interfaces.Column) error {
	return fmt.Errorf("create table not supported in transaction")
}

// DropTable is not supported for transaction dialect
func (td *TransactionDialect) DropTable(tableName string) error {
	return fmt.Errorf("drop table not supported in transaction")
}

// TableExists is not supported for transaction dialect
func (td *TransactionDialect) TableExists(tableName string) (bool, error) {
	return false, fmt.Errorf("table exists not supported in transaction")
}

// GetSQLType is not supported for transaction dialect
func (td *TransactionDialect) GetSQLType(goType reflect.Type) string {
	return ""
}

// GetPlaceholder is not supported for transaction dialect
func (td *TransactionDialect) GetPlaceholder(index int) string {
	return "?"
}

// Transaction executes a function within a transaction
func Transaction(orm *connection.ORMImpl, fn func(interfaces.ORM) error) error {
	tx, err := orm.GetDialect().Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Create a transaction-scoped ORM
	txORM := &connection.ORMImpl{
		Dialect:         &TransactionDialect{tx: tx},
		MetadataManager: orm.MetadataManager,
		Models:          orm.Models,
		Connected:       true,
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(txORM); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %w, rollback failed: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransactionWithContext executes a function within a transaction with context
func TransactionWithContext(orm *connection.ORMImpl, ctx context.Context, fn func(interfaces.ORM) error) error {
	tx, err := orm.GetDialect().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Create a transaction-scoped ORM
	txORM := &connection.ORMImpl{
		Dialect:         &TransactionDialect{tx: tx},
		MetadataManager: orm.MetadataManager,
		Models:          orm.Models,
		Connected:       true,
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(txORM); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %w, rollback failed: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
