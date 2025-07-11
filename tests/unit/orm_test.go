package unit

import (
	"context"
	"testing"

	"github.com/ESGI-M2/GO/orm"
	"github.com/ESGI-M2/GO/orm/core/interfaces"
	"github.com/ESGI-M2/GO/orm/dialect"
)

type TestModel struct {
	ID   int    `orm:"pk,auto"`
	Name string `orm:"column:name"`
}

func setupORM() orm.ORM {
	mockDialect := &dialect.MockDialect{}
	ormInstance := orm.New(mockDialect)
	ormInstance.RegisterModel(&TestModel{})
	return ormInstance
}

func TestORM_Connect_Close(t *testing.T) {
	ormInstance := setupORM()
	config := orm.ConnectionConfig{}
	err := ormInstance.Connect(config)
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	if !ormInstance.IsConnected() {
		t.Error("Should be connected after Connect")
	}
	err = ormInstance.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}
	if ormInstance.IsConnected() {
		t.Error("Should not be connected after Close")
	}
}

func TestORM_RegisterModel_GetMetadata(t *testing.T) {
	ormInstance := setupORM()
	model := &TestModel{}
	err := ormInstance.RegisterModel(model)
	if err != nil {
		t.Errorf("RegisterModel failed: %v", err)
	}
	meta, err := ormInstance.GetMetadata(model)
	if err != nil {
		t.Errorf("GetMetadata failed: %v", err)
	}
	if meta == nil {
		t.Error("Metadata should not be nil")
	}
}

func TestORM_Repository_Query_Raw(t *testing.T) {
	ormInstance := setupORM()
	model := &TestModel{}
	repo := ormInstance.Repository(model)
	if repo == nil {
		t.Error("Repository should not be nil")
	}
	query := ormInstance.Query(model)
	if query == nil {
		t.Error("Query should not be nil")
	}
	raw := ormInstance.Raw("SELECT 1")
	if raw == nil {
		t.Error("Raw should not be nil")
	}
}

func TestORM_Transaction(t *testing.T) {
	ormInstance := setupORM()

	// Connect first
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

	err = ormInstance.Transaction(func(tx orm.ORM) error {
		if tx == nil {
			t.Error("Transaction ORM should not be nil")
		}
		return nil
	})
	if err != nil {
		t.Errorf("Transaction failed: %v", err)
	}
}

func TestORM_TransactionWithContext(t *testing.T) {
	ormInstance := setupORM()
	ctx := context.Background()

	// Connect first
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

	err = ormInstance.TransactionWithContext(ctx, func(tx orm.ORM) error {
		if tx == nil {
			t.Error("TransactionWithContext ORM should not be nil")
		}
		return nil
	})
	if err != nil {
		t.Errorf("TransactionWithContext failed: %v", err)
	}
}

func TestORM_CreateTable_DropTable_Migrate(t *testing.T) {
	ormInstance := setupORM()
	model := &TestModel{}
	err := ormInstance.CreateTable(model)
	if err != nil {
		t.Errorf("CreateTable failed: %v", err)
	}
	err = ormInstance.DropTable(model)
	if err != nil {
		t.Errorf("DropTable failed: %v", err)
	}
	err = ormInstance.Migrate()
	if err != nil {
		t.Errorf("Migrate failed: %v", err)
	}
}

func TestORM_ErrorCases(t *testing.T) {
	// Simulate nil dialect
	ormInstance := orm.New(nil)
	err := ormInstance.Connect(orm.ConnectionConfig{})
	if err == nil {
		t.Error("Expected error when dialect is nil")
	}
}
