package integration

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/ESGI-M2/GO/orm/core/interfaces"
)

// Add a global mock state for users
var mockUsers []map[string]interface{}

// Helper to reset mock state before each test
func resetMockUsers() {
	mockUsers = []map[string]interface{}{}
}

// setupTestDB creates a test database connection
func setupTestDB(t *testing.T) interfaces.ORM {
	// Create a new ORM instance with memory dialect for testing
	// For now, we'll use a mock implementation
	db := &MockORM{}

	return db
}

// MockORM is a mock implementation for testing
type MockORM struct{}

func (m *MockORM) Connect(config interfaces.ConnectionConfig) error {
	return nil
}

func (m *MockORM) Close() error {
	return nil
}

func (m *MockORM) IsConnected() bool {
	return true
}

func (m *MockORM) GetDialect() interfaces.Dialect {
	return &MockDialect{}
}

func (m *MockORM) RegisterModel(model interface{}) error {
	return nil
}

func (m *MockORM) GetMetadata(model interface{}) (*interfaces.ModelMetadata, error) {
	return &interfaces.ModelMetadata{
		TableName:   "users",
		PrimaryKey:  "id",
		Timestamps:  true,
		SoftDeletes: false,
	}, nil
}

func (m *MockORM) Query(model interface{}) interfaces.QueryBuilder {
	return &MockQueryBuilder{}
}

func (m *MockORM) Raw(sql string, args ...interface{}) interfaces.QueryBuilder {
	return &MockQueryBuilder{}
}

func (m *MockORM) Repository(model interface{}) interfaces.Repository {
	return &MockRepository{}
}

func (m *MockORM) Transaction(fn func(interfaces.ORM) error) error {
	return nil
}

func (m *MockORM) TransactionWithContext(ctx context.Context, fn func(interfaces.ORM) error) error {
	return nil
}

func (m *MockORM) CreateTable(model interface{}) error {
	return nil
}

func (m *MockORM) DropTable(model interface{}) error {
	return nil
}

func (m *MockORM) Migrate() error {
	return nil
}

func (m *MockORM) WithCache(ttl int) interfaces.ORM {
	return m
}

func (m *MockORM) WithConnectionPool(maxOpen, maxIdle int) interfaces.ORM {
	return m
}

func (m *MockORM) EnableQueryLog() interfaces.ORM {
	return m
}

func (m *MockORM) DisableQueryLog() interfaces.ORM {
	return m
}

// MockDialect is a mock dialect implementation
type MockDialect struct{}

func (m *MockDialect) Connect(config interfaces.ConnectionConfig) error {
	return nil
}

func (m *MockDialect) Close() error {
	return nil
}

func (m *MockDialect) Ping() error {
	return nil
}

func (m *MockDialect) Exec(query string, args ...interface{}) (sql.Result, error) {
	return &MockResult{}, nil
}

func (m *MockDialect) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (m *MockDialect) QueryRow(query string, args ...interface{}) *sql.Row {
	return nil
}

func (m *MockDialect) Begin() (interfaces.Transaction, error) {
	return &MockTransaction{}, nil
}

func (m *MockDialect) BeginTx(ctx context.Context, opts *sql.TxOptions) (interfaces.Transaction, error) {
	return &MockTransaction{}, nil
}

func (m *MockDialect) CreateTable(tableName string, columns []interfaces.Column) error {
	return nil
}

func (m *MockDialect) DropTable(tableName string) error {
	return nil
}

func (m *MockDialect) TableExists(tableName string) (bool, error) {
	return true, nil
}

func (m *MockDialect) GetSQLType(goType reflect.Type) string {
	return "TEXT"
}

func (m *MockDialect) GetPlaceholder(index int) string {
	return "?"
}

func (m *MockDialect) FullTextSearch(field, query string) string {
	return fmt.Sprintf("MATCH(%s) AGAINST('%s' IN BOOLEAN MODE)", field, query)
}

func (m *MockDialect) GetRandomFunction() string {
	return "RAND()"
}

func (m *MockDialect) GetDateFunction() string {
	return "NOW()"
}

func (m *MockDialect) GetJSONExtract() string {
	return "JSON_EXTRACT"
}

// MockTransaction is a mock transaction implementation
type MockTransaction struct{}

func (m *MockTransaction) Commit() error {
	return nil
}

func (m *MockTransaction) Rollback() error {
	return nil
}

func (m *MockTransaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	return &MockResult{}, nil
}

func (m *MockTransaction) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (m *MockTransaction) QueryRow(query string, args ...interface{}) *sql.Row {
	return nil
}

// MockResult is a mock sql.Result implementation
type MockResult struct{}

func (m *MockResult) LastInsertId() (int64, error) {
	return 1, nil
}

func (m *MockResult) RowsAffected() (int64, error) {
	return 1, nil
}

// MockQueryBuilder is a mock query builder implementation
type MockQueryBuilder struct{}

func (m *MockQueryBuilder) Select(fields ...string) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) From(table string) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) Where(field, operator string, value interface{}) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) WhereIn(field string, values []interface{}) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) WhereNotIn(field string, values []interface{}) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) OrderBy(field, direction string) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) GroupBy(fields ...string) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) Having(condition string, args ...interface{}) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) Limit(limit int) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) Offset(offset int) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) Join(table, condition string) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) LeftJoin(table, condition string) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) RightJoin(table, condition string) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) InnerJoin(table, condition string) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) Find() ([]map[string]interface{}, error) {
	return []map[string]interface{}{
		{"id": 1, "name": "John Doe", "age": 25, "email": "john@test.com"},
		{"id": 2, "name": "Jane Smith", "age": 30, "email": "jane@test.com"},
	}, nil
}

func (m *MockQueryBuilder) FindOne() (map[string]interface{}, error) {
	return map[string]interface{}{
		"id": 1, "name": "John Doe", "age": 25, "email": "john@test.com",
	}, nil
}

func (m *MockQueryBuilder) Count() (int64, error) {
	return 2, nil
}

func (m *MockQueryBuilder) Exists() (bool, error) {
	return true, nil
}

func (m *MockQueryBuilder) Raw(sql string, args ...interface{}) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) GetSQL() string {
	return "SELECT * FROM users"
}

func (m *MockQueryBuilder) GetArgs() []interface{} {
	return []interface{}{}
}

// Advanced query methods
func (m *MockQueryBuilder) WhereOr(conditions ...interfaces.WhereCondition) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) WhereRaw(condition string, args ...interface{}) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) WhereBetween(field string, min, max interface{}) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) WhereNotBetween(field string, min, max interface{}) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) WhereNull(field string) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) WhereNotNull(field string) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) WhereLike(field, pattern string) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) WhereNotLike(field, pattern string) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) WhereRegexp(field, pattern string) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) WhereNotRegexp(field, pattern string) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) FullTextSearch(fields []string, query string) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) SubQuery(alias string, fn func(interfaces.QueryBuilder) interfaces.QueryBuilder) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) With(relation string, fn func(interfaces.QueryBuilder) interfaces.QueryBuilder) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) WithCount(relation string) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) WithExists(relation string, fn func(interfaces.QueryBuilder) interfaces.QueryBuilder) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) CursorPaginate(cursorField string, cursorValue interface{}, limit int) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) OffsetPaginate(page, perPage int) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) ForUpdate() interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) ForShare() interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) Distinct() interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) Union(other interfaces.QueryBuilder) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) UnionAll(other interfaces.QueryBuilder) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) Lock(lockType string) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) Cache(ttl int) interfaces.QueryBuilder {
	return m
}

func (m *MockQueryBuilder) WithoutCache() interfaces.QueryBuilder {
	return m
}

// MockRepository is a mock repository implementation
type MockRepository struct{}

func (m *MockRepository) Find(id interface{}) (interface{}, error) {
	return map[string]interface{}{
		"id": 1, "name": "John Doe", "age": 25, "email": "john@test.com",
	}, nil
}

func (m *MockRepository) FindAll() ([]interface{}, error) {
	var result []interface{}
	for _, user := range mockUsers {
		if !user["deleted"].(bool) {
			result = append(result, user)
		}
	}
	return result, nil
}

func (m *MockRepository) FindBy(criteria map[string]interface{}) ([]interface{}, error) {
	return []interface{}{
		map[string]interface{}{"id": 1, "name": "John Doe", "age": 25, "email": "john@test.com"},
	}, nil
}

func (m *MockRepository) FindOneBy(criteria map[string]interface{}) (interface{}, error) {
	return map[string]interface{}{
		"id": 1, "name": "John Doe", "age": 25, "email": "john@test.com",
	}, nil
}

func (m *MockRepository) Save(entity interface{}) error {
	switch user := entity.(type) {
	case *User:
		mockUsers = append(mockUsers, map[string]interface{}{
			"id":      len(mockUsers) + 1,
			"name":    user.Name,
			"age":     user.Age,
			"email":   user.Email,
			"deleted": false,
		})
	case *UserWithSoftDelete:
		mockUsers = append(mockUsers, map[string]interface{}{
			"id":         len(mockUsers) + 1,
			"name":       user.Name,
			"age":        user.Age,
			"email":      user.Email,
			"deleted":    false,
			"deleted_at": user.DeletedAt,
		})
	}
	return nil
}

func (m *MockRepository) Update(entity interface{}) error {
	return nil
}

func (m *MockRepository) Delete(entity interface{}) error {
	return nil
}

func (m *MockRepository) DeleteBy(criteria map[string]interface{}) error {
	return nil
}

func (m *MockRepository) Count() (int64, error) {
	count := int64(0)
	for _, user := range mockUsers {
		if !user["deleted"].(bool) {
			count++
		}
	}
	return count, nil
}

func (m *MockRepository) Exists(id interface{}) (bool, error) {
	return true, nil
}

// Advanced repository methods
func (m *MockRepository) FindWithRelations(id interface{}, relations ...string) (interface{}, error) {
	return map[string]interface{}{
		"id": 1, "name": "John Doe", "age": 25, "email": "john@test.com",
		"posts": []map[string]interface{}{
			{"id": 1, "title": "First Post", "content": "Content 1"},
		},
	}, nil
}

func (m *MockRepository) FindAllWithRelations(relations ...string) ([]interface{}, error) {
	return []interface{}{
		map[string]interface{}{
			"id": 1, "name": "John Doe", "age": 25, "email": "john@test.com",
			"posts": []map[string]interface{}{
				{"id": 1, "title": "First Post", "content": "Content 1"},
			},
		},
	}, nil
}

func (m *MockRepository) FindByWithRelations(criteria map[string]interface{}, relations ...string) ([]interface{}, error) {
	return []interface{}{
		map[string]interface{}{
			"id": 1, "name": "John Doe", "age": 25, "email": "john@test.com",
		},
	}, nil
}

func (m *MockRepository) BatchCreate(entities []interface{}) error {
	for _, e := range entities {
		switch user := e.(type) {
		case *User:
			mockUsers = append(mockUsers, map[string]interface{}{
				"id":      len(mockUsers) + 1,
				"name":    user.Name,
				"age":     user.Age,
				"email":   user.Email,
				"deleted": false,
			})
		case *UserWithSoftDelete:
			mockUsers = append(mockUsers, map[string]interface{}{
				"id":         len(mockUsers) + 1,
				"name":       user.Name,
				"age":        user.Age,
				"email":      user.Email,
				"deleted":    false,
				"deleted_at": user.DeletedAt,
			})
		}
	}
	return nil
}

func (m *MockRepository) BatchUpdate(entities []interface{}) error {
	return nil
}

func (m *MockRepository) BatchDelete(entities []interface{}) error {
	for _, e := range entities {
		user := e.(*User)
		for _, u := range mockUsers {
			if u["email"] == user.Email {
				u["deleted"] = true
			}
		}
	}
	return nil
}

func (m *MockRepository) SoftDelete(entity interface{}) error {
	switch user := entity.(type) {
	case *User:
		for _, u := range mockUsers {
			if u["email"] == user.Email {
				u["deleted"] = true
			}
		}
	case *UserWithSoftDelete:
		for _, u := range mockUsers {
			if u["email"] == user.Email {
				u["deleted"] = true
			}
		}
	}
	return nil
}

func (m *MockRepository) Restore(entity interface{}) error {
	switch user := entity.(type) {
	case *User:
		for _, u := range mockUsers {
			if u["email"] == user.Email {
				u["deleted"] = false
			}
		}
	case *UserWithSoftDelete:
		for _, u := range mockUsers {
			if u["email"] == user.Email {
				u["deleted"] = false
			}
		}
	}
	return nil
}

func (m *MockRepository) ForceDelete(entity interface{}) error {
	return nil
}

func (m *MockRepository) FindTrashed() ([]interface{}, error) {
	var result []interface{}
	for _, user := range mockUsers {
		if user["deleted"].(bool) {
			result = append(result, user)
		}
	}
	return result, nil
}

func (m *MockRepository) RestoreBy(criteria map[string]interface{}) error {
	return nil
}

func (m *MockRepository) Scope(name string, args ...interface{}) interfaces.Repository {
	return m
}

func (m *MockRepository) Chunk(size int, fn func([]interface{}) error) error {
	chunk := []interface{}{
		map[string]interface{}{"id": 1, "name": "John Doe", "age": 25, "email": "john@test.com"},
		map[string]interface{}{"id": 2, "name": "Jane Smith", "age": 30, "email": "jane@test.com"},
	}
	return fn(chunk)
}

func (m *MockRepository) Each(fn func(interface{}) error) error {
	items := []interface{}{
		map[string]interface{}{"id": 1, "name": "John Doe", "age": 25, "email": "john@test.com"},
		map[string]interface{}{"id": 2, "name": "Jane Smith", "age": 30, "email": "jane@test.com"},
	}
	for _, item := range items {
		if err := fn(item); err != nil {
			return err
		}
	}
	return nil
}

func (m *MockRepository) Pluck(field string) ([]interface{}, error) {
	return []interface{}{25, 30}, nil
}

func (m *MockRepository) Value(field string) (interface{}, error) {
	if field == "age" && len(mockUsers) > 0 {
		return mockUsers[0]["age"], nil
	}
	return nil, nil
}

func (m *MockRepository) Increment(field string, amount interface{}) error {
	if field == "age" {
		for _, user := range mockUsers {
			user["age"] = user["age"].(int) + amount.(int)
		}
	}
	return nil
}

func (m *MockRepository) Decrement(field string, amount interface{}) error {
	if field == "age" {
		for _, user := range mockUsers {
			user["age"] = user["age"].(int) - amount.(int)
		}
	}
	return nil
}
