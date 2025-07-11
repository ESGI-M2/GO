package dialect

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"project/orm/core/interfaces"
)

// MockDialect implements the Dialect interface for testing
type MockDialect struct {
	connected  bool
	execError  error
	queryError error
	beginError error
}

func (m *MockDialect) Connect(config interfaces.ConnectionConfig) error {
	m.connected = true
	return nil
}

func (m *MockDialect) Close() error {
	m.connected = false
	return nil
}

func (m *MockDialect) Ping() error {
	if !m.connected {
		return fmt.Errorf("not connected")
	}
	return nil
}

func (m *MockDialect) Exec(query string, args ...interface{}) (sql.Result, error) {
	if !m.connected {
		return nil, fmt.Errorf("not connected")
	}
	if m.execError != nil {
		return nil, m.execError
	}

	// Check if it's a SELECT query
	isSelect := len(query) > 6 && query[:6] == "SELECT"

	return &MockSQLResult{
		lastInsertID: 1,
		rowsAffected: 1,
		isSelect:     isSelect,
	}, nil
}

func (m *MockDialect) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if !m.connected {
		return nil, fmt.Errorf("not connected")
	}
	if m.queryError != nil {
		return nil, m.queryError
	}
	// Return nil to indicate no rows, which is safer than empty rows
	return nil, nil
}

func (m *MockDialect) QueryRow(query string, args ...interface{}) *sql.Row {
	if !m.connected {
		return nil
	}
	// Return nil to avoid panics with mock rows
	return nil
}

func (m *MockDialect) Begin() (interfaces.Transaction, error) {
	if !m.connected {
		return nil, fmt.Errorf("not connected")
	}
	if m.beginError != nil {
		return nil, m.beginError
	}
	return &MockTransaction{}, nil
}

func (m *MockDialect) BeginTx(ctx context.Context, opts *sql.TxOptions) (interfaces.Transaction, error) {
	if !m.connected {
		return nil, fmt.Errorf("not connected")
	}
	if m.beginError != nil {
		return nil, m.beginError
	}
	return &MockTransaction{}, nil
}

func (m *MockDialect) CreateTable(tableName string, columns []interfaces.Column) error {
	return nil
}

func (m *MockDialect) DropTable(tableName string) error {
	return nil
}

func (m *MockDialect) TableExists(tableName string) (bool, error) {
	return false, nil
}

func (m *MockDialect) GetSQLType(goType reflect.Type) string {
	switch goType.Kind() {
	case reflect.String:
		return "VARCHAR(255)"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "INT"
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.Float32, reflect.Float64:
		return "DOUBLE"
	default:
		if goType == reflect.TypeOf(time.Time{}) {
			return "DATETIME"
		}
		return "TEXT"
	}
}

func (m *MockDialect) GetPlaceholder(index int) string {
	return "?"
}

// MockSQLResult for testing
type MockSQLResult struct {
	lastInsertID int64
	rowsAffected int64
	isSelect     bool
}

func (m *MockSQLResult) LastInsertId() (int64, error) {
	if m.isSelect {
		return 0, fmt.Errorf("LastInsertId not supported for SELECT")
	}
	return m.lastInsertID, nil
}

func (m *MockSQLResult) RowsAffected() (int64, error) {
	if m.isSelect {
		return 0, nil
	}
	return m.rowsAffected, nil
}

// MockTransaction for testing
type MockTransaction struct {
	committed  bool
	rolledBack bool
}

func (m *MockTransaction) Commit() error {
	if m.rolledBack {
		return fmt.Errorf("already rolled back")
	}
	if m.committed {
		return fmt.Errorf("already committed")
	}
	m.committed = true
	return nil
}

func (m *MockTransaction) Rollback() error {
	if m.committed {
		return fmt.Errorf("already committed")
	}
	if m.rolledBack {
		return fmt.Errorf("already rolled back")
	}
	m.rolledBack = true
	return nil
}

func (m *MockTransaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	return &MockSQLResult{lastInsertID: 1, rowsAffected: 1}, nil
}

func (m *MockTransaction) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return &sql.Rows{}, nil
}

func (m *MockTransaction) QueryRow(query string, args ...interface{}) *sql.Row {
	// Return a dummy non-nil *sql.Row
	return &sql.Row{}
}
