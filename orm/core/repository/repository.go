package repository

import (
	"fmt"
	"reflect"

	"github.com/ESGI-M2/GO/orm/core/interfaces"
)

// RepositoryImpl implements the Repository interface
type RepositoryImpl struct {
	orm      interfaces.ORM
	metadata *interfaces.ModelMetadata
	err      error
}

// NewRepository creates a new repository instance
func NewRepository(orm interfaces.ORM, metadata *interfaces.ModelMetadata) *RepositoryImpl {
	return &RepositoryImpl{
		orm:      orm,
		metadata: metadata,
	}
}

// NewErrorRepository creates a repository that will return an error
func NewErrorRepository(orm interfaces.ORM, err error) *RepositoryImpl {
	return &RepositoryImpl{
		orm: orm,
		err: err,
	}
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
