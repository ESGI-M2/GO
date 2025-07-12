package metadata

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/ESGI-M2/GO/orm/core/interfaces"
)

// Tag constants for better maintainability
const (
	TagORM        = "orm"
	TagDB         = "db"
	TagPrimary    = "primary"
	TagAutoIncr   = "autoincrement"
	TagUnique     = "unique"
	TagIndex      = "index"
	TagForeignKey = "foreign"
	TagRelation   = "relation"
	TagColumn     = "column"
	TagLength     = "length"
	TagDefault    = "default"
	TagNullable   = "nullable"
)

// ORMTag represents parsed ORM tag data
type ORMTag struct {
	Column       string
	PrimaryKey   bool
	AutoIncr     bool
	Unique       bool
	Index        bool
	ForeignKey   string
	Length       int
	Default      string
	Nullable     bool
	Relation     string
	RelationType string
	SoftDelete   bool
}

// extractColumn extracts column information from a struct field
func (mm *Manager) extractColumn(field reflect.StructField) (*interfaces.Column, error) {
	// Try ORM tag first, then fall back to DB tag for backward compatibility
	ormTagStr := field.Tag.Get(TagORM)
	dbTag := field.Tag.Get(TagDB)

	if ormTagStr == "" && dbTag == "" {
		return nil, nil
	}

	if ormTagStr == "-" || dbTag == "-" {
		return nil, nil
	}

	column := &interfaces.Column{
		Type:     getSQLType(field.Type),
		Nullable: true, // Default to nullable
	}

	// Parse ORM tag if present
	if ormTagStr != "" {
		ormTag := parseORMTag(ormTagStr)

		// If this tag defines a relation only, skip column generation
		if ormTag.Relation != "" {
			return nil, nil
		}

		// Set column name
		if ormTag.Column != "" {
			column.Name = ormTag.Column
		} else {
			column.Name = strings.ToLower(field.Name)
		}

		// Set boolean flags
		column.PrimaryKey = ormTag.PrimaryKey
		column.AutoIncrement = ormTag.AutoIncr
		column.Unique = ormTag.Unique
		column.Index = ormTag.Index
		column.Nullable = ormTag.Nullable

		// Set soft delete flag
		column.SoftDelete = ormTag.SoftDelete
		if column.SoftDelete {
			column.Nullable = true // soft delete columns must allow NULL
		}

		// Set length
		if ormTag.Length > 0 {
			column.Length = ormTag.Length
		}

		// Set default value
		if ormTag.Default != "" {
			column.Default = parseDefaultValue(ormTag.Default, field.Type)
		}

		// Set foreign key
		if ormTag.ForeignKey != "" {
			parts := strings.Split(ormTag.ForeignKey, ".")
			if len(parts) == 2 {
				column.ForeignKey = &interfaces.ForeignKey{
					ReferencedTable:  parts[0],
					ReferencedColumn: parts[1],
				}
			}
		}
	} else {
		// Fall back to old DB tag parsing for backward compatibility
		column.Name = dbTag

		// Extract primary key
		if field.Tag.Get(TagPrimary) == "true" {
			column.PrimaryKey = true
			column.Nullable = false
		}

		// Extract auto increment
		if field.Tag.Get(TagAutoIncr) == "true" {
			column.AutoIncrement = true
		}

		// Extract unique constraint
		if field.Tag.Get(TagUnique) == "true" {
			column.Unique = true
		}

		// Extract index
		if field.Tag.Get(TagIndex) == "true" {
			column.Index = true
		}

		// Extract length for string types
		if length := field.Tag.Get(TagLength); length != "" {
			if l, err := parseInt(length); err == nil {
				column.Length = l
			}
		}

		// Extract default value
		if defaultValue := field.Tag.Get(TagDefault); defaultValue != "" {
			column.Default = parseDefaultValue(defaultValue, field.Type)
		}

		// Extract foreign key
		if fk := field.Tag.Get(TagForeignKey); fk != "" {
			parts := strings.Split(fk, ".")
			if len(parts) == 2 {
				column.ForeignKey = &interfaces.ForeignKey{
					ReferencedTable:  parts[0],
					ReferencedColumn: parts[1],
					OnDelete:         field.Tag.Get("ondelete"),
					OnUpdate:         field.Tag.Get("onupdate"),
				}
			}
		}
	}

	return column, nil
}

// extractRelations extracts relationship information from struct fields
func (mm *Manager) extractRelations(t reflect.Type, metadata *interfaces.ModelMetadata) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Check for relation tags (both old and new format)
		relationType := field.Tag.Get("relation")
		ormTagStr := field.Tag.Get(TagORM)

		if relationType == "" && ormTagStr != "" {
			// Parse ORM tag to check for relation
			ormTag := parseORMTag(ormTagStr)
			if ormTag.Relation != "" {
				relationType = ormTag.Relation
			}
		}

		if relationType != "" {
			relation := &interfaces.Relation{
				Type:        parseRelationType(relationType),
				TargetModel: field.Type,
				Lazy:        field.Tag.Get("lazy") == "true",
			}

			// Extract foreign key information from ORM tag
			if ormTagStr != "" {
				ormTag := parseORMTag(ormTagStr)
				if ormTag.ForeignKey != "" {
					relation.ForeignKey = ormTag.ForeignKey
				}
			} else {
				// Fall back to old tag format
				if fk := field.Tag.Get("foreign_key"); fk != "" {
					relation.ForeignKey = fk
				}
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
func (mm *Manager) extractIndexes(t reflect.Type, metadata *interfaces.ModelMetadata) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Check for index tags (both old and new format)
		indexName := field.Tag.Get("index_name")
		ormTagStr := field.Tag.Get(TagORM)

		// Check if field has index in ORM tag
		if ormTagStr != "" {
			ormTag := parseORMTag(ormTagStr)
			if ormTag.Index {
				indexName = "idx_" + strings.ToLower(field.Name)
			}
		}

		if indexName != "" {
			columnName := field.Tag.Get(TagDB)
			if columnName == "" {
				columnName = strings.ToLower(field.Name)
			}

			index := interfaces.Index{
				Name:    indexName,
				Columns: []string{columnName},
				Unique:  field.Tag.Get(TagUnique) == "true",
			}
			metadata.Indexes = append(metadata.Indexes, index)
		}
	}
}

// parseORMTag parses ORM tags like "primary,auto" or "column:title,index"
func parseORMTag(tag string) *ORMTag {
	if tag == "" {
		return &ORMTag{}
	}

	ormTag := &ORMTag{}
	parts := strings.Split(tag, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Handle key-value pairs like "column:title"
		if strings.Contains(part, ":") {
			kv := strings.SplitN(part, ":", 2)
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])

			switch key {
			case "column":
				ormTag.Column = value
			case "fk", "foreign_key":
				ormTag.ForeignKey = value
			case "length":
				if len, err := strconv.Atoi(value); err == nil {
					ormTag.Length = len
				}
			case "default":
				ormTag.Default = value
			case "relation":
				ormTag.Relation = value
			case "type":
				ormTag.RelationType = value
			}
		} else {
			// Handle boolean flags like "primary", "auto"
			switch part {
			case "pk", "primary":
				ormTag.PrimaryKey = true
			case "auto", "auto_increment":
				ormTag.AutoIncr = true
			case "unique":
				ormTag.Unique = true
			case "index":
				ormTag.Index = true
			case "nullable":
				ormTag.Nullable = true
			case "soft":
				ormTag.SoftDelete = true
			}
		}
	}

	return ormTag
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
		// Check for time.Time specifically
		if t.String() == "time.Time" {
			return "TIMESTAMP"
		}
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

// parseRelationType parses relation type from string
func parseRelationType(relationType string) interfaces.RelationType {
	switch strings.ToLower(relationType) {
	case "one_to_one":
		return interfaces.OneToOne
	case "one_to_many":
		return interfaces.OneToMany
	case "many_to_one":
		return interfaces.ManyToOne
	case "many_to_many":
		return interfaces.ManyToMany
	default:
		return interfaces.OneToOne
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
