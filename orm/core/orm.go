package core

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"sync"
)

// ORMImpl implements the ORM interface
type ORMImpl struct {
	dialect         Dialect
	metadataManager *MetadataManager
	models          map[reflect.Type]*ModelMetadata
	connected       bool
	mu              sync.RWMutex
}

// NewORM creates a new ORM instance
func NewORM(dialect Dialect) *ORMImpl {
	return &ORMImpl{
		dialect:         dialect,
		metadataManager: NewMetadataManager(),
		models:          make(map[reflect.Type]*ModelMetadata),
		connected:       false,
	}
}

// Connect establishes a connection to the database
func (o *ORMImpl) Connect(config ConnectionConfig) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.dialect == nil {
		return fmt.Errorf("dialect is not set")
	}

	if err := o.dialect.Connect(config); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	o.connected = true
	return nil
}

// Close closes the database connection
func (o *ORMImpl) Close() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.connected {
		o.connected = false
		if o.dialect != nil {
			return o.dialect.Close()
		}
	}
	return nil
}

// RegisterModel registers a model with the ORM
func (o *ORMImpl) RegisterModel(model interface{}) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	metadata, err := o.metadataManager.ExtractMetadata(model)
	if err != nil {
		return fmt.Errorf("failed to extract metadata: %w", err)
	}

	o.models[metadata.Type] = metadata
	return nil
}

// GetMetadata returns metadata for a model
func (o *ORMImpl) GetMetadata(model interface{}) (*ModelMetadata, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return o.metadataManager.GetMetadata(model)
}

// Query creates a new query builder for the given model
func (o *ORMImpl) Query(model interface{}) QueryBuilder {
	metadata, err := o.GetMetadata(model)
	if err != nil {
		// Return a query builder that will fail on execution
		return &QueryBuilderImpl{
			orm:      o,
			metadata: nil,
			err:      err,
		}
	}

	return &QueryBuilderImpl{
		orm:      o,
		metadata: metadata,
		table:    metadata.TableName,
		fields:   []string{"*"},
		where:    make([]WhereCondition, 0),
		orderBy:  make([]OrderBy, 0),
		joins:    make([]Join, 0),
		limit:    0,
		offset:   0,
		args:     make([]interface{}, 0),
	}
}

// Raw creates a raw SQL query builder
func (o *ORMImpl) Raw(sql string, args ...interface{}) QueryBuilder {
	return &QueryBuilderImpl{
		orm:      o,
		rawSQL:   sql,
		rawArgs:  args,
		metadata: nil,
	}
}

// Repository creates a repository for the given model
func (o *ORMImpl) Repository(model interface{}) Repository {
	metadata, err := o.GetMetadata(model)
	if err != nil {
		return &RepositoryImpl{
			orm:      o,
			metadata: nil,
			err:      err,
		}
	}

	return &RepositoryImpl{
		orm:      o,
		metadata: metadata,
	}
}

// Transaction executes a function within a transaction
func (o *ORMImpl) Transaction(fn func(ORM) error) error {
	tx, err := o.dialect.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Create a transaction-scoped ORM
	txORM := &ORMImpl{
		dialect:         &TransactionDialect{tx: tx},
		metadataManager: o.metadataManager,
		models:          o.models,
		connected:       true,
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
func (o *ORMImpl) TransactionWithContext(ctx context.Context, fn func(ORM) error) error {
	tx, err := o.dialect.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Create a transaction-scoped ORM
	txORM := &ORMImpl{
		dialect:         &TransactionDialect{tx: tx},
		metadataManager: o.metadataManager,
		models:          o.models,
		connected:       true,
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

	return o.dialect.CreateTable(metadata.TableName, metadata.Columns)
}

// DropTable drops the table for the given model
func (o *ORMImpl) DropTable(model interface{}) error {
	metadata, err := o.GetMetadata(model)
	if err != nil {
		return fmt.Errorf("failed to get metadata: %w", err)
	}

	return o.dialect.DropTable(metadata.TableName)
}

// Migrate performs database migrations
func (o *ORMImpl) Migrate() error {
	o.mu.RLock()
	defer o.mu.RUnlock()

	for _, metadata := range o.models {
		exists, err := o.dialect.TableExists(metadata.TableName)
		if err != nil {
			return fmt.Errorf("failed to check if table %s exists: %w", metadata.TableName, err)
		}

		if !exists {
			if err := o.dialect.CreateTable(metadata.TableName, metadata.Columns); err != nil {
				return fmt.Errorf("failed to create table %s: %w", metadata.TableName, err)
			}
		}
	}

	return nil
}

// GetDialect returns the current dialect
func (o *ORMImpl) GetDialect() Dialect {
	return o.dialect
}

// IsConnected returns whether the ORM is connected to the database
func (o *ORMImpl) IsConnected() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.connected
}

// TransactionDialect wraps a transaction as a dialect
type TransactionDialect struct {
	tx Transaction
}

// Connect is a no-op for transaction dialect
func (td *TransactionDialect) Connect(config ConnectionConfig) error {
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

// Exec delegates to the transaction
func (td *TransactionDialect) Exec(query string, args ...interface{}) (sql.Result, error) {
	return td.tx.Exec(query, args...)
}

// Query delegates to the transaction
func (td *TransactionDialect) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return td.tx.Query(query, args...)
}

// QueryRow delegates to the transaction
func (td *TransactionDialect) QueryRow(query string, args ...interface{}) *sql.Row {
	return td.tx.QueryRow(query, args...)
}

// Begin is not supported in transaction dialect
func (td *TransactionDialect) Begin() (Transaction, error) {
	return nil, fmt.Errorf("nested transactions not supported")
}

// BeginTx is not supported in transaction dialect
func (td *TransactionDialect) BeginTx(ctx context.Context, opts *sql.TxOptions) (Transaction, error) {
	return nil, fmt.Errorf("nested transactions not supported")
}

// CreateTable is not supported in transaction dialect
func (td *TransactionDialect) CreateTable(tableName string, columns []Column) error {
	return fmt.Errorf("schema operations not supported in transactions")
}

// DropTable is not supported in transaction dialect
func (td *TransactionDialect) DropTable(tableName string) error {
	return fmt.Errorf("schema operations not supported in transactions")
}

// TableExists is not supported in transaction dialect
func (td *TransactionDialect) TableExists(tableName string) (bool, error) {
	return false, fmt.Errorf("schema operations not supported in transactions")
}

// GetSQLType delegates to the original dialect (not available in transaction)
func (td *TransactionDialect) GetSQLType(goType reflect.Type) string {
	return "TEXT" // Default fallback
}

// GetPlaceholder returns the placeholder for parameterized queries
func (td *TransactionDialect) GetPlaceholder(index int) string {
	return "?"
}
