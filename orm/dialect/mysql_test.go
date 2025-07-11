package dialect

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	"project/orm/core"
)

// TestDialectInterface tests the dialect interface implementation
func TestDialectInterface(t *testing.T) {
	// This test verifies that the dialect interface is properly defined
	// and can be implemented by different database drivers

	// Test interface compliance
	var _ core.Dialect = (*MockDialect)(nil)

	t.Log("Dialect interface is properly defined")
}

// MockDialect implements the Dialect interface for testing
type MockDialect struct {
	connected  bool
	execError  error
	queryError error
	beginError error
}

func (m *MockDialect) Connect(config core.ConnectionConfig) error {
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
	return &MockSQLResult{}, m.execError
}

func (m *MockDialect) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, m.queryError
}

func (m *MockDialect) QueryRow(query string, args ...interface{}) *sql.Row {
	return nil
}

func (m *MockDialect) Begin() (core.Transaction, error) {
	if m.beginError != nil {
		return nil, m.beginError
	}
	return &MockTransaction{}, nil
}

func (m *MockDialect) BeginTx(ctx context.Context, opts *sql.TxOptions) (core.Transaction, error) {
	if m.beginError != nil {
		return nil, m.beginError
	}
	return &MockTransaction{}, nil
}

func (m *MockDialect) CreateTable(tableName string, columns []core.Column) error {
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
}

func (m *MockSQLResult) LastInsertId() (int64, error) {
	return m.lastInsertID, nil
}

func (m *MockSQLResult) RowsAffected() (int64, error) {
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
	m.committed = true
	return nil
}

func (m *MockTransaction) Rollback() error {
	if m.committed {
		return fmt.Errorf("already committed")
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

func TestMockDialect_Connect(t *testing.T) {
	dialect := &MockDialect{}

	config := core.ConnectionConfig{
		Driver:   "mock",
		Host:     "localhost",
		Port:     3306,
		Database: "test",
		Username: "root",
		Password: "password",
	}

	err := dialect.Connect(config)
	if err != nil {
		t.Errorf("Connect should not return error: %v", err)
	}

	if !dialect.connected {
		t.Error("Dialect should be connected after Connect")
	}
}

func TestMockDialect_Close(t *testing.T) {
	dialect := &MockDialect{
		connected: true,
	}

	err := dialect.Close()
	if err != nil {
		t.Errorf("Close should not return error: %v", err)
	}

	if dialect.connected {
		t.Error("Dialect should not be connected after Close")
	}
}

func TestMockDialect_Ping(t *testing.T) {
	dialect := &MockDialect{
		connected: true,
	}

	err := dialect.Ping()
	if err != nil {
		t.Errorf("Ping should not return error: %v", err)
	}

	// Test when not connected
	dialect.connected = false
	err = dialect.Ping()
	if err == nil {
		t.Error("Ping should return error when not connected")
	}
}

func TestMockDialect_Exec(t *testing.T) {
	dialect := &MockDialect{}

	result, err := dialect.Exec("INSERT INTO users (name) VALUES (?)", "John")
	if err != nil {
		t.Errorf("Exec should not return error: %v", err)
	}

	if result == nil {
		t.Error("Exec should return result")
	}

	// Test with exec error
	dialect.execError = fmt.Errorf("exec error")
	_, err = dialect.Exec("INSERT INTO users (name) VALUES (?)", "John")
	if err == nil {
		t.Error("Exec should return error when exec fails")
	}
}

func TestMockDialect_Query(t *testing.T) {
	dialect := &MockDialect{}

	_, err := dialect.Query("SELECT * FROM users")
	if err != nil {
		t.Errorf("Query should not return error: %v", err)
	}

	// Test with query error
	dialect.queryError = fmt.Errorf("query error")
	_, err = dialect.Query("SELECT * FROM users")
	if err == nil {
		t.Error("Query should return error when query fails")
	}
}

func TestMockDialect_Begin(t *testing.T) {
	dialect := &MockDialect{}

	tx, err := dialect.Begin()
	if err != nil {
		t.Errorf("Begin should not return error: %v", err)
	}

	if tx == nil {
		t.Error("Begin should return transaction")
	}

	// Test with begin error
	dialect.beginError = fmt.Errorf("begin error")
	_, err = dialect.Begin()
	if err == nil {
		t.Error("Begin should return error when begin fails")
	}
}

func TestMockDialect_BeginTx(t *testing.T) {
	dialect := &MockDialect{}

	ctx := context.Background()
	opts := &sql.TxOptions{}

	tx, err := dialect.BeginTx(ctx, opts)
	if err != nil {
		t.Errorf("BeginTx should not return error: %v", err)
	}

	if tx == nil {
		t.Error("BeginTx should return transaction")
	}

	// Test with begin error
	dialect.beginError = fmt.Errorf("begin error")
	_, err = dialect.BeginTx(ctx, opts)
	if err == nil {
		t.Error("BeginTx should return error when begin fails")
	}
}

func TestMockDialect_CreateTable(t *testing.T) {
	dialect := &MockDialect{}

	columns := []core.Column{
		{Name: "id", Type: "INT", PrimaryKey: true, AutoIncrement: true},
		{Name: "name", Type: "VARCHAR(255)"},
		{Name: "age", Type: "INT"},
	}

	err := dialect.CreateTable("users", columns)
	if err != nil {
		t.Errorf("CreateTable should not return error: %v", err)
	}
}

func TestMockDialect_DropTable(t *testing.T) {
	dialect := &MockDialect{}

	err := dialect.DropTable("users")
	if err != nil {
		t.Errorf("DropTable should not return error: %v", err)
	}
}

func TestMockDialect_TableExists(t *testing.T) {
	dialect := &MockDialect{}

	exists, err := dialect.TableExists("users")
	if err != nil {
		t.Errorf("TableExists should not return error: %v", err)
	}

	// Mock dialect returns false by default
	if exists {
		t.Error("TableExists should return false for non-existent table")
	}
}

func TestMockDialect_GetSQLType(t *testing.T) {
	dialect := &MockDialect{}

	// Test string type
	sqlType := dialect.GetSQLType(reflect.TypeOf(""))
	if sqlType != "VARCHAR(255)" {
		t.Errorf("Expected VARCHAR(255) for string, got %s", sqlType)
	}

	// Test int type
	sqlType = dialect.GetSQLType(reflect.TypeOf(0))
	if sqlType != "INT" {
		t.Errorf("Expected INT for int, got %s", sqlType)
	}

	// Test bool type
	sqlType = dialect.GetSQLType(reflect.TypeOf(true))
	if sqlType != "BOOLEAN" {
		t.Errorf("Expected BOOLEAN for bool, got %s", sqlType)
	}

	// Test float64 type
	sqlType = dialect.GetSQLType(reflect.TypeOf(0.0))
	if sqlType != "DOUBLE" {
		t.Errorf("Expected DOUBLE for float64, got %s", sqlType)
	}

	// Test time.Time type
	sqlType = dialect.GetSQLType(reflect.TypeOf(time.Time{}))
	if sqlType != "DATETIME" {
		t.Errorf("Expected DATETIME for time.Time, got %s", sqlType)
	}

	// Test unknown type
	sqlType = dialect.GetSQLType(reflect.TypeOf([]string{}))
	if sqlType != "TEXT" {
		t.Errorf("Expected TEXT for unknown type, got %s", sqlType)
	}
}

func TestMockDialect_GetPlaceholder(t *testing.T) {
	dialect := &MockDialect{}

	placeholder := dialect.GetPlaceholder(1)
	if placeholder != "?" {
		t.Errorf("Expected ?, got %s", placeholder)
	}

	placeholder = dialect.GetPlaceholder(5)
	if placeholder != "?" {
		t.Errorf("Expected ?, got %s", placeholder)
	}
}

func TestMockTransaction_Methods(t *testing.T) {
	tx := &MockTransaction{}

	// Test Commit
	err := tx.Commit()
	if err != nil {
		t.Errorf("Commit should not return error: %v", err)
	}

	if !tx.committed {
		t.Error("Transaction should be marked as committed")
	}

	// Test Rollback on committed transaction
	err = tx.Rollback()
	if err == nil {
		t.Error("Rollback should return error when already committed")
	}

	// Test new transaction
	tx2 := &MockTransaction{}

	// Test Rollback
	err = tx2.Rollback()
	if err != nil {
		t.Errorf("Rollback should not return error: %v", err)
	}

	if !tx2.rolledBack {
		t.Error("Transaction should be marked as rolled back")
	}

	// Test Commit on rolled back transaction
	err = tx2.Commit()
	if err == nil {
		t.Error("Commit should return error when already rolled back")
	}

	// Test Exec
	result, err := tx2.Exec("INSERT INTO users (name) VALUES (?)", "John")
	if err != nil {
		t.Errorf("Exec should not return error: %v", err)
	}

	if result == nil {
		t.Error("Exec should return result")
	}

	// Test Query
	_, err = tx2.Query("SELECT * FROM users")
	if err != nil {
		t.Errorf("Query should not return error: %v", err)
	}

	// Test QueryRow
	row := tx2.QueryRow("SELECT * FROM users WHERE id = ?", 1)
	if row == nil {
		t.Error("QueryRow should return row")
	}
}

func TestMockSQLResult_Methods(t *testing.T) {
	result := &MockSQLResult{
		lastInsertID: 1,
		rowsAffected: 5,
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		t.Errorf("LastInsertId should not return error: %v", err)
	}

	if lastID != 1 {
		t.Errorf("Expected LastInsertId 1, got %d", lastID)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Errorf("RowsAffected should not return error: %v", err)
	}

	if rowsAffected != 5 {
		t.Errorf("Expected RowsAffected 5, got %d", rowsAffected)
	}
}
