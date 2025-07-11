package connection

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"sync"

	"project/orm/core/interfaces"
	"project/orm/core/metadata"
	"project/orm/core/query"
	"project/orm/core/repository"
)

// ORMImpl implements the ORM interface
type ORMImpl struct {
	Dialect         interfaces.Dialect
	MetadataManager *metadata.Manager
	Models          map[reflect.Type]*interfaces.ModelMetadata
	Connected       bool
	mu              sync.RWMutex
}

// NewORM creates a new ORM instance
func NewORM(dialect interfaces.Dialect) *ORMImpl {
	return &ORMImpl{
		Dialect:         dialect,
		MetadataManager: metadata.NewManager(),
		Models:          make(map[reflect.Type]*interfaces.ModelMetadata),
		Connected:       false,
	}
}

// Connect establishes a connection to the database
func (o *ORMImpl) Connect(config interfaces.ConnectionConfig) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.Dialect == nil {
		return fmt.Errorf("dialect is not set")
	}

	if err := o.Dialect.Connect(config); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	o.Connected = true
	return nil
}

// Close closes the database connection
func (o *ORMImpl) Close() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.Connected {
		o.Connected = false
		if o.Dialect != nil {
			return o.Dialect.Close()
		}
		return fmt.Errorf("dialect is not set")
	}
	return nil
}

// IsConnected returns whether the ORM is connected to the database
func (o *ORMImpl) IsConnected() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.Connected
}

// GetDialect returns the current dialect
func (o *ORMImpl) GetDialect() interfaces.Dialect {
	return o.Dialect
}

// RegisterModel registers a model with the ORM
func (o *ORMImpl) RegisterModel(model interface{}) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	metadata, err := o.MetadataManager.ExtractMetadata(model)
	if err != nil {
		return fmt.Errorf("failed to extract metadata: %w", err)
	}

	o.Models[metadata.Type] = metadata
	return nil
}

// GetMetadata returns metadata for a model
func (o *ORMImpl) GetMetadata(model interface{}) (*interfaces.ModelMetadata, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return o.MetadataManager.GetMetadata(model)
}

// Query creates a new query builder for the given model
func (o *ORMImpl) Query(model interface{}) interfaces.QueryBuilder {
	metadata, err := o.GetMetadata(model)
	if err != nil {
		// Return a query builder that will fail on execution
		return &query.BuilderImpl{
			Orm:      o,
			Metadata: nil,
			Err:      err,
		}
	}

	return query.NewBuilder(o, metadata)
}

// Raw creates a raw SQL query builder
func (o *ORMImpl) Raw(sql string, args ...interface{}) interfaces.QueryBuilder {
	return query.NewRawBuilder(o, sql, args...)
}

// Repository creates a repository for the given model
func (o *ORMImpl) Repository(model interface{}) interfaces.Repository {
	metadata, err := o.GetMetadata(model)
	if err != nil {
		return repository.NewErrorRepository(o, err)
	}

	return repository.NewRepository(o, metadata)
}

// Transaction executes a function within a transaction
func (o *ORMImpl) Transaction(fn func(interfaces.ORM) error) error {
	tx, err := o.Dialect.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Create a transaction-scoped ORM
	txORM := &ORMImpl{
		Dialect:         &TransactionDialect{tx: tx},
		MetadataManager: o.MetadataManager,
		Models:          o.Models,
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
func (o *ORMImpl) TransactionWithContext(ctx context.Context, fn func(interfaces.ORM) error) error {
	tx, err := o.Dialect.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Create a transaction-scoped ORM
	txORM := &ORMImpl{
		Dialect:         &TransactionDialect{tx: tx},
		MetadataManager: o.MetadataManager,
		Models:          o.Models,
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

// CreateTable creates a table for the given model
func (o *ORMImpl) CreateTable(model interface{}) error {
	metadata, err := o.GetMetadata(model)
	if err != nil {
		return fmt.Errorf("failed to get metadata: %w", err)
	}

	return o.Dialect.CreateTable(metadata.TableName, metadata.Columns)
}

// DropTable drops the table for the given model
func (o *ORMImpl) DropTable(model interface{}) error {
	metadata, err := o.GetMetadata(model)
	if err != nil {
		return fmt.Errorf("failed to get metadata: %w", err)
	}

	return o.Dialect.DropTable(metadata.TableName)
}

// Migrate performs database migrations
func (o *ORMImpl) Migrate() error {
	o.mu.RLock()
	defer o.mu.RUnlock()

	for _, metadata := range o.Models {
		exists, err := o.Dialect.TableExists(metadata.TableName)
		if err != nil {
			return fmt.Errorf("failed to check if table %s exists: %w", metadata.TableName, err)
		}

		if !exists {
			if err := o.Dialect.CreateTable(metadata.TableName, metadata.Columns); err != nil {
				return fmt.Errorf("failed to create table %s: %w", metadata.TableName, err)
			}
		}
	}

	return nil
}

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
