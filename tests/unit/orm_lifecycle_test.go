package unit

import (
	"testing"

	"github.com/ESGI-M2/GO/orm/builder"
	"github.com/ESGI-M2/GO/orm/factory"
)

// TestUser represents a test user model
type TestUser struct {
	ID    int    `orm:"pk,auto"`
	Name  string `orm:"column:name"`
	Email string `orm:"column:email,unique"`
}

// TestBasicORM tests basic ORM functionality
func TestBasicORM(t *testing.T) {
	// Create ORM instance using new approach
	orm := builder.NewSimpleORM().
		WithDialect(factory.Mock).
		RegisterModel(&TestUser{})

	// Test that ORM is created successfully
	if orm == nil {
		t.Fatal("ORM instance should not be nil")
	}

	// Test that ORM is not connected initially
	if orm.IsConnected() {
		t.Error("ORM should not be connected initially")
	}

	// Connect the ORM
	if err := orm.Connect(); err != nil {
		t.Errorf("Failed to connect: %v", err)
	}

	// Test that ORM is connected after Connect()
	if !orm.IsConnected() {
		t.Error("ORM should be connected after Connect()")
	}

	// Test getting underlying ORM
	underlyingORM := orm.GetORM()
	if underlyingORM == nil {
		t.Fatal("Underlying ORM should not be nil")
	}

	// Test getting metadata
	user := &TestUser{}
	metadata, err := underlyingORM.GetMetadata(user)
	if err != nil {
		t.Errorf("Failed to get metadata: %v", err)
	}

	if metadata == nil {
		t.Fatal("Metadata should not be nil")
	}

	// Test query builder creation
	query := orm.Query(user)
	if query == nil {
		t.Fatal("Query builder should not be nil")
	}

	// Test repository creation
	repo := orm.Repository(user)
	if repo == nil {
		t.Fatal("Repository should not be nil")
	}

	// Test close
	if err := orm.Close(); err != nil {
		t.Errorf("Failed to close: %v", err)
	}

	if orm.IsConnected() {
		t.Error("ORM should not be connected after Close()")
	}
}

// TestConnection tests connection functionality
func TestConnection(t *testing.T) {
	// Create ORM instance with manual config
	orm := builder.NewSimpleORM().
		WithDialect(factory.Mock).
		WithQuickConfig("localhost", "test", "test", "test").
		RegisterModel(&TestUser{})

	// Test connection
	err := orm.Connect()
	if err != nil {
		t.Errorf("Failed to connect: %v", err)
	}

	if !orm.IsConnected() {
		t.Error("ORM should be connected after Connect()")
	}

	// Test close
	err = orm.Close()
	if err != nil {
		t.Errorf("Failed to close: %v", err)
	}

	if orm.IsConnected() {
		t.Error("ORM should not be connected after Close()")
	}
}
