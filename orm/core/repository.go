package core

import (
	"fmt"
	"reflect"
	"strings"
)

// RepositoryImpl implements the Repository interface
type RepositoryImpl struct {
	orm      *ORMImpl
	metadata *ModelMetadata
	err      error
}

// Find finds a record by ID
func (r *RepositoryImpl) Find(id interface{}) (interface{}, error) {
	if r.err != nil {
		return nil, r.err
	}

	if r.metadata == nil {
		return nil, fmt.Errorf("metadata not available")
	}

	query := r.orm.Query(reflect.New(r.metadata.Type).Interface())
	query = query.Where(r.metadata.PrimaryKey, "=", id)

	results, err := query.Find()
	if err != nil {
		return nil, fmt.Errorf("failed to find record: %w", err)
	}

	if len(results) == 0 {
		return nil, nil
	}

	return r.mapToStruct(results[0])
}

// FindAll finds all records
func (r *RepositoryImpl) FindAll() ([]interface{}, error) {
	if r.err != nil {
		return nil, r.err
	}

	if r.metadata == nil {
		return nil, fmt.Errorf("metadata not available")
	}

	query := r.orm.Query(reflect.New(r.metadata.Type).Interface())
	results, err := query.Find()
	if err != nil {
		return nil, fmt.Errorf("failed to find all records: %w", err)
	}

	var entities []interface{}
	for _, result := range results {
		entity, err := r.mapToStruct(result)
		if err != nil {
			return nil, fmt.Errorf("failed to map result to struct: %w", err)
		}
		entities = append(entities, entity)
	}

	return entities, nil
}

// FindBy finds records by criteria
func (r *RepositoryImpl) FindBy(criteria map[string]interface{}) ([]interface{}, error) {
	if r.err != nil {
		return nil, r.err
	}

	if r.metadata == nil {
		return nil, fmt.Errorf("metadata not available")
	}

	query := r.orm.Query(reflect.New(r.metadata.Type).Interface())

	for field, value := range criteria {
		query = query.Where(field, "=", value)
	}

	results, err := query.Find()
	if err != nil {
		return nil, fmt.Errorf("failed to find records by criteria: %w", err)
	}

	var entities []interface{}
	for _, result := range results {
		entity, err := r.mapToStruct(result)
		if err != nil {
			return nil, fmt.Errorf("failed to map result to struct: %w", err)
		}
		entities = append(entities, entity)
	}

	return entities, nil
}

// FindOneBy finds one record by criteria
func (r *RepositoryImpl) FindOneBy(criteria map[string]interface{}) (interface{}, error) {
	if r.err != nil {
		return nil, r.err
	}

	if r.metadata == nil {
		return nil, fmt.Errorf("metadata not available")
	}

	query := r.orm.Query(reflect.New(r.metadata.Type).Interface())

	for field, value := range criteria {
		query = query.Where(field, "=", value)
	}

	result, err := query.FindOne()
	if err != nil {
		return nil, fmt.Errorf("failed to find one record by criteria: %w", err)
	}

	if result == nil {
		return nil, nil
	}

	return r.mapToStruct(result)
}

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

	idField := entityValue.FieldByName(r.metadata.PrimaryKey)
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

	idField := entityValue.FieldByName(r.metadata.PrimaryKey)
	if !idField.IsValid() {
		return fmt.Errorf("primary key field %s not found", r.metadata.PrimaryKey)
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?",
		r.metadata.TableName, r.metadata.PrimaryKey)

	_, err := r.orm.dialect.Exec(query, idField.Interface())
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

	_, err := r.orm.dialect.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete records by criteria: %w", err)
	}

	return nil
}

// Count counts all records
func (r *RepositoryImpl) Count() (int64, error) {
	if r.err != nil {
		return 0, r.err
	}

	if r.metadata == nil {
		return 0, fmt.Errorf("metadata not available")
	}

	query := r.orm.Query(reflect.New(r.metadata.Type).Interface())
	return query.Count()
}

// Exists checks if a record exists by ID
func (r *RepositoryImpl) Exists(id interface{}) (bool, error) {
	if r.err != nil {
		return false, r.err
	}

	if r.metadata == nil {
		return false, fmt.Errorf("metadata not available")
	}

	query := r.orm.Query(reflect.New(r.metadata.Type).Interface())
	query = query.Where(r.metadata.PrimaryKey, "=", id)

	return query.Exists()
}

// insert inserts a new entity
func (r *RepositoryImpl) insert(entity interface{}) error {
	entityValue := reflect.ValueOf(entity)
	if entityValue.Kind() == reflect.Ptr {
		entityValue = entityValue.Elem()
	}

	var columns []string
	var placeholders []string
	var values []interface{}

	for _, column := range r.metadata.Columns {
		// Skip auto-increment fields
		if column.AutoIncrement {
			continue
		}

		field := entityValue.FieldByName(column.Name)
		if !field.IsValid() {
			continue
		}

		columns = append(columns, column.Name)
		placeholders = append(placeholders, "?")
		values = append(values, field.Interface())
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		r.metadata.TableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	result, err := r.orm.dialect.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("failed to insert entity: %w", err)
	}

	// Set the generated ID if auto-increment
	if r.metadata.AutoIncrement != "" {
		lastID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID: %w", err)
		}

		idField := entityValue.FieldByName(r.metadata.PrimaryKey)
		if idField.IsValid() && idField.CanSet() {
			idField.SetInt(lastID)
		}
	}

	return nil
}

// update updates an existing entity
func (r *RepositoryImpl) update(entity interface{}) error {
	if r.metadata == nil {
		return fmt.Errorf("metadata not available")
	}

	entityValue := reflect.ValueOf(entity)
	if entityValue.Kind() == reflect.Ptr {
		entityValue = entityValue.Elem()
	}

	var setClauses []string
	var values []interface{}

	for _, column := range r.metadata.Columns {
		// Skip primary key and auto-increment fields
		if column.PrimaryKey || column.AutoIncrement {
			continue
		}

		field := entityValue.FieldByName(column.Name)
		if !field.IsValid() {
			continue
		}

		setClauses = append(setClauses, fmt.Sprintf("%s = ?", column.Name))
		values = append(values, field.Interface())
	}

	// Add the primary key value for the WHERE clause
	idField := entityValue.FieldByName(r.metadata.PrimaryKey)
	if !idField.IsValid() {
		return fmt.Errorf("primary key field %s not found", r.metadata.PrimaryKey)
	}
	values = append(values, idField.Interface())

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?",
		r.metadata.TableName,
		strings.Join(setClauses, ", "),
		r.metadata.PrimaryKey)

	_, err := r.orm.dialect.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("failed to update entity: %w", err)
	}

	return nil
}

// mapToStruct maps a database result map to a struct
func (r *RepositoryImpl) mapToStruct(result map[string]interface{}) (interface{}, error) {
	if r.metadata == nil {
		return nil, fmt.Errorf("metadata not available")
	}

	// Create a new instance of the struct
	entity := reflect.New(r.metadata.Type).Interface()
	entityValue := reflect.ValueOf(entity).Elem()

	for _, column := range r.metadata.Columns {
		field := entityValue.FieldByName(column.Name)
		if !field.IsValid() || !field.CanSet() {
			continue
		}

		value, exists := result[column.Name]
		if !exists {
			continue
		}

		// Convert the value to the appropriate type
		if err := setFieldValue(field, value); err != nil {
			return nil, fmt.Errorf("failed to set field %s: %w", column.Name, err)
		}
	}

	return entity, nil
}

// isZeroValue checks if a value is the zero value for its type
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

// setFieldValue sets a field value with type conversion
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

	// Handle common type conversions
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
		switch v := value.(type) {
		case string:
			field.SetString(v)
		case []byte:
			field.SetString(string(v))
		default:
			return fmt.Errorf("cannot convert %v to %s", value, fieldType)
		}
	case reflect.Bool:
		switch v := value.(type) {
		case bool:
			field.SetBool(v)
		case int64:
			field.SetBool(v != 0)
		case int:
			field.SetBool(v != 0)
		default:
			return fmt.Errorf("cannot convert %v to %s", value, fieldType)
		}
	default:
		return fmt.Errorf("unsupported field type: %s", fieldType)
	}

	return nil
}
