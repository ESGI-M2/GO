package unit

import (
	"testing"

	"project/orm"
	"project/orm/core/interfaces"
	"project/orm/dialect"
)

func setupConnectionTest() orm.ORM {
	mockDialect := &dialect.MockDialect{}
	return orm.New(mockDialect)
}

func TestConnection_Connect(t *testing.T) {
	ormInstance := setupConnectionTest()

	config := interfaces.ConnectionConfig{
		Host:     "localhost",
		Port:     3306,
		Database: "test",
		Username: "root",
		Password: "password",
	}

	err := ormInstance.Connect(config)
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	// Test that connection is established
	if !ormInstance.IsConnected() {
		t.Error("ORM should be connected after Connect()")
	}
}

func TestConnection_Connect_AlreadyConnected(t *testing.T) {
	ormInstance := setupConnectionTest()

	config := interfaces.ConnectionConfig{
		Host:     "localhost",
		Port:     3306,
		Database: "test",
		Username: "root",
		Password: "password",
	}

	// First connection
	err := ormInstance.Connect(config)
	if err != nil {
		t.Errorf("First Connect failed: %v", err)
	}

	// Second connection should not fail
	err = ormInstance.Connect(config)
	if err != nil {
		t.Errorf("Second Connect failed: %v", err)
	}
}

func TestConnection_Close(t *testing.T) {
	ormInstance := setupConnectionTest()

	config := interfaces.ConnectionConfig{
		Host:     "localhost",
		Port:     3306,
		Database: "test",
		Username: "root",
		Password: "password",
	}

	// Connect first
	err := ormInstance.Connect(config)
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	// Then close
	err = ormInstance.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	// Test that connection is closed
	if ormInstance.IsConnected() {
		t.Error("ORM should not be connected after Close()")
	}
}

func TestConnection_Close_NotConnected(t *testing.T) {
	ormInstance := setupConnectionTest()

	// Close without connecting first
	err := ormInstance.Close()
	if err != nil {
		t.Errorf("Close when not connected failed: %v", err)
	}
}

func TestConnection_IsConnected(t *testing.T) {
	ormInstance := setupConnectionTest()

	// Initially not connected
	if ormInstance.IsConnected() {
		t.Error("ORM should not be connected initially")
	}

	config := interfaces.ConnectionConfig{
		Host:     "localhost",
		Port:     3306,
		Database: "test",
		Username: "root",
		Password: "password",
	}

	// Connect
	err := ormInstance.Connect(config)
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	// Should be connected
	if !ormInstance.IsConnected() {
		t.Error("ORM should be connected after Connect()")
	}

	// Close
	err = ormInstance.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	// Should not be connected
	if ormInstance.IsConnected() {
		t.Error("ORM should not be connected after Close()")
	}
}

func TestConnection_GetDialect(t *testing.T) {
	mockDialect := &dialect.MockDialect{}
	ormInstance := orm.New(mockDialect)

	dialect := ormInstance.GetDialect()
	if dialect == nil {
		t.Error("GetDialect should return the dialect")
	}

	if dialect != mockDialect {
		t.Error("GetDialect should return the same dialect instance")
	}
}

func TestConnection_ErrorCases(t *testing.T) {
	// Test with nil dialect
	ormInstance := orm.New(nil)

	config := interfaces.ConnectionConfig{
		Host:     "localhost",
		Port:     3306,
		Database: "test",
		Username: "root",
		Password: "password",
	}

	err := ormInstance.Connect(config)
	if err == nil {
		t.Error("Connect should fail with nil dialect")
	}

	// Close should not fail when not connected
	err = ormInstance.Close()
	if err != nil {
		t.Errorf("Close should not fail when not connected: %v", err)
	}

	if ormInstance.IsConnected() {
		t.Error("Should not be connected with nil dialect")
	}
}
