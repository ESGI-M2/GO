package unit

import (
	"context"
	"testing"

	"github.com/ESGI-M2/GO/orm/builder"
	"github.com/ESGI-M2/GO/orm/core/interfaces"
	"github.com/ESGI-M2/GO/orm/factory"
)

type TestModel struct {
	ID   int    `orm:"pk,auto"`
	Name string `orm:"column:name"`
}

func setupORM() *builder.SimpleORM {
	return builder.NewSimpleORM().
		WithDialect(factory.Mock).
		RegisterModel(&TestModel{})
}

func TestORM_Connect_Close(t *testing.T) {
	orm := setupORM()
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	if !orm.IsConnected() {
		t.Error("Should be connected after Connect")
	}
	err = orm.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}
	if orm.IsConnected() {
		t.Error("Should not be connected after Close")
	}
}

func TestORM_RegisterModel_GetMetadata(t *testing.T) {
	orm := setupORM()
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Get underlying ORM to access metadata
	underlyingORM := orm.GetORM()
	if underlyingORM == nil {
		t.Fatal("Underlying ORM should not be nil")
	}

	model := &TestModel{}
	meta, err := underlyingORM.GetMetadata(model)
	if err != nil {
		t.Errorf("GetMetadata failed: %v", err)
	}
	if meta == nil {
		t.Error("Metadata should not be nil")
	}
}

func TestORM_Repository_Query_Raw(t *testing.T) {
	orm := setupORM()
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	model := &TestModel{}
	repo := orm.Repository(model)
	if repo == nil {
		t.Error("Repository should not be nil")
	}
	query := orm.Query(model)
	if query == nil {
		t.Error("Query should not be nil")
	}
	raw := orm.Raw("SELECT 1")
	if raw == nil {
		t.Error("Raw should not be nil")
	}
}

func TestORM_Transaction(t *testing.T) {
	orm := setupORM()

	// Connect first
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	err = orm.Transaction(func(tx interfaces.ORM) error {
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
	orm := setupORM()
	ctx := context.Background()

	// Connect first
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Get underlying ORM for transaction
	underlyingORM := orm.GetORM()
	err = underlyingORM.TransactionWithContext(ctx, func(tx interfaces.ORM) error {
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
	orm := setupORM()
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Get underlying ORM for table operations
	underlyingORM := orm.GetORM()
	model := &TestModel{}

	err = underlyingORM.CreateTable(model)
	if err != nil {
		t.Errorf("CreateTable failed: %v", err)
	}
	err = underlyingORM.DropTable(model)
	if err != nil {
		t.Errorf("DropTable failed: %v", err)
	}
	err = underlyingORM.Migrate()
	if err != nil {
		t.Errorf("Migrate failed: %v", err)
	}
}

func TestORM_ErrorCases(t *testing.T) {
	// Test with invalid dialect
	orm := builder.NewSimpleORM().WithDialect("invalid")
	err := orm.Connect()
	if err == nil {
		t.Error("Expected error when dialect is invalid")
	}
}
