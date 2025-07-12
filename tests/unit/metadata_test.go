package unit

import (
	"testing"

	"github.com/ESGI-M2/GO/orm/builder"
	"github.com/ESGI-M2/GO/orm/factory"
)

type MetadataTestModel struct {
	ID       int    `orm:"pk,auto"`
	Name     string `orm:"column:name"`
	Email    string `orm:"column:email,unique"`
	Age      int    `orm:"column:age,index"`
	IsActive bool   `orm:"column:is_active,default:true"`
}

type MetadataTestModelWithRelations struct {
	ID       int    `orm:"pk,auto"`
	Name     string `orm:"column:name"`
	UserID   int    `orm:"column:user_id,fk:users.id"`
	Category string `orm:"column:category,index"`
}

func setupMetadataTest() *builder.SimpleORM {
	return builder.NewSimpleORM().WithDialect(factory.Mock)
}

func TestMetadata_ExtractMetadata(t *testing.T) {
	orm := setupMetadataTest()
	model := &MetadataTestModel{}

	// Register model and connect
	orm.RegisterModel(model)
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Get underlying ORM to access metadata
	underlyingORM := orm.GetORM()
	metadata, err := underlyingORM.GetMetadata(model)
	if err != nil {
		t.Errorf("GetMetadata failed: %v", err)
	}

	if metadata == nil {
		t.Fatal("Metadata should not be nil")
	}

	// Check that metadata has expected properties
	if metadata.TableName == "" {
		t.Error("Metadata should have a table name")
	}

	if len(metadata.Columns) == 0 {
		t.Error("Metadata should have columns")
	}
}

func TestMetadata_PrimaryKey_AutoIncrement(t *testing.T) {
	orm := setupMetadataTest()
	model := &MetadataTestModel{}

	// Register model and connect
	orm.RegisterModel(model)
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Get underlying ORM to access metadata
	underlyingORM := orm.GetORM()
	metadata, err := underlyingORM.GetMetadata(model)
	if err != nil {
		t.Errorf("GetMetadata failed: %v", err)
	}

	if metadata.PrimaryKey == "" {
		t.Error("Expected primary key to be set")
	}

	if metadata.AutoIncrement == "" {
		t.Error("Expected auto increment to be set")
	}
}

func TestMetadata_ColumnTypes(t *testing.T) {
	orm := setupMetadataTest()
	model := &MetadataTestModel{}

	// Register model and connect
	orm.RegisterModel(model)
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Get underlying ORM to access metadata
	underlyingORM := orm.GetORM()
	metadata, err := underlyingORM.GetMetadata(model)
	if err != nil {
		t.Errorf("GetMetadata failed: %v", err)
	}

	// Check that columns were extracted
	if len(metadata.Columns) == 0 {
		t.Error("Expected columns to be extracted")
	}

	// Check for specific column types (SQL types, not Go types)
	hasStringColumn := false
	hasIntColumn := false
	hasBoolColumn := false

	for _, col := range metadata.Columns {
		if col.Type == "VARCHAR(255)" {
			hasStringColumn = true
		}
		if col.Type == "INT" {
			hasIntColumn = true
		}
		if col.Type == "BOOLEAN" {
			hasBoolColumn = true
		}
	}

	if !hasStringColumn {
		t.Error("Expected to find VARCHAR(255) column")
	}
	if !hasIntColumn {
		t.Error("Expected to find INT column")
	}
	if !hasBoolColumn {
		t.Error("Expected to find BOOLEAN column")
	}
}

func TestMetadata_ForeignKey(t *testing.T) {
	orm := setupMetadataTest()
	model := &MetadataTestModelWithRelations{}

	// Register model and connect
	orm.RegisterModel(model)
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Get underlying ORM to access metadata
	underlyingORM := orm.GetORM()
	metadata, err := underlyingORM.GetMetadata(model)
	if err != nil {
		t.Errorf("GetMetadata failed: %v", err)
	}

	// Check for foreign key
	hasForeignKey := false
	for _, col := range metadata.Columns {
		if col.Name == "user_id" && col.ForeignKey != nil {
			hasForeignKey = true
			if col.ForeignKey.ReferencedTable != "users" {
				t.Errorf("Expected referenced table 'users', got '%s'", col.ForeignKey.ReferencedTable)
			}
			if col.ForeignKey.ReferencedColumn != "id" {
				t.Errorf("Expected referenced column 'id', got '%s'", col.ForeignKey.ReferencedColumn)
			}
			break
		}
	}

	if !hasForeignKey {
		t.Error("Should have foreign key column")
	}
}

func TestMetadata_Cache(t *testing.T) {
	orm := setupMetadataTest()
	model := &MetadataTestModel{}

	// Register model and connect
	orm.RegisterModel(model)
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Get underlying ORM to access metadata
	underlyingORM := orm.GetORM()

	// First call should extract metadata
	metadata1, err := underlyingORM.GetMetadata(model)
	if err != nil {
		t.Errorf("First GetMetadata failed: %v", err)
	}

	// Second call should use cache
	metadata2, err := underlyingORM.GetMetadata(model)
	if err != nil {
		t.Errorf("Second GetMetadata failed: %v", err)
	}

	// Both should be the same instance (cached)
	if metadata1 != metadata2 {
		t.Error("Metadata should be cached and return same instance")
	}
}

func TestMetadata_ErrorCases(t *testing.T) {
	orm := setupMetadataTest()

	// Connect without registering model
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Get underlying ORM to access metadata
	underlyingORM := orm.GetORM()

	// Test with non-struct type
	var nonStruct int = 42
	_, err = underlyingORM.GetMetadata(nonStruct)
	if err == nil {
		t.Error("Expected error when model is not a struct")
	}
}

func TestMetadata_TagParsing(t *testing.T) {
	orm := setupMetadataTest()
	model := &MetadataTestModel{}

	// Register model and connect
	orm.RegisterModel(model)
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Get underlying ORM to access metadata
	underlyingORM := orm.GetORM()
	metadata, err := underlyingORM.GetMetadata(model)
	if err != nil {
		t.Errorf("GetMetadata failed: %v", err)
	}

	// Check that ORM tags were parsed correctly
	for _, col := range metadata.Columns {
		if col.Name == "name" && col.Type != "VARCHAR(255)" {
			t.Errorf("Expected column type 'VARCHAR(255)', got '%s'", col.Type)
		}
		if col.Name == "email" && col.Type != "VARCHAR(255)" {
			t.Errorf("Expected column type 'VARCHAR(255)', got '%s'", col.Type)
		}
		if col.Name == "id" && !col.PrimaryKey {
			t.Error("Expected ID column to be primary key")
		}
	}
}
