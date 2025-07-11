package unit

import (
	"testing"

	"github.com/ESGI-M2/GO/orm"
	"github.com/ESGI-M2/GO/orm/dialect"
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

func setupMetadataTest() orm.ORM {
	mockDialect := &dialect.MockDialect{}
	ormInstance := orm.New(mockDialect)
	return ormInstance
}

func TestMetadata_ExtractMetadata(t *testing.T) {
	ormInstance := setupMetadataTest()
	model := &MetadataTestModel{}

	err := ormInstance.RegisterModel(model)
	if err != nil {
		t.Errorf("RegisterModel failed: %v", err)
	}

	metadata, err := ormInstance.GetMetadata(model)
	if err != nil {
		t.Errorf("GetMetadata failed: %v", err)
	}

	if metadata == nil {
		t.Fatal("Metadata should not be nil")
	}

	if metadata.TableName != "metadatatestmodel" {
		t.Errorf("Expected table name 'metadatatestmodel', got '%s'", metadata.TableName)
	}

	if len(metadata.Columns) == 0 {
		t.Error("Metadata should have columns")
	}
}

func TestMetadata_PrimaryKey_AutoIncrement(t *testing.T) {
	ormInstance := setupMetadataTest()
	model := &MetadataTestModel{}

	err := ormInstance.RegisterModel(model)
	if err != nil {
		t.Errorf("RegisterModel failed: %v", err)
	}

	metadata, err := ormInstance.GetMetadata(model)
	if err != nil {
		t.Errorf("GetMetadata failed: %v", err)
	}

	if metadata.PrimaryKey != "id" {
		t.Errorf("Expected primary key 'id', got '%s'", metadata.PrimaryKey)
	}

	if metadata.AutoIncrement != "id" {
		t.Errorf("Expected auto increment 'id', got '%s'", metadata.AutoIncrement)
	}
}

func TestMetadata_ColumnTypes(t *testing.T) {
	ormInstance := setupMetadataTest()
	model := &MetadataTestModel{}

	err := ormInstance.RegisterModel(model)
	if err != nil {
		t.Errorf("RegisterModel failed: %v", err)
	}

	metadata, err := ormInstance.GetMetadata(model)
	if err != nil {
		t.Errorf("GetMetadata failed: %v", err)
	}

	// Check that we have the expected number of columns
	expectedColumns := 5 // ID, Name, Email, Age, IsActive
	if len(metadata.Columns) != expectedColumns {
		t.Errorf("Expected %d columns, got %d", expectedColumns, len(metadata.Columns))
	}

	// Check for specific columns
	hasName := false
	hasEmail := false
	hasAge := false
	hasIsActive := false

	for _, col := range metadata.Columns {
		switch col.Name {
		case "name":
			hasName = true
		case "email":
			hasEmail = true
			if !col.Unique {
				t.Error("Email column should be unique")
			}
		case "age":
			hasAge = true
			if !col.Index {
				t.Error("Age column should be indexed")
			}
		case "is_active":
			hasIsActive = true
		}
	}

	if !hasName {
		t.Error("Should have name column")
	}
	if !hasEmail {
		t.Error("Should have email column")
	}
	if !hasAge {
		t.Error("Should have age column")
	}
	if !hasIsActive {
		t.Error("Should have is_active column")
	}
}

func TestMetadata_ForeignKey(t *testing.T) {
	ormInstance := setupMetadataTest()
	model := &MetadataTestModelWithRelations{}

	err := ormInstance.RegisterModel(model)
	if err != nil {
		t.Errorf("RegisterModel failed: %v", err)
	}

	metadata, err := ormInstance.GetMetadata(model)
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
	ormInstance := setupMetadataTest()
	model := &MetadataTestModel{}

	// First call should extract metadata
	metadata1, err := ormInstance.GetMetadata(model)
	if err != nil {
		t.Errorf("First GetMetadata failed: %v", err)
	}

	// Second call should use cache
	metadata2, err := ormInstance.GetMetadata(model)
	if err != nil {
		t.Errorf("Second GetMetadata failed: %v", err)
	}

	// Both should be the same instance (cached)
	if metadata1 != metadata2 {
		t.Error("Metadata should be cached and return same instance")
	}
}

func TestMetadata_ErrorCases(t *testing.T) {
	ormInstance := setupMetadataTest()

	// Test with non-struct type
	var nonStruct int = 42
	_, err := ormInstance.GetMetadata(nonStruct)
	if err == nil {
		t.Error("Expected error when model is not a struct")
	}
}

func TestMetadata_TagParsing(t *testing.T) {
	ormInstance := setupMetadataTest()
	model := &MetadataTestModel{}

	err := ormInstance.RegisterModel(model)
	if err != nil {
		t.Errorf("RegisterModel failed: %v", err)
	}

	metadata, err := ormInstance.GetMetadata(model)
	if err != nil {
		t.Errorf("GetMetadata failed: %v", err)
	}

	// Check that tags are parsed correctly
	for _, col := range metadata.Columns {
		switch col.Name {
		case "id":
			if !col.PrimaryKey {
				t.Error("ID column should be primary key")
			}
			if !col.AutoIncrement {
				t.Error("ID column should be auto increment")
			}
		case "email":
			if !col.Unique {
				t.Error("Email column should be unique")
			}
		case "age":
			if !col.Index {
				t.Error("Age column should be indexed")
			}
		}
	}
}
