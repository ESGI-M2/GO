package unit

import (
	"testing"

	"github.com/ESGI-M2/GO/orm"
	"github.com/ESGI-M2/GO/orm/dialect"
)

// TestUser represents a test user model
type TestUser struct {
	ID    int    `orm:"pk,auto"`
	Name  string `orm:"column:name"`
	Email string `orm:"column:email,unique"`
}

// TestBasicORM tests basic ORM functionality
func TestBasicORM(t *testing.T) {
	// Create a mock dialect
	mockDialect := &dialect.MockDialect{}

	// Create ORM instance
	ormInstance := orm.New(mockDialect)

	// Test that ORM is created successfully
	if ormInstance == nil {
		t.Fatal("ORM instance should not be nil")
	}

	// Test that ORM is not connected initially
	if ormInstance.IsConnected() {
		t.Error("ORM should not be connected initially")
	}

	// Test model registration
	user := &TestUser{}
	err := ormInstance.RegisterModel(user)
	if err != nil {
		t.Errorf("Failed to register model: %v", err)
	}

	// Test getting metadata
	metadata, err := ormInstance.GetMetadata(user)
	if err != nil {
		t.Errorf("Failed to get metadata: %v", err)
	}

	if metadata == nil {
		t.Fatal("Metadata should not be nil")
	}

	if metadata.TableName != "testuser" {
		t.Errorf("Expected table name 'testuser', got '%s'", metadata.TableName)
	}

	// Test query builder creation
	query := ormInstance.Query(user)
	if query == nil {
		t.Fatal("Query builder should not be nil")
	}

	// Test repository creation
	repo := ormInstance.Repository(user)
	if repo == nil {
		t.Fatal("Repository should not be nil")
	}
}

// TestConnection tests connection functionality
func TestConnection(t *testing.T) {
	mockDialect := &dialect.MockDialect{}
	ormInstance := orm.New(mockDialect)

	config := orm.ConnectionConfig{
		Host:     "localhost",
		Port:     3306,
		Username: "test",
		Password: "test",
		Database: "test",
	}

	err := ormInstance.Connect(config)
	if err != nil {
		t.Errorf("Failed to connect: %v", err)
	}

	if !ormInstance.IsConnected() {
		t.Error("ORM should be connected after Connect()")
	}

	err = ormInstance.Close()
	if err != nil {
		t.Errorf("Failed to close: %v", err)
	}

	if ormInstance.IsConnected() {
		t.Error("ORM should not be connected after Close()")
	}
}
