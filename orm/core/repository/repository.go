package repository

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/ESGI-M2/GO/orm/core/interfaces"
)

// RepositoryImpl implements the Repository interface
type RepositoryImpl struct {
	orm      interfaces.ORM
	metadata *interfaces.ModelMetadata
	model    interface{}
}

// NewRepository creates a new repository instance
func NewRepository(orm interfaces.ORM, metadata *interfaces.ModelMetadata, model interface{}) *RepositoryImpl {
	return &RepositoryImpl{
		orm:      orm,
		metadata: metadata,
		model:    model,
	}
}

// Find finds a record by ID
func (r *RepositoryImpl) Find(id interface{}) (interface{}, error) {
	query := r.orm.Query(r.model).Where("id", "=", id)
	result, err := query.FindOne()
	if err != nil {
		return nil, fmt.Errorf("failed to find record: %w", err)
	}
	return result, nil
}

// FindWithRelations finds a record by ID with relations
func (r *RepositoryImpl) FindWithRelations(id interface{}, relations ...string) (interface{}, error) {
	query := r.orm.Query(r.model).Where("id", "=", id)

	// Add relations
	for _, relation := range relations {
		query = query.With(relation, func(q interfaces.QueryBuilder) interfaces.QueryBuilder {
			return q
		})
	}

	result, err := query.FindOne()
	if err != nil {
		return nil, fmt.Errorf("failed to find record with relations: %w", err)
	}
	return result, nil
}

// FindAll finds all records
func (r *RepositoryImpl) FindAll() ([]interface{}, error) {
	query := r.orm.Query(r.model)
	results, err := query.Find()
	if err != nil {
		return nil, fmt.Errorf("failed to find all records: %w", err)
	}
	return r.convertToInterface(results), nil
}

// FindAllWithRelations finds all records with relations
func (r *RepositoryImpl) FindAllWithRelations(relations ...string) ([]interface{}, error) {
	query := r.orm.Query(r.model)

	// Add relations
	for _, relation := range relations {
		query = query.With(relation, func(q interfaces.QueryBuilder) interfaces.QueryBuilder {
			return q
		})
	}

	results, err := query.Find()
	if err != nil {
		return nil, fmt.Errorf("failed to find all records with relations: %w", err)
	}
	return r.convertToInterface(results), nil
}

// FindBy finds records by criteria
func (r *RepositoryImpl) FindBy(criteria map[string]interface{}) ([]interface{}, error) {
	query := r.orm.Query(r.model)

	for field, value := range criteria {
		query = query.Where(field, "=", value)
	}

	results, err := query.Find()
	if err != nil {
		return nil, fmt.Errorf("failed to find records by criteria: %w", err)
	}
	return r.convertToInterface(results), nil
}

// FindByWithRelations finds records by criteria with relations
func (r *RepositoryImpl) FindByWithRelations(criteria map[string]interface{}, relations ...string) ([]interface{}, error) {
	query := r.orm.Query(r.model)

	for field, value := range criteria {
		query = query.Where(field, "=", value)
	}

	// Add relations
	for _, relation := range relations {
		query = query.With(relation, func(q interfaces.QueryBuilder) interfaces.QueryBuilder {
			return q
		})
	}

	results, err := query.Find()
	if err != nil {
		return nil, fmt.Errorf("failed to find records by criteria with relations: %w", err)
	}
	return r.convertToInterface(results), nil
}

// FindOneBy finds one record by criteria
func (r *RepositoryImpl) FindOneBy(criteria map[string]interface{}) (interface{}, error) {
	query := r.orm.Query(r.model)

	for field, value := range criteria {
		query = query.Where(field, "=", value)
	}

	result, err := query.FindOne()
	if err != nil {
		return nil, fmt.Errorf("failed to find one record by criteria: %w", err)
	}
	return result, nil
}

// Create creates a new record
func (r *RepositoryImpl) Create(entity interface{}) error {
	// Execute before hooks
	if err := r.executeHooks("BeforeCreate", entity); err != nil {
		return err
	}

	// Set timestamps if enabled
	if r.metadata.Timestamps {
		r.setTimestamps(entity, true)
	}

	// Execute before save hooks
	if err := r.executeHooks("BeforeSave", entity); err != nil {
		return err
	}

	// Create the record
	if err := r.orm.Repository(r.model).Save(entity); err != nil {
		return fmt.Errorf("failed to create record: %w", err)
	}

	// Execute after hooks
	if err := r.executeHooks("AfterCreate", entity); err != nil {
		return err
	}

	// Execute after save hooks
	if err := r.executeHooks("AfterSave", entity); err != nil {
		return err
	}

	return nil
}

// SoftDeleteBy soft deletes records by criteria
func (r *RepositoryImpl) SoftDeleteBy(criteria map[string]interface{}) error {
	if !r.metadata.SoftDeletes {
		return fmt.Errorf("soft deletes not enabled for this model")
	}

	query := r.orm.Query(r.model)

	for field, value := range criteria {
		query = query.Where(field, "=", value)
	}

	results, err := query.Find()
	if err != nil {
		return fmt.Errorf("failed to find records for soft deletion: %w", err)
	}

	for _, result := range results {
		if err := r.SoftDelete(result); err != nil {
			return err
		}
	}

	return nil
}

// BatchCreate creates multiple records in batch
func (r *RepositoryImpl) BatchCreate(entities []interface{}) error {
	if len(entities) == 0 {
		return nil
	}

	// Execute before hooks for each entity
	for _, entity := range entities {
		if err := r.executeHooks("BeforeCreate", entity); err != nil {
			return err
		}

		if r.metadata.Timestamps {
			r.setTimestamps(entity, true)
		}

		if err := r.executeHooks("BeforeSave", entity); err != nil {
			return err
		}
	}

	// Batch create
	for _, entity := range entities {
		if err := r.orm.Repository(r.model).Save(entity); err != nil {
			return fmt.Errorf("failed to batch create record: %w", err)
		}
	}

	// Execute after hooks for each entity
	for _, entity := range entities {
		if err := r.executeHooks("AfterCreate", entity); err != nil {
			return err
		}

		if err := r.executeHooks("AfterSave", entity); err != nil {
			return err
		}
	}

	return nil
}

// BatchUpdate updates multiple records in batch
func (r *RepositoryImpl) BatchUpdate(entities []interface{}) error {
	if len(entities) == 0 {
		return nil
	}

	// Execute before hooks for each entity
	for _, entity := range entities {
		if err := r.executeHooks("BeforeUpdate", entity); err != nil {
			return err
		}

		if r.metadata.Timestamps {
			r.setTimestamps(entity, false)
		}

		if err := r.executeHooks("BeforeSave", entity); err != nil {
			return err
		}
	}

	// Batch update
	for _, entity := range entities {
		if err := r.orm.Repository(r.model).Update(entity); err != nil {
			return fmt.Errorf("failed to batch update record: %w", err)
		}
	}

	// Execute after hooks for each entity
	for _, entity := range entities {
		if err := r.executeHooks("AfterUpdate", entity); err != nil {
			return err
		}

		if err := r.executeHooks("AfterSave", entity); err != nil {
			return err
		}
	}

	return nil
}

// BatchDelete deletes multiple records in batch
func (r *RepositoryImpl) BatchDelete(entities []interface{}) error {
	if len(entities) == 0 {
		return nil
	}

	// Execute before hooks for each entity
	for _, entity := range entities {
		if err := r.executeHooks("BeforeDelete", entity); err != nil {
			return err
		}
	}

	// Batch delete
	for _, entity := range entities {
		if err := r.Delete(entity); err != nil {
			return err
		}
	}

	// Execute after hooks for each entity
	for _, entity := range entities {
		if err := r.executeHooks("AfterDelete", entity); err != nil {
			return err
		}
	}

	return nil
}

// SoftDelete soft deletes a record
func (r *RepositoryImpl) SoftDelete(entity interface{}) error {
	if !r.metadata.SoftDeletes {
		return fmt.Errorf("soft deletes not enabled for this model")
	}

	// Set deleted_at timestamp
	r.setDeletedAt(entity)

	// Update the record
	return r.Update(entity)
}

// Restore restores a soft-deleted record
func (r *RepositoryImpl) Restore(entity interface{}) error {
	if !r.metadata.SoftDeletes {
		return fmt.Errorf("soft deletes not enabled for this model")
	}

	// Clear deleted_at timestamp
	r.clearDeletedAt(entity)

	// Update the record
	return r.Update(entity)
}

// ForceDelete force deletes a record (ignores soft deletes)
func (r *RepositoryImpl) ForceDelete(entity interface{}) error {
	// Execute before hooks
	if err := r.executeHooks("BeforeDelete", entity); err != nil {
		return err
	}

	// Hard delete
	if err := r.orm.Repository(r.model).Delete(entity); err != nil {
		return fmt.Errorf("failed to force delete record: %w", err)
	}

	// Execute after hooks
	if err := r.executeHooks("AfterDelete", entity); err != nil {
		return err
	}

	return nil
}

// FindTrashed finds soft-deleted records
func (r *RepositoryImpl) FindTrashed() ([]interface{}, error) {
	if !r.metadata.SoftDeletes {
		return nil, fmt.Errorf("soft deletes not enabled for this model")
	}

	query := r.orm.Query(r.model).WhereNotNull(r.metadata.DeletedAt)
	results, err := query.Find()
	if err != nil {
		return nil, fmt.Errorf("failed to find trashed records: %w", err)
	}
	return r.convertToInterface(results), nil
}

// RestoreBy restores soft-deleted records by criteria
func (r *RepositoryImpl) RestoreBy(criteria map[string]interface{}) error {
	if !r.metadata.SoftDeletes {
		return fmt.Errorf("soft deletes not enabled for this model")
	}

	query := r.orm.Query(r.model).WhereNotNull(r.metadata.DeletedAt)

	for field, value := range criteria {
		query = query.Where(field, "=", value)
	}

	results, err := query.Find()
	if err != nil {
		return fmt.Errorf("failed to find trashed records: %w", err)
	}

	for _, result := range results {
		if err := r.Restore(result); err != nil {
			return err
		}
	}

	return nil
}

// Scope applies a named scope
func (r *RepositoryImpl) Scope(name string, args ...interface{}) interfaces.Repository {
	if _, exists := r.metadata.Scopes[name]; exists {
		// Create a new repository with the scope applied
		// This is a simplified implementation
		return r
	}
	return r
}

// Chunk processes records in chunks
func (r *RepositoryImpl) Chunk(size int, fn func([]interface{}) error) error {
	offset := 0

	for {
		query := r.orm.Query(r.model).Limit(size).Offset(offset)
		results, err := query.Find()
		if err != nil {
			return fmt.Errorf("failed to get chunk: %w", err)
		}

		if len(results) == 0 {
			break
		}

		chunk := r.convertToInterface(results)
		if err := fn(chunk); err != nil {
			return err
		}

		if len(results) < size {
			break
		}

		offset += size
	}

	return nil
}

// Each processes records one by one
func (r *RepositoryImpl) Each(fn func(interface{}) error) error {
	return r.Chunk(1, func(chunk []interface{}) error {
		if len(chunk) > 0 {
			return fn(chunk[0])
		}
		return nil
	})
}

// Pluck gets a single column's value from the first result
func (r *RepositoryImpl) Pluck(field string) ([]interface{}, error) {
	query := r.orm.Query(r.model).Select(field)
	results, err := query.Find()
	if err != nil {
		return nil, fmt.Errorf("failed to pluck field: %w", err)
	}

	var values []interface{}
	for _, result := range results {
		if value, exists := result[field]; exists {
			values = append(values, value)
		}
	}

	return values, nil
}

// Value gets a single value from the first result
func (r *RepositoryImpl) Value(field string) (interface{}, error) {
	query := r.orm.Query(r.model).Select(field).Limit(1)
	result, err := query.FindOne()
	if err != nil {
		return nil, fmt.Errorf("failed to get value: %w", err)
	}

	if value, exists := result[field]; exists {
		return value, nil
	}

	return nil, nil
}

// Count counts all records
func (r *RepositoryImpl) Count() (int64, error) {
	return r.orm.Query(r.model).Count()
}

// Exists checks if a record exists
func (r *RepositoryImpl) Exists(id interface{}) (bool, error) {
	return r.orm.Query(r.model).Where("id", "=", id).Exists()
}

// Increment increments a field value for all records
func (r *RepositoryImpl) Increment(field string, amount interface{}) error {
	if r.metadata == nil {
		return fmt.Errorf("metadata not available")
	}
	query := fmt.Sprintf("UPDATE %s SET %s = %s + %s WHERE 1=1", r.metadata.TableName, field, field, r.orm.GetDialect().GetPlaceholder(0))
	_, err := r.orm.GetDialect().Exec(query, amount)
	if err != nil {
		return fmt.Errorf("failed to increment field: %w", err)
	}
	return nil
}

// Decrement decrements a field value for all records
func (r *RepositoryImpl) Decrement(field string, amount interface{}) error {
	if r.metadata == nil {
		return fmt.Errorf("metadata not available")
	}
	query := fmt.Sprintf("UPDATE %s SET %s = %s - %s WHERE 1=1", r.metadata.TableName, field, field, r.orm.GetDialect().GetPlaceholder(0))
	_, err := r.orm.GetDialect().Exec(query, amount)
	if err != nil {
		return fmt.Errorf("failed to decrement field: %w", err)
	}
	return nil
}

// Helper methods

// executeHooks executes model hooks
func (r *RepositoryImpl) executeHooks(hookType string, entity interface{}) error {
	if r.metadata.Hooks == nil {
		return nil
	}

	var hooks []func(interface{}) error

	switch hookType {
	case "BeforeCreate":
		hooks = r.metadata.Hooks.BeforeCreate
	case "AfterCreate":
		hooks = r.metadata.Hooks.AfterCreate
	case "BeforeUpdate":
		hooks = r.metadata.Hooks.BeforeUpdate
	case "AfterUpdate":
		hooks = r.metadata.Hooks.AfterUpdate
	case "BeforeDelete":
		hooks = r.metadata.Hooks.BeforeDelete
	case "AfterDelete":
		hooks = r.metadata.Hooks.AfterDelete
	case "BeforeSave":
		hooks = r.metadata.Hooks.BeforeSave
	case "AfterSave":
		hooks = r.metadata.Hooks.AfterSave
	}

	for _, hook := range hooks {
		if err := hook(entity); err != nil {
			return fmt.Errorf("hook %s failed: %w", hookType, err)
		}
	}

	return nil
}

// setTimestamps sets created_at and updated_at timestamps
func (r *RepositoryImpl) setTimestamps(entity interface{}, isCreate bool) {
	if !r.metadata.Timestamps {
		return
	}

	entityValue := reflect.ValueOf(entity)
	if entityValue.Kind() == reflect.Ptr {
		entityValue = entityValue.Elem()
	}

	now := time.Now()

	if isCreate && r.metadata.CreatedAt != "" {
		if field := entityValue.FieldByName(r.metadata.CreatedAt); field.IsValid() && field.CanSet() {
			field.Set(reflect.ValueOf(now))
		}
	}

	if r.metadata.UpdatedAt != "" {
		if field := entityValue.FieldByName(r.metadata.UpdatedAt); field.IsValid() && field.CanSet() {
			field.Set(reflect.ValueOf(now))
		}
	}
}

// setDeletedAt sets the deleted_at timestamp
func (r *RepositoryImpl) setDeletedAt(entity interface{}) {
	if !r.metadata.SoftDeletes {
		return
	}

	entityVal := reflect.ValueOf(entity)
	if entityVal.Kind() == reflect.Ptr {
		entityVal = entityVal.Elem()
	}

	entityType := entityVal.Type()

	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		tag := field.Tag.Get("orm")
		if strings.Contains(tag, "soft") {
			fieldVal := entityVal.Field(i)
			if fieldVal.IsValid() && fieldVal.CanSet() {
				now := time.Now()
				if fieldVal.Kind() == reflect.Ptr {
					fieldVal.Set(reflect.ValueOf(&now))
				} else {
					fieldVal.Set(reflect.ValueOf(now))
				}
			}
			return
		}
	}
}

// clearDeletedAt clears the deleted_at timestamp
func (r *RepositoryImpl) clearDeletedAt(entity interface{}) {
	if !r.metadata.SoftDeletes {
		return
	}

	entityVal := reflect.ValueOf(entity)
	if entityVal.Kind() == reflect.Ptr {
		entityVal = entityVal.Elem()
	}

	entityType := entityVal.Type()

	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		tag := field.Tag.Get("orm")
		if strings.Contains(tag, "soft") {
			fieldVal := entityVal.Field(i)
			if fieldVal.IsValid() && fieldVal.CanSet() {
				fieldVal.Set(reflect.Zero(fieldVal.Type()))
			}
			return
		}
	}
}

// convertToInterface converts []map[string]interface{} to []interface{}
func (r *RepositoryImpl) convertToInterface(data []map[string]interface{}) []interface{} {
	result := make([]interface{}, len(data))
	for i, item := range data {
		result[i] = item
	}
	return result
}
