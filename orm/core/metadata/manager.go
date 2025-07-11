package metadata

import (
	"fmt"
	"reflect"

	"project/orm/core/interfaces"
)

// Manager handles model metadata extraction and caching
type Manager struct {
	metadata map[reflect.Type]*interfaces.ModelMetadata
}

// NewManager creates a new metadata manager
func NewManager() *Manager {
	return &Manager{
		metadata: make(map[reflect.Type]*interfaces.ModelMetadata),
	}
}

// ExtractMetadata extracts metadata from a model struct
func (mm *Manager) ExtractMetadata(model interface{}) (*interfaces.ModelMetadata, error) {
	t := reflect.TypeOf(model)

	// Handle pointer types
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Check if it's a struct type
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("model must be a struct type, got %s", t.Kind())
	}

	// Check if metadata is already cached
	if cached, exists := mm.metadata[t]; exists {
		return cached, nil
	}

	metadata := &interfaces.ModelMetadata{
		Type:      t,
		TableName: getTableName(t),
		Columns:   make([]interfaces.Column, 0),
		Relations: make(map[string]*interfaces.Relation),
		Indexes:   make([]interfaces.Index, 0),
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

// GetMetadata returns cached metadata for a model
func (mm *Manager) GetMetadata(model interface{}) (*interfaces.ModelMetadata, error) {
	return mm.ExtractMetadata(model)
}

// ClearCache clears the metadata cache
func (mm *Manager) ClearCache() {
	mm.metadata = make(map[reflect.Type]*interfaces.ModelMetadata)
}
