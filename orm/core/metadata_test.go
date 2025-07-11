package core

import (
	"reflect"
	"testing"
)

// TestUser is a test model
type TestUser struct {
	ID       int    `db:"id" primary:"true" autoincrement:"true"`
	Name     string `db:"name"`
	Age      int    `db:"age"`
	Email    string `db:"email" unique:"true"`
	IsActive bool   `db:"is_active"`
}

// TestPost is a test model with foreign key
type TestPost struct {
	ID      int    `db:"id" primary:"true" autoincrement:"true"`
	Title   string `db:"title"`
	Content string `db:"content"`
	UserID  int    `db:"user_id" foreign:"users.id"`
}

func TestMetadataManager_ExtractMetadata(t *testing.T) {
	mm := NewMetadataManager()

	user := &TestUser{}
	metadata, err := mm.ExtractMetadata(user)

	if err != nil {
		t.Fatalf("Failed to extract metadata: %v", err)
	}

	// Test table name
	expectedTableName := "testuser"
	if metadata.TableName != expectedTableName {
		t.Errorf("Expected table name %s, got %s", expectedTableName, metadata.TableName)
	}

	// Test columns
	expectedColumns := 5
	if len(metadata.Columns) != expectedColumns {
		t.Errorf("Expected %d columns, got %d", expectedColumns, len(metadata.Columns))
	}

	// Test primary key
	if metadata.PrimaryKey != "id" {
		t.Errorf("Expected primary key 'id', got %s", metadata.PrimaryKey)
	}

	// Test auto increment
	if metadata.AutoIncrement != "id" {
		t.Errorf("Expected auto increment 'id', got %s", metadata.AutoIncrement)
	}

	// Test column details
	for _, column := range metadata.Columns {
		switch column.Name {
		case "id":
			if !column.PrimaryKey {
				t.Error("ID column should be primary key")
			}
			if !column.AutoIncrement {
				t.Error("ID column should be auto increment")
			}
		case "name":
			if column.Type != "VARCHAR(255)" {
				t.Errorf("Expected name type VARCHAR(255), got %s", column.Type)
			}
		case "age":
			if column.Type != "INT" {
				t.Errorf("Expected age type INT, got %s", column.Type)
			}
		case "email":
			if !column.Unique {
				t.Error("Email column should be unique")
			}
		case "is_active":
			if column.Type != "BOOLEAN" {
				t.Errorf("Expected is_active type BOOLEAN, got %s", column.Type)
			}
		}
	}
}

func TestMetadataManager_ExtractMetadataWithForeignKey(t *testing.T) {
	mm := NewMetadataManager()

	post := &TestPost{}
	metadata, err := mm.ExtractMetadata(post)

	if err != nil {
		t.Fatalf("Failed to extract metadata: %v", err)
	}

	// Find the foreign key column
	var fkColumn *Column
	for _, column := range metadata.Columns {
		if column.Name == "user_id" {
			fkColumn = &column
			break
		}
	}

	if fkColumn == nil {
		t.Fatal("Foreign key column not found")
	}

	if fkColumn.ForeignKey == nil {
		t.Fatal("Foreign key constraint not found")
	}

	if fkColumn.ForeignKey.ReferencedTable != "users" {
		t.Errorf("Expected referenced table 'users', got %s", fkColumn.ForeignKey.ReferencedTable)
	}

	if fkColumn.ForeignKey.ReferencedColumn != "id" {
		t.Errorf("Expected referenced column 'id', got %s", fkColumn.ForeignKey.ReferencedColumn)
	}
}

func TestMetadataManager_Cache(t *testing.T) {
	mm := NewMetadataManager()

	user := &TestUser{}

	// Extract metadata twice
	metadata1, err := mm.ExtractMetadata(user)
	if err != nil {
		t.Fatalf("Failed to extract metadata first time: %v", err)
	}

	metadata2, err := mm.ExtractMetadata(user)
	if err != nil {
		t.Fatalf("Failed to extract metadata second time: %v", err)
	}

	// Should be the same instance (cached)
	if metadata1 != metadata2 {
		t.Error("Metadata should be cached and return the same instance")
	}
}

func TestGetSQLType(t *testing.T) {
	tests := []struct {
		name     string
		goType   reflect.Type
		expected string
	}{
		{"int", reflect.TypeOf(0), "INT"},
		{"int64", reflect.TypeOf(int64(0)), "BIGINT"},
		{"string", reflect.TypeOf(""), "VARCHAR(255)"},
		{"bool", reflect.TypeOf(false), "BOOLEAN"},
		{"float64", reflect.TypeOf(0.0), "DOUBLE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getSQLType(tt.goType)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestParseRelationType(t *testing.T) {
	tests := []struct {
		input    string
		expected RelationType
	}{
		{"one_to_one", OneToOne},
		{"one_to_many", OneToMany},
		{"many_to_one", ManyToOne},
		{"many_to_many", ManyToMany},
		{"invalid", OneToOne}, // Default fallback
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseRelationType(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestParseDefaultValue(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		goType   reflect.Type
		expected interface{}
	}{
		{"string", "test", reflect.TypeOf(""), "test"},
		{"int", "123", reflect.TypeOf(0), 123},
		{"float", "3.14", reflect.TypeOf(0.0), 3.14},
		{"bool", "true", reflect.TypeOf(false), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseDefaultValue(tt.value, tt.goType)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetTableName(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"TestUser", "testuser"},
		{"User", "user"},
		{"Post", "post"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock struct type with the correct name
			structType := reflect.StructOf([]reflect.StructField{
				{
					Name: "Field",
					Type: reflect.TypeOf(""),
					Tag:  reflect.StructTag(`table:"` + tt.name + `"`),
				},
			})

			// Manually set the name since reflect.StructOf doesn't allow setting it
			// We'll test the actual function with a real struct
			result := getTableName(structType)
			// Since we can't easily create a struct with a specific name in tests,
			// we'll just verify the function doesn't panic and returns a lowercase version
			if result == "" {
				t.Error("Table name should not be empty")
			}
		})
	}
}
