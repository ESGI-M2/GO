package core

import (
	"fmt"
	"reflect"
	"strings"
)

// MetadataManager handles model metadata extraction and caching
type MetadataManager struct {
	metadata map[reflect.Type]*ModelMetadata
}

// NewMetadataManager creates a new metadata manager
func NewMetadataManager() *MetadataManager {
	return &MetadataManager{
		metadata: make(map[reflect.Type]*ModelMetadata),
	}
}

// ExtractMetadata extracts metadata from a model struct
func (mm *MetadataManager) ExtractMetadata(model interface{}) (*ModelMetadata, error) {
	t := reflect.TypeOf(model)

	// Handle pointer types
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Check if metadata is already cached
	if cached, exists := mm.metadata[t]; exists {
		return cached, nil
	}

	metadata := &ModelMetadata{
		Type:      t,
		TableName: getTableName(t),
		Columns:   make([]Column, 0),
		Relations: make(map[string]*Relation),
		Indexes:   make([]Index, 0),
	}

	// Extract columns from struct fields
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		column, err := mm.extractColumn(field)
		if err != nil {
			return nil, fmt.Errorf("error extracting column from field %s: %w", field.Name, err)
		}

		if column != nil {
			metadata.Columns = append(metadata.Columns, *column)

			// Track primary key and auto increment
			if column.PrimaryKey {
				metadata.PrimaryKey = column.Name
			}
			if column.AutoIncrement {
				metadata.AutoIncrement = column.Name
			}
		}
	}

	// Extract relations
	mm.extractRelations(t, metadata)

	// Extract indexes
	mm.extractIndexes(t, metadata)

	// Cache the metadata
	mm.metadata[t] = metadata

	return metadata, nil
}

// extractColumn extracts column information from a struct field
func (mm *MetadataManager) extractColumn(field reflect.StructField) (*Column, error) {
	dbTag := field.Tag.Get("db")
	if dbTag == "" || dbTag == "-" {
		return nil, nil
	}

	column := &Column{
		Name:     dbTag,
		Type:     getSQLType(field.Type),
		Nullable: true, // Default to nullable
	}

	// Extract primary key
	if field.Tag.Get("primary") == "true" {
		column.PrimaryKey = true
		column.Nullable = false
	}

	// Extract auto increment
	if field.Tag.Get("autoincrement") == "true" {
		column.AutoIncrement = true
	}

	// Extract unique constraint
	if field.Tag.Get("unique") == "true" {
		column.Unique = true
	}

	// Extract index
	if field.Tag.Get("index") == "true" {
		column.Index = true
	}

	// Extract length for string types
	if length := field.Tag.Get("length"); length != "" {
		if l, err := parseInt(length); err == nil {
			column.Length = l
		}
	}

	// Extract default value
	if defaultValue := field.Tag.Get("default"); defaultValue != "" {
		column.Default = parseDefaultValue(defaultValue, field.Type)
	}

	// Extract foreign key
	if fk := field.Tag.Get("foreign"); fk != "" {
		parts := strings.Split(fk, ".")
		if len(parts) == 2 {
			column.ForeignKey = &ForeignKey{
				ReferencedTable:  parts[0],
				ReferencedColumn: parts[1],
				OnDelete:         field.Tag.Get("ondelete"),
				OnUpdate:         field.Tag.Get("onupdate"),
			}
		}
	}

	return column, nil
}

// extractRelations extracts relationship information from struct fields
func (mm *MetadataManager) extractRelations(t reflect.Type, metadata *ModelMetadata) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Check for relation tags
		if relationType := field.Tag.Get("relation"); relationType != "" {
			relation := &Relation{
				Type:        parseRelationType(relationType),
				TargetModel: field.Type,
				Lazy:        field.Tag.Get("lazy") == "true",
			}

			// Extract foreign key information
			if fk := field.Tag.Get("foreign_key"); fk != "" {
				relation.ForeignKey = fk
			}
			if rk := field.Tag.Get("referenced_key"); rk != "" {
				relation.ReferencedKey = rk
			}
			if jt := field.Tag.Get("join_table"); jt != "" {
				relation.JoinTable = jt
			}

			metadata.Relations[field.Name] = relation
		}
	}
}

// extractIndexes extracts index information from struct tags
func (mm *MetadataManager) extractIndexes(t reflect.Type, metadata *ModelMetadata) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Check for index tags
		if indexName := field.Tag.Get("index_name"); indexName != "" {
			index := Index{
				Name:    indexName,
				Columns: []string{field.Tag.Get("db")},
				Unique:  field.Tag.Get("unique") == "true",
			}
			metadata.Indexes = append(metadata.Indexes, index)
		}
	}
}

// getSQLType maps Go types to SQL types
func getSQLType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return "INT"
	case reflect.Int64:
		return "BIGINT"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return "INT UNSIGNED"
	case reflect.Uint64:
		return "BIGINT UNSIGNED"
	case reflect.Float32:
		return "FLOAT"
	case reflect.Float64:
		return "DOUBLE"
	case reflect.String:
		return "VARCHAR(255)"
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.Struct:
		if t == reflect.TypeOf([]byte{}) {
			return "BLOB"
		}
		return "TEXT"
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			return "BLOB"
		}
		return "TEXT"
	default:
		return "TEXT"
	}
}

// getTableName extracts the table name from a struct
func getTableName(t reflect.Type) string {
	// Check for table tag on the struct
	if t.NumField() > 0 {
		if tableTag := t.Field(0).Tag.Get("table"); tableTag != "" {
			return tableTag
		}
	}

	// Default to lowercase struct name
	return strings.ToLower(t.Name())
}

// parseRelationType parses relation type from string
func parseRelationType(relationType string) RelationType {
	switch strings.ToLower(relationType) {
	case "one_to_one":
		return OneToOne
	case "one_to_many":
		return OneToMany
	case "many_to_one":
		return ManyToOne
	case "many_to_many":
		return ManyToMany
	default:
		return OneToOne
	}
}

// parseDefaultValue parses default value from string
func parseDefaultValue(value string, t reflect.Type) interface{} {
	switch t.Kind() {
	case reflect.String:
		return value
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if i, err := parseInt(value); err == nil {
			return i
		}
	case reflect.Float32, reflect.Float64:
		if f, err := parseFloat(value); err == nil {
			return f
		}
	case reflect.Bool:
		return strings.ToLower(value) == "true"
	}
	return value
}

// parseInt parses an integer from string
func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

// parseFloat parses a float from string
func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

// GetMetadata returns cached metadata for a model
func (mm *MetadataManager) GetMetadata(model interface{}) (*ModelMetadata, error) {
	return mm.ExtractMetadata(model)
}

// ClearCache clears the metadata cache
func (mm *MetadataManager) ClearCache() {
	mm.metadata = make(map[reflect.Type]*ModelMetadata)
}
