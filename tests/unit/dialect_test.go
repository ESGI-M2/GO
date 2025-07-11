package unit

import (
	"testing"

	"project/orm/core/interfaces"
	"project/orm/dialect"
)

func setupMockDialect() *dialect.MockDialect {
	return &dialect.MockDialect{}
}

func TestMockDialect_Connect(t *testing.T) {
	mock := setupMockDialect()

	config := interfaces.ConnectionConfig{
		Host:     "localhost",
		Port:     3306,
		Database: "test",
		Username: "root",
		Password: "password",
	}

	err := mock.Connect(config)
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
}

func TestMockDialect_Close(t *testing.T) {
	mock := setupMockDialect()

	config := interfaces.ConnectionConfig{
		Host:     "localhost",
		Port:     3306,
		Database: "test",
		Username: "root",
		Password: "password",
	}

	// Connect first
	err := mock.Connect(config)
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	// Then close
	err = mock.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}
}

func TestMockDialect_Query(t *testing.T) {
	mock := setupMockDialect()

	config := interfaces.ConnectionConfig{
		Host:     "localhost",
		Port:     3306,
		Database: "test",
		Username: "root",
		Password: "password",
	}

	// Connect first
	err := mock.Connect(config)
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	// Execute query
	rows, err := mock.Query("SELECT * FROM users")
	if err != nil {
		t.Errorf("Query failed: %v", err)
	}

	// rows can be nil for empty results, which is acceptable
	_ = rows
}

func TestMockDialect_Exec(t *testing.T) {
	mock := setupMockDialect()

	config := interfaces.ConnectionConfig{
		Host:     "localhost",
		Port:     3306,
		Database: "test",
		Username: "root",
		Password: "password",
	}

	// Connect first
	err := mock.Connect(config)
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	// Execute statement
	result, err := mock.Exec("INSERT INTO users (name) VALUES (?)", "test")
	if err != nil {
		t.Errorf("Exec failed: %v", err)
	}

	if result == nil {
		t.Error("Exec should return result")
	}
}

func TestMockDialect_Begin(t *testing.T) {
	mock := setupMockDialect()

	config := interfaces.ConnectionConfig{
		Host:     "localhost",
		Port:     3306,
		Database: "test",
		Username: "root",
		Password: "password",
	}

	// Connect first
	err := mock.Connect(config)
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	// Begin transaction
	tx, err := mock.Begin()
	if err != nil {
		t.Errorf("Begin failed: %v", err)
	}

	if tx == nil {
		t.Error("Begin should return transaction")
	}
}

func TestMockDialect_Commit(t *testing.T) {
	mock := setupMockDialect()

	config := interfaces.ConnectionConfig{
		Host:     "localhost",
		Port:     3306,
		Database: "test",
		Username: "root",
		Password: "password",
	}

	// Connect first
	err := mock.Connect(config)
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	// Begin transaction
	tx, err := mock.Begin()
	if err != nil {
		t.Errorf("Begin failed: %v", err)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		t.Errorf("Commit failed: %v", err)
	}
}

func TestMockDialect_Rollback(t *testing.T) {
	mock := setupMockDialect()

	config := interfaces.ConnectionConfig{
		Host:     "localhost",
		Port:     3306,
		Database: "test",
		Username: "root",
		Password: "password",
	}

	// Connect first
	err := mock.Connect(config)
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	// Begin transaction
	tx, err := mock.Begin()
	if err != nil {
		t.Errorf("Begin failed: %v", err)
	}

	// Rollback transaction
	err = tx.Rollback()
	if err != nil {
		t.Errorf("Rollback failed: %v", err)
	}
}

func TestMockDialect_GetLastInsertID(t *testing.T) {
	mock := setupMockDialect()

	config := interfaces.ConnectionConfig{
		Host:     "localhost",
		Port:     3306,
		Database: "test",
		Username: "root",
		Password: "password",
	}

	// Connect first
	err := mock.Connect(config)
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	// Execute insert
	result, err := mock.Exec("INSERT INTO users (name) VALUES (?)", "test")
	if err != nil {
		t.Errorf("Exec failed: %v", err)
	}

	// Get last insert ID
	id, err := result.LastInsertId()
	if err != nil {
		t.Errorf("LastInsertId failed: %v", err)
	}

	if id < 0 {
		t.Error("LastInsertId should be non-negative")
	}
}

func TestMockDialect_GetRowsAffected(t *testing.T) {
	mock := setupMockDialect()

	config := interfaces.ConnectionConfig{
		Host:     "localhost",
		Port:     3306,
		Database: "test",
		Username: "root",
		Password: "password",
	}

	// Connect first
	err := mock.Connect(config)
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	// Execute update
	result, err := mock.Exec("UPDATE users SET name = ? WHERE id = ?", "updated", 1)
	if err != nil {
		t.Errorf("Exec failed: %v", err)
	}

	// Get rows affected
	affected, err := result.RowsAffected()
	if err != nil {
		t.Errorf("RowsAffected failed: %v", err)
	}

	if affected < 0 {
		t.Error("RowsAffected should be non-negative")
	}
}

func TestMockDialect_ErrorCases(t *testing.T) {
	mock := setupMockDialect()

	// Test query without connecting
	_, err := mock.Query("SELECT * FROM users")
	if err == nil {
		t.Error("Query should fail when not connected")
	}

	// Test exec without connecting
	_, err = mock.Exec("INSERT INTO users (name) VALUES (?)", "test")
	if err == nil {
		t.Error("Exec should fail when not connected")
	}

	// Test begin without connecting
	_, err = mock.Begin()
	if err == nil {
		t.Error("Begin should fail when not connected")
	}
}

func TestMockDialect_TransactionErrorCases(t *testing.T) {
	mock := setupMockDialect()

	config := interfaces.ConnectionConfig{
		Host:     "localhost",
		Port:     3306,
		Database: "test",
		Username: "root",
		Password: "password",
	}

	// Connect first
	err := mock.Connect(config)
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	// Begin transaction
	tx, err := mock.Begin()
	if err != nil {
		t.Errorf("Begin failed: %v", err)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		t.Errorf("Commit failed: %v", err)
	}

	// Try to commit again (should fail)
	err = tx.Commit()
	if err == nil {
		t.Error("Commit should fail after already committed")
	}

	// Try to rollback after commit (should fail)
	err = tx.Rollback()
	if err == nil {
		t.Error("Rollback should fail after commit")
	}
}

func TestMockDialect_ResultErrorCases(t *testing.T) {
	mock := setupMockDialect()

	config := interfaces.ConnectionConfig{
		Host:     "localhost",
		Port:     3306,
		Database: "test",
		Username: "root",
		Password: "password",
	}

	// Connect first
	err := mock.Connect(config)
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	// Execute query (not insert/update)
	result, err := mock.Exec("SELECT * FROM users")
	if err != nil {
		t.Errorf("Exec failed: %v", err)
	}

	// LastInsertId should fail for SELECT
	_, err = result.LastInsertId()
	if err == nil {
		t.Error("LastInsertId should fail for SELECT")
	}

	// RowsAffected should be 0 for SELECT
	affected, err := result.RowsAffected()
	if err != nil {
		t.Errorf("RowsAffected failed: %v", err)
	}
	if affected != 0 {
		t.Errorf("Expected 0 rows affected for SELECT, got %d", affected)
	}
}
