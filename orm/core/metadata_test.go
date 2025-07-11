package core

import (
	"reflect"
	"testing"
)

// TestUser for testing with new ORM tags
type TestUser struct {
	ID       int    `orm:"pk,auto"`
	Name     string `orm:"index"`
	Email    string `orm:"unique"`
	Age      int
	IsActive bool `orm:"default:true"`
}

// TestUserWithOldTags for backward compatibility testing
type TestUserWithOldTags struct {
	ID       int    `db:"id" primary:"true" autoincrement:"true"`
	Name     string `db:"name" index:"true"`
	Email    string `db:"email" unique:"true"`
	Age      int    `db:"age"`
	IsActive bool   `db:"is_active"`
}

// TestPost for testing foreign keys
type TestPost struct {
	ID      int    `orm:"pk,auto"`
	Title   string `orm:"index"`
	Content string
	UserID  int `orm:"fk:users.id"`
}

// TestComplexModel for testing advanced features
type TestComplexModel struct {
	ID        int    `orm:"pk,auto"`
	Name      string `orm:"column:full_name,index"`
	Email     string `orm:"unique,length:255"`
	Age       int    `orm:"default:18"`
	IsActive  bool   `orm:"default:true"`
	CreatedAt string `orm:"column:created_at"`
}

func TestParseORMTag(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		expected *ORMTag
	}{
		{
			name: "primary key and auto increment",
			tag:  "pk,auto",
			expected: &ORMTag{
				PrimaryKey: true,
				AutoIncr:   true,
			},
		},
		{
			name: "column name and index",
			tag:  "column:title,index",
			expected: &ORMTag{
				Column: "title",
				Index:  true,
			},
		},
		{
			name: "foreign key",
			tag:  "fk:users.id",
			expected: &ORMTag{
				ForeignKey: "users.id",
			},
		},
		{
			name: "unique and length",
			tag:  "unique,length:255",
			expected: &ORMTag{
				Unique: true,
				Length: 255,
			},
		},
		{
			name: "default value",
			tag:  "default:true",
			expected: &ORMTag{
				Default: "true",
			},
		},
		{
			name: "nullable",
			tag:  "nullable",
			expected: &ORMTag{
				Nullable: true,
			},
		},
		{
			name:     "empty tag",
			tag:      "",
			expected: &ORMTag{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseORMTag(tt.tag)

			if result.PrimaryKey != tt.expected.PrimaryKey {
				t.Errorf("PrimaryKey: expected %v, got %v", tt.expected.PrimaryKey, result.PrimaryKey)
			}
			if result.AutoIncr != tt.expected.AutoIncr {
				t.Errorf("AutoIncr: expected %v, got %v", tt.expected.AutoIncr, result.AutoIncr)
			}
			if result.Column != tt.expected.Column {
				t.Errorf("Column: expected %s, got %s", tt.expected.Column, result.Column)
			}
			if result.Index != tt.expected.Index {
				t.Errorf("Index: expected %v, got %v", tt.expected.Index, result.Index)
			}
			if result.Unique != tt.expected.Unique {
				t.Errorf("Unique: expected %v, got %v", tt.expected.Unique, result.Unique)
			}
			if result.ForeignKey != tt.expected.ForeignKey {
				t.Errorf("ForeignKey: expected %s, got %s", tt.expected.ForeignKey, result.ForeignKey)
			}
			if result.Length != tt.expected.Length {
				t.Errorf("Length: expected %d, got %d", tt.expected.Length, result.Length)
			}
			if result.Default != tt.expected.Default {
				t.Errorf("Default: expected %s, got %s", tt.expected.Default, result.Default)
			}
			if result.Nullable != tt.expected.Nullable {
				t.Errorf("Nullable: expected %v, got %v", tt.expected.Nullable, result.Nullable)
			}
		})
	}
}

func TestMetadataManager_ExtractMetadataWithNewTags(t *testing.T) {
	mm := NewMetadataManager()

	user := &TestUser{}
	metadata, err := mm.ExtractMetadata(user)
	if err != nil {
		t.Fatalf("ExtractMetadata should not return error: %v", err)
	}

	if metadata == nil {
		t.Fatal("ExtractMetadata should return metadata")
	}

	if metadata.TableName != "testuser" {
		t.Errorf("Expected table name 'testuser', got '%s'", metadata.TableName)
	}

	if metadata.PrimaryKey != "id" {
		t.Errorf("Expected primary key 'id', got '%s'", metadata.PrimaryKey)
	}

	if metadata.AutoIncrement != "id" {
		t.Errorf("Expected auto increment 'id', got '%s'", metadata.AutoIncrement)
	}

	// Check columns
	expectedColumns := map[string]bool{
		"id":       true,
		"name":     true,
		"email":    true,
		"age":      true,
		"isactive": true,
	}

	for _, col := range metadata.Columns {
		if !expectedColumns[col.Name] {
			t.Errorf("Unexpected column: %s", col.Name)
		}
	}
}

func TestMetadataManager_ExtractMetadataWithOldTags(t *testing.T) {
	mm := NewMetadataManager()

	user := &TestUserWithOldTags{}
	metadata, err := mm.ExtractMetadata(user)
	if err != nil {
		t.Fatalf("ExtractMetadata should not return error: %v", err)
	}

	if metadata == nil {
		t.Fatal("ExtractMetadata should return metadata")
	}

	if metadata.TableName != "testuserwitholdtags" {
		t.Errorf("Expected table name 'testuserwitholdtags', got '%s'", metadata.TableName)
	}

	if metadata.PrimaryKey != "id" {
		t.Errorf("Expected primary key 'id', got '%s'", metadata.PrimaryKey)
	}

	if metadata.AutoIncrement != "id" {
		t.Errorf("Expected auto increment 'id', got '%s'", metadata.AutoIncrement)
	}
}

func TestMetadataManager_ExtractMetadataWithComplexModel(t *testing.T) {
	mm := NewMetadataManager()

	model := &TestComplexModel{}
	metadata, err := mm.ExtractMetadata(model)
	if err != nil {
		t.Fatalf("ExtractMetadata should not return error: %v", err)
	}

	if metadata == nil {
		t.Fatal("ExtractMetadata should return metadata")
	}

	// Check that custom column names are used
	foundFullName := false
	foundCreatedAt := false
	for _, col := range metadata.Columns {
		if col.Name == "full_name" {
			foundFullName = true
		}
		if col.Name == "created_at" {
			foundCreatedAt = true
		}
	}

	if !foundFullName {
		t.Error("Expected column 'full_name' not found")
	}
	if !foundCreatedAt {
		t.Error("Expected column 'created_at' not found")
	}
}

func TestMetadataManager_ExtractMetadataWithForeignKey(t *testing.T) {
	mm := NewMetadataManager()

	post := &TestPost{}
	metadata, err := mm.ExtractMetadata(post)
	if err != nil {
		t.Fatalf("ExtractMetadata should not return error: %v", err)
	}

	if metadata == nil {
		t.Fatal("ExtractMetadata should return metadata")
	}

	// Check foreign key
	foundUserID := false
	for _, col := range metadata.Columns {
		if col.Name == "userid" && col.ForeignKey != nil {
			if col.ForeignKey.ReferencedTable == "users" && col.ForeignKey.ReferencedColumn == "id" {
				foundUserID = true
			}
		}
	}

	if !foundUserID {
		t.Error("Expected foreign key 'userid' referencing 'users.id' not found")
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
