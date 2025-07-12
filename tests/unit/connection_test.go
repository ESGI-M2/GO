package unit

import (
	"testing"

	"github.com/ESGI-M2/GO/orm/builder"
	"github.com/ESGI-M2/GO/orm/factory"
)

func setupConnectionTest() *builder.SimpleORM {
	return builder.NewSimpleORM().WithDialect(factory.Mock)
}

func TestConnection_Connect(t *testing.T) {
	orm := setupConnectionTest()

	// Use the new approach with quick config
	orm.WithQuickConfig("localhost", "test", "root", "password")

	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	// Test that connection is established
	if !orm.IsConnected() {
		t.Error("ORM should be connected after Connect()")
	}

	// Clean up
	orm.Close()
}

func TestConnection_Connect_AlreadyConnected(t *testing.T) {
	orm := setupConnectionTest()

	// Use the new approach with quick config
	orm.WithQuickConfig("localhost", "test", "root", "password")

	// First connection
	err := orm.Connect()
	if err != nil {
		t.Errorf("First Connect failed: %v", err)
	}

	// Second connection should not fail
	err = orm.Connect()
	if err != nil {
		t.Errorf("Second Connect failed: %v", err)
	}

	// Clean up
	orm.Close()
}

func TestConnection_Close(t *testing.T) {
	orm := setupConnectionTest()

	// Use the new approach with quick config
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

	// Test that connection is closed
	if orm.IsConnected() {
		t.Error("ORM should not be connected after Close()")
	}
}

func TestConnection_Close_NotConnected(t *testing.T) {
	orm := setupConnectionTest()

	// Close without connecting first
	err := orm.Close()
	if err != nil {
		t.Errorf("Close when not connected failed: %v", err)
	}
}

func TestConnection_IsConnected(t *testing.T) {
	orm := setupConnectionTest()

	// Initially not connected
	if orm.IsConnected() {
		t.Error("ORM should not be connected initially")
	}

	// Use the new approach with quick config
	orm.WithQuickConfig("localhost", "test", "root", "password")

	// Connect
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	// Should be connected
	if !orm.IsConnected() {
		t.Error("ORM should be connected after Connect()")
	}

	// Close
	err = orm.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	// Should not be connected
	if orm.IsConnected() {
		t.Error("ORM should not be connected after Close()")
	}
}

func TestConnection_GetDialect(t *testing.T) {
	orm := setupConnectionTest()

	// Connect first
	orm.WithQuickConfig("localhost", "test", "root", "password")
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	// Get dialect type
	dialectType := orm.GetDialectType()
	if dialectType != factory.Mock {
		t.Errorf("Expected dialect type Mock, got %v", dialectType)
	}

	// Clean up
	orm.Close()
}

func TestConnection_ErrorCases(t *testing.T) {
	// Test with invalid dialect
	orm := builder.NewSimpleORM().WithDialect("invalid")

	err := orm.Connect()
	if err == nil {
		t.Error("Connect should fail with invalid dialect")
	}

	// Should not be connected
	if orm.IsConnected() {
		t.Error("Should not be connected with invalid dialect")
	}
}
