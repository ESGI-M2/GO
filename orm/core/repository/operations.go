package repository

import (
	"fmt"
	"reflect"
	"strings"
)

// Save saves an entity (insert or update)
func (r *RepositoryImpl) Save(entity interface{}) error {
	if r.err != nil {
		return r.err
	}

	if r.metadata == nil {
		return fmt.Errorf("metadata not available")
	}

	// Check if entity has an ID to determine if it's an insert or update
	entityValue := reflect.ValueOf(entity)
	if entityValue.Kind() == reflect.Ptr {
		entityValue = entityValue.Elem()
	}

	// Find field by name (case-insensitive)
	var idField reflect.Value
	for i := 0; i < entityValue.NumField(); i++ {
		field := entityValue.Type().Field(i)
		if strings.EqualFold(field.Name, r.metadata.PrimaryKey) {
			idField = entityValue.Field(i)
			break
		}
	}
	if !idField.IsValid() {
		return fmt.Errorf("primary key field %s not found", r.metadata.PrimaryKey)
	}

	// If ID is zero value, it's an insert
	if isZeroValue(idField) {
		return r.insert(entity)
	}

	// Otherwise, it's an update
	return r.update(entity)
}

// Update updates an entity
func (r *RepositoryImpl) Update(entity interface{}) error {
	if r.err != nil {
		return r.err
	}

	return r.update(entity)
}

// Delete deletes an entity
func (r *RepositoryImpl) Delete(entity interface{}) error {
	if r.err != nil {
		return r.err
	}

	if r.metadata == nil {
		return fmt.Errorf("metadata not available")
	}

	entityValue := reflect.ValueOf(entity)
	if entityValue.Kind() == reflect.Ptr {
		entityValue = entityValue.Elem()
	}

	// Find field by name (case-insensitive)
	var idField reflect.Value
	for i := 0; i < entityValue.NumField(); i++ {
		field := entityValue.Type().Field(i)
		if strings.EqualFold(field.Name, r.metadata.PrimaryKey) {
			idField = entityValue.Field(i)
			break
		}
	}
	if !idField.IsValid() {
		return fmt.Errorf("primary key field %s not found", r.metadata.PrimaryKey)
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?",
		r.metadata.TableName, r.metadata.PrimaryKey)

	_, err := r.orm.GetDialect().Exec(query, idField.Interface())
	if err != nil {
		return fmt.Errorf("failed to delete entity: %w", err)
	}

	return nil
}

// DeleteBy deletes records by criteria
func (r *RepositoryImpl) DeleteBy(criteria map[string]interface{}) error {
	if r.err != nil {
		return r.err
	}

	if r.metadata == nil {
		return fmt.Errorf("metadata not available")
	}

	var conditions []string
	var args []interface{}

	for field, value := range criteria {
		conditions = append(conditions, fmt.Sprintf("%s = ?", field))
		args = append(args, value)
	}

	whereClause := strings.Join(conditions, " AND ")
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", r.metadata.TableName, whereClause)

	_, err := r.orm.GetDialect().Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete records by criteria: %w", err)
	}

	return nil
}

// insert inserts a new entity
func (r *RepositoryImpl) insert(entity interface{}) error {
	entityValue := reflect.ValueOf(entity)
	if entityValue.Kind() == reflect.Ptr {
		entityValue = entityValue.Elem()
	}

	var columns []string
	var values []interface{}
	var placeholders []string

	for _, column := range r.metadata.Columns {
		// Find field by name (case-insensitive)
		var field reflect.Value
		for i := 0; i < entityValue.NumField(); i++ {
			fieldType := entityValue.Type().Field(i)
			if strings.EqualFold(fieldType.Name, column.Name) {
				field = entityValue.Field(i)
				break
			}
		}
		if !field.IsValid() {
			continue
		}

		// Skip auto-increment fields for insert
		if column.AutoIncrement {
			continue
		}

		columns = append(columns, column.Name)
		values = append(values, field.Interface())
		placeholders = append(placeholders, "?")
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		r.metadata.TableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	_, err := r.orm.GetDialect().Exec(query, values...)
	if err != nil {
		return fmt.Errorf("failed to insert entity: %w", err)
	}

	return nil
}

// update updates an existing entity
func (r *RepositoryImpl) update(entity interface{}) error {
	entityValue := reflect.ValueOf(entity)
	if entityValue.Kind() == reflect.Ptr {
		entityValue = entityValue.Elem()
	}

	var sets []string
	var values []interface{}

	for _, column := range r.metadata.Columns {
		// Find field by name (case-insensitive)
		var field reflect.Value
		for i := 0; i < entityValue.NumField(); i++ {
			fieldType := entityValue.Type().Field(i)
			if strings.EqualFold(fieldType.Name, column.Name) {
				field = entityValue.Field(i)
				break
			}
		}
		if !field.IsValid() || (field.Kind() == reflect.Ptr && field.IsNil()) {
			continue // skip unset or nil fields
		}

		// Skip primary key for update
		if column.Name == r.metadata.PrimaryKey {
			continue
		}

		sets = append(sets, fmt.Sprintf("%s = ?", column.Name))
		values = append(values, field.Interface())
	}

	// Add WHERE condition for primary key
	// Find field by name (case-insensitive)
	var idField reflect.Value
	for i := 0; i < entityValue.NumField(); i++ {
		field := entityValue.Type().Field(i)
		if strings.EqualFold(field.Name, r.metadata.PrimaryKey) {
			idField = entityValue.Field(i)
			break
		}
	}
	if !idField.IsValid() {
		return fmt.Errorf("primary key field %s not found", r.metadata.PrimaryKey)
	}
	values = append(values, idField.Interface())

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?",
		r.metadata.TableName,
		strings.Join(sets, ", "),
		r.metadata.PrimaryKey)

	_, err := r.orm.GetDialect().Exec(query, values...)
	if err != nil {
		return fmt.Errorf("failed to update entity: %w", err)
	}

	return nil
}

// mapToStruct maps a database result to a struct
func (r *RepositoryImpl) mapToStruct(result map[string]interface{}) (interface{}, error) {
	if r.metadata == nil {
		return nil, fmt.Errorf("metadata not available")
	}

	// Create a new instance of the struct
	entity := reflect.New(r.metadata.Type).Interface()
	entityValue := reflect.ValueOf(entity).Elem()

	// Map database columns to struct fields
	for _, column := range r.metadata.Columns {
		// Find field by name (case-insensitive)
		var field reflect.Value
		for i := 0; i < entityValue.NumField(); i++ {
			fieldType := entityValue.Type().Field(i)
			if strings.EqualFold(fieldType.Name, column.Name) {
				field = entityValue.Field(i)
				break
			}
		}
		if !field.IsValid() {
			continue
		}

		value, exists := result[column.Name]
		if !exists {
			continue
		}

		if err := setFieldValue(field, value); err != nil {
			return nil, fmt.Errorf("failed to set field %s: %w", column.Name, err)
		}
	}

	return entity, nil
}

// isZeroValue checks if a reflect.Value is a zero value
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	default:
		return false
	}
}

// setFieldValue sets a field value with proper type conversion
func setFieldValue(field reflect.Value, value interface{}) error {
	if value == nil {
		field.Set(reflect.Zero(field.Type()))
		return nil
	}

	valueType := reflect.TypeOf(value)
	fieldType := field.Type()

	// If types match, set directly
	if valueType == fieldType {
		field.Set(reflect.ValueOf(value))
		return nil
	}

	// Handle type conversions
	switch fieldType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v := value.(type) {
		case int64:
			field.SetInt(v)
		case int:
			field.SetInt(int64(v))
		case float64:
			field.SetInt(int64(v))
		default:
			return fmt.Errorf("cannot convert %v to %s", value, fieldType)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch v := value.(type) {
		case int64:
			field.SetUint(uint64(v))
		case int:
			field.SetUint(uint64(v))
		case float64:
			field.SetUint(uint64(v))
		default:
			return fmt.Errorf("cannot convert %v to %s", value, fieldType)
		}
	case reflect.Float32, reflect.Float64:
		switch v := value.(type) {
		case float64:
			field.SetFloat(v)
		case int64:
			field.SetFloat(float64(v))
		case int:
			field.SetFloat(float64(v))
		default:
			return fmt.Errorf("cannot convert %v to %s", value, fieldType)
		}
	case reflect.String:
		if str, ok := value.(string); ok {
			field.SetString(str)
		} else {
			return fmt.Errorf("cannot convert %v to string", value)
		}
	case reflect.Bool:
		if b, ok := value.(bool); ok {
			field.SetBool(b)
		} else {
			return fmt.Errorf("cannot convert %v to bool", value)
		}
	default:
		return fmt.Errorf("unsupported field type: %s", fieldType)
	}

	return nil
}
