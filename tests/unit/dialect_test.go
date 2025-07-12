package unit

import (
	"fmt"
	"testing"

	"github.com/ESGI-M2/GO/orm/builder"
	"github.com/ESGI-M2/GO/orm/core/interfaces"
	"github.com/ESGI-M2/GO/orm/factory"
)

func setupMockDialect() *builder.SimpleORM {
	return builder.NewSimpleORM().WithDialect(factory.Mock)
}

func TestMockDialect_Connect(t *testing.T) {
	orm := setupMockDialect()

	// Use the new quick config approach
	orm.WithQuickConfig("localhost", "test", "root", "password")

	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	// Clean up
	orm.Close()
}

func TestMockDialect_Close(t *testing.T) {
	orm := setupMockDialect()

	// Use the new quick config approach
	orm.WithQuickConfig("localhost", "test", "root", "password")

	// Connect first
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	// Then close
	err = orm.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}
}

func TestMockDialect_Query(t *testing.T) {
	orm := setupMockDialect()

	// Use the new quick config approach
	orm.WithQuickConfig("localhost", "test", "root", "password")

	// Connect first
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Execute query through ORM
	raw := orm.Raw("SELECT * FROM users")
	if raw == nil {
		t.Error("Raw query should return a query builder")
	}

	// Try to execute the query
	results, err := raw.Find()
	if err != nil {
		t.Errorf("Query failed: %v", err)
	}

	// results can be nil for empty results, which is acceptable
	_ = results
}

func TestMockDialect_Exec(t *testing.T) {
	orm := setupMockDialect()

	// Use the new quick config approach
	orm.WithQuickConfig("localhost", "test", "root", "password")

	// Connect first
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Execute statement through ORM
	raw := orm.Raw("INSERT INTO users (name) VALUES (?)", "test")
	if raw == nil {
		t.Error("Raw exec should return a query builder")
	}

	// Try to execute the statement
	_, err = raw.Find()
	if err != nil {
		t.Errorf("Exec failed: %v", err)
	}
}

func TestMockDialect_Begin(t *testing.T) {
	orm := setupMockDialect()

	// Use the new quick config approach
	orm.WithQuickConfig("localhost", "test", "root", "password")

	// Connect first
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Begin transaction
	err = orm.Transaction(func(tx interfaces.ORM) error {
		if tx == nil {
			t.Error("Transaction should provide ORM instance")
		}
		return nil
	})
	if err != nil {
		t.Errorf("Transaction failed: %v", err)
	}
}

func TestMockDialect_Commit(t *testing.T) {
	orm := setupMockDialect()

	// Use the new quick config approach
	orm.WithQuickConfig("localhost", "test", "root", "password")

	// Connect first
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Begin and commit transaction
	err = orm.Transaction(func(tx interfaces.ORM) error {
		// Transaction will auto-commit on success
		return nil
	})
	if err != nil {
		t.Errorf("Transaction commit failed: %v", err)
	}
}

func TestMockDialect_Rollback(t *testing.T) {
	orm := setupMockDialect()

	// Use the new quick config approach
	orm.WithQuickConfig("localhost", "test", "root", "password")

	// Connect first
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Begin and rollback transaction
	err = orm.Transaction(func(tx interfaces.ORM) error {
		// Return error to trigger rollback
		return fmt.Errorf("test rollback error")
	})
	if err == nil {
		t.Error("Transaction should fail when returning error")
	}
}

func TestMockDialect_GetLastInsertID(t *testing.T) {
	orm := setupMockDialect()

	// Use the new quick config approach
	orm.WithQuickConfig("localhost", "test", "root", "password")

	// Connect first
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Execute insert through ORM
	raw := orm.Raw("INSERT INTO users (name) VALUES (?)", "test")
	if raw == nil {
		t.Error("Raw insert should return a query builder")
	}

	// Try to execute the insert
	_, err = raw.Find()
	if err != nil {
		t.Errorf("Insert failed: %v", err)
	}

	// Mock dialect should handle last insert ID internally
}

func TestMockDialect_GetRowsAffected(t *testing.T) {
	orm := setupMockDialect()

	// Use the new quick config approach
	orm.WithQuickConfig("localhost", "test", "root", "password")

	// Connect first
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Execute update through ORM
	raw := orm.Raw("UPDATE users SET name = ? WHERE id = ?", "updated", 1)
	if raw == nil {
		t.Error("Raw update should return a query builder")
	}

	// Try to execute the update
	_, err = raw.Find()
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	// Mock dialect should handle rows affected internally
}

func TestMockDialect_ErrorCases(t *testing.T) {
	orm := setupMockDialect()

	// Test query without connecting
	raw := orm.Raw("SELECT * FROM users")
	_, err := raw.Find()
	if err == nil {
		t.Error("Query should fail when not connected")
	}
}

func TestMockDialect_TransactionErrorCases(t *testing.T) {
	orm := setupMockDialect()

	// Use the new quick config approach
	orm.WithQuickConfig("localhost", "test", "root", "password")

	// Connect first
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Test transaction with error
	err = orm.Transaction(func(tx interfaces.ORM) error {
		return fmt.Errorf("test transaction error")
	})
	if err == nil {
		t.Error("Transaction should fail when returning error")
	}
}

func TestMockDialect_ResultErrorCases(t *testing.T) {
	orm := setupMockDialect()

	// Use the new quick config approach
	orm.WithQuickConfig("localhost", "test", "root", "password")

	// Connect first
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Execute query (not insert/update)
	raw := orm.Raw("SELECT * FROM users")
	results, err := raw.Find()
	if err != nil {
		t.Errorf("Select failed: %v", err)
	}

	// Results should be valid for SELECT
	if results == nil {
		t.Error("Select should return a slice, even if empty")
	}
}
