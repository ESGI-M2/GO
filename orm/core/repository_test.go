package core

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"
)

// TestUserForRepository for testing
type TestUserForRepository struct {
	ID       int    `orm:"pk,auto"`
	Name     string `orm:"index"`
	Email    string `orm:"unique"`
	Age      int
	IsActive bool `orm:"default:true"`
}

// MockDialectForRepository for testing
type MockDialectForRepository struct {
	queryResults []map[string]interface{}
	queryError   error
	countResult  int64
	countError   error
	execResult   sql.Result
	execError    error
}

func (m *MockDialectForRepository) Connect(config ConnectionConfig) error {
	return nil
}

func (m *MockDialectForRepository) Close() error {
	return nil
}

func (m *MockDialectForRepository) Ping() error {
	return nil
}

func (m *MockDialectForRepository) Exec(query string, args ...interface{}) (sql.Result, error) {
	return m.execResult, m.execError
}

func (m *MockDialectForRepository) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if m.queryError != nil {
		return nil, m.queryError
	}
	// Return an error to prevent nil pointer dereference in tests
	return nil, fmt.Errorf("mock query not implemented")
}

func (m *MockDialectForRepository) QueryRow(query string, args ...interface{}) *sql.Row {
	// Return nil to simulate no result
	return nil
}

func (m *MockDialectForRepository) Begin() (Transaction, error) {
	return &MockTransactionForRepository{}, nil
}

func (m *MockDialectForRepository) BeginTx(ctx context.Context, opts *sql.TxOptions) (Transaction, error) {
	return &MockTransactionForRepository{}, nil
}

func (m *MockDialectForRepository) CreateTable(tableName string, columns []Column) error {
	return nil
}

func (m *MockDialectForRepository) DropTable(tableName string) error {
	return nil
}

func (m *MockDialectForRepository) TableExists(tableName string) (bool, error) {
	return false, nil
}

func (m *MockDialectForRepository) GetSQLType(goType reflect.Type) string {
	return "TEXT"
}

func (m *MockDialectForRepository) GetPlaceholder(index int) string {
	return "?"
}

type MockTransactionForRepository struct{}

func (mt *MockTransactionForRepository) Commit() error {
	return nil
}

func (mt *MockTransactionForRepository) Rollback() error {
	return nil
}

func (mt *MockTransactionForRepository) Exec(query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (mt *MockTransactionForRepository) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (mt *MockTransactionForRepository) QueryRow(query string, args ...interface{}) *sql.Row {
	return nil
}

func TestRepositoryImpl_Save(t *testing.T) {
	orm := &ORMImpl{
		dialect:         &MockDialectForRepository{},
		metadataManager: NewMetadataManager(),
		models:          make(map[reflect.Type]*ModelMetadata),
	}

	metadata := &ModelMetadata{
		TableName:  "users",
		Type:       reflect.TypeOf(TestUserForRepository{}),
		PrimaryKey: "ID",
		Columns: []Column{
			{Name: "id", Type: "INT", PrimaryKey: true, AutoIncrement: true},
			{Name: "name", Type: "VARCHAR(255)"},
			{Name: "age", Type: "INT"},
			{Name: "email", Type: "VARCHAR(255)"},
			{Name: "is_active", Type: "BOOLEAN"},
		},
	}

	repo := &RepositoryImpl{
		orm:      orm,
		metadata: metadata,
	}

	user := &TestUserForRepository{
		Name:     "John Doe",
		Age:      30,
		Email:    "john@example.com",
		IsActive: true,
	}

	err := repo.Save(user)
	if err != nil {
		t.Errorf("Save should not return error: %v", err)
	}
}

func TestRepositoryImpl_Find(t *testing.T) {
	orm := &ORMImpl{
		dialect:         &MockDialectForRepository{},
		metadataManager: NewMetadataManager(),
		models:          make(map[reflect.Type]*ModelMetadata),
	}

	metadata := &ModelMetadata{
		TableName:  "users",
		Type:       reflect.TypeOf(TestUserForRepository{}),
		PrimaryKey: "ID",
		Columns: []Column{
			{Name: "id", Type: "INT", PrimaryKey: true, AutoIncrement: true},
			{Name: "name", Type: "VARCHAR(255)"},
			{Name: "age", Type: "INT"},
			{Name: "email", Type: "VARCHAR(255)"},
			{Name: "is_active", Type: "BOOLEAN"},
		},
	}

	repo := &RepositoryImpl{
		orm:      orm,
		metadata: metadata,
	}

	user, err := repo.Find(1)
	if err == nil {
		t.Error("Find should return error when mock dialect returns error")
	}

	if user != nil {
		t.Error("Find should return nil for non-existent record")
	}
}

func TestRepositoryImpl_FindAll(t *testing.T) {
	orm := &ORMImpl{
		dialect:         &MockDialectForRepository{},
		metadataManager: NewMetadataManager(),
		models:          make(map[reflect.Type]*ModelMetadata),
	}

	metadata := &ModelMetadata{
		TableName:  "users",
		Type:       reflect.TypeOf(TestUserForRepository{}),
		PrimaryKey: "ID",
		Columns: []Column{
			{Name: "id", Type: "INT", PrimaryKey: true, AutoIncrement: true},
			{Name: "name", Type: "VARCHAR(255)"},
			{Name: "age", Type: "INT"},
			{Name: "email", Type: "VARCHAR(255)"},
			{Name: "is_active", Type: "BOOLEAN"},
		},
	}

	repo := &RepositoryImpl{
		orm:      orm,
		metadata: metadata,
	}

	users, err := repo.FindAll()
	if err == nil {
		t.Error("FindAll should return error when mock dialect returns error")
	}

	if users != nil {
		t.Error("FindAll should return nil when error occurs")
	}
}

func TestRepositoryImpl_Update(t *testing.T) {
	orm := &ORMImpl{
		dialect:         &MockDialectForRepository{},
		metadataManager: NewMetadataManager(),
		models:          make(map[reflect.Type]*ModelMetadata),
	}

	metadata := &ModelMetadata{
		TableName:  "users",
		Type:       reflect.TypeOf(TestUserForRepository{}),
		PrimaryKey: "ID",
		Columns: []Column{
			{Name: "id", Type: "INT", PrimaryKey: true, AutoIncrement: true},
			{Name: "name", Type: "VARCHAR(255)"},
			{Name: "age", Type: "INT"},
			{Name: "email", Type: "VARCHAR(255)"},
			{Name: "is_active", Type: "BOOLEAN"},
		},
	}

	repo := &RepositoryImpl{
		orm:      orm,
		metadata: metadata,
	}

	user := &TestUserForRepository{
		ID:       1,
		Name:     "John Doe",
		Age:      30,
		Email:    "john@example.com",
		IsActive: true,
	}

	err := repo.Update(user)
	if err != nil {
		t.Errorf("Update should not return error: %v", err)
	}
}

func TestRepositoryImpl_Delete(t *testing.T) {
	orm := &ORMImpl{
		dialect:         &MockDialectForRepository{},
		metadataManager: NewMetadataManager(),
		models:          make(map[reflect.Type]*ModelMetadata),
	}

	metadata := &ModelMetadata{
		TableName:  "users",
		Type:       reflect.TypeOf(TestUserForRepository{}),
		PrimaryKey: "ID",
		Columns: []Column{
			{Name: "id", Type: "INT", PrimaryKey: true, AutoIncrement: true},
			{Name: "name", Type: "VARCHAR(255)"},
			{Name: "age", Type: "INT"},
			{Name: "email", Type: "VARCHAR(255)"},
			{Name: "is_active", Type: "BOOLEAN"},
		},
	}

	repo := &RepositoryImpl{
		orm:      orm,
		metadata: metadata,
	}

	user := &TestUserForRepository{
		ID:       1,
		Name:     "John Doe",
		Age:      30,
		Email:    "john@example.com",
		IsActive: true,
	}

	err := repo.Delete(user)
	if err != nil {
		t.Errorf("Delete should not return error: %v", err)
	}
}

func TestRepositoryImpl_Count(t *testing.T) {
	orm := &ORMImpl{
		dialect:         &MockDialectForRepository{},
		metadataManager: NewMetadataManager(),
		models:          make(map[reflect.Type]*ModelMetadata),
	}

	metadata := &ModelMetadata{
		TableName:  "users",
		Type:       reflect.TypeOf(TestUserForRepository{}),
		PrimaryKey: "ID",
		Columns: []Column{
			{Name: "id", Type: "INT", PrimaryKey: true, AutoIncrement: true},
			{Name: "name", Type: "VARCHAR(255)"},
			{Name: "age", Type: "INT"},
			{Name: "email", Type: "VARCHAR(255)"},
			{Name: "is_active", Type: "BOOLEAN"},
		},
	}

	repo := &RepositoryImpl{
		orm:      orm,
		metadata: metadata,
	}

	count, err := repo.Count()
	if err == nil {
		t.Error("Count should return error when QueryRow returns nil")
	}

	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
}

func TestRepositoryImpl_Exists(t *testing.T) {
	orm := &ORMImpl{
		dialect:         &MockDialectForRepository{},
		metadataManager: NewMetadataManager(),
		models:          make(map[reflect.Type]*ModelMetadata),
	}

	metadata := &ModelMetadata{
		TableName:  "users",
		Type:       reflect.TypeOf(TestUserForRepository{}),
		PrimaryKey: "ID",
		Columns: []Column{
			{Name: "id", Type: "INT", PrimaryKey: true, AutoIncrement: true},
			{Name: "name", Type: "VARCHAR(255)"},
			{Name: "age", Type: "INT"},
			{Name: "email", Type: "VARCHAR(255)"},
			{Name: "is_active", Type: "BOOLEAN"},
		},
	}

	repo := &RepositoryImpl{
		orm:      orm,
		metadata: metadata,
	}

	exists, err := repo.Exists(1)
	if err == nil {
		t.Error("Exists should return error when mock dialect returns error")
	}

	if exists {
		t.Error("Exists should return false when error occurs")
	}
}

func TestRepositoryImpl_FindBy(t *testing.T) {
	orm := &ORMImpl{
		dialect:         &MockDialectForRepository{},
		metadataManager: NewMetadataManager(),
		models:          make(map[reflect.Type]*ModelMetadata),
	}

	metadata := &ModelMetadata{
		TableName:  "users",
		Type:       reflect.TypeOf(TestUserForRepository{}),
		PrimaryKey: "ID",
		Columns: []Column{
			{Name: "id", Type: "INT", PrimaryKey: true, AutoIncrement: true},
			{Name: "name", Type: "VARCHAR(255)"},
			{Name: "age", Type: "INT"},
			{Name: "email", Type: "VARCHAR(255)"},
			{Name: "is_active", Type: "BOOLEAN"},
		},
	}

	repo := &RepositoryImpl{
		orm:      orm,
		metadata: metadata,
	}

	criteria := map[string]interface{}{
		"name": "John Doe",
		"age":  30,
	}

	users, err := repo.FindBy(criteria)
	if err == nil {
		t.Error("FindBy should return error when mock dialect returns error")
	}

	if users != nil {
		t.Error("FindBy should return nil when error occurs")
	}
}

func TestRepositoryImpl_FindOneBy(t *testing.T) {
	orm := &ORMImpl{
		dialect:         &MockDialectForRepository{},
		metadataManager: NewMetadataManager(),
		models:          make(map[reflect.Type]*ModelMetadata),
	}

	metadata := &ModelMetadata{
		TableName:  "users",
		Type:       reflect.TypeOf(TestUserForRepository{}),
		PrimaryKey: "ID",
		Columns: []Column{
			{Name: "id", Type: "INT", PrimaryKey: true, AutoIncrement: true},
			{Name: "name", Type: "VARCHAR(255)"},
			{Name: "age", Type: "INT"},
			{Name: "email", Type: "VARCHAR(255)"},
			{Name: "is_active", Type: "BOOLEAN"},
		},
	}

	repo := &RepositoryImpl{
		orm:      orm,
		metadata: metadata,
	}

	criteria := map[string]interface{}{
		"name": "John Doe",
	}

	user, err := repo.FindOneBy(criteria)
	if err == nil {
		t.Error("FindOneBy should return error when mock dialect returns error")
	}

	if user != nil {
		t.Error("FindOneBy should return nil when error occurs")
	}
}

func TestRepositoryImpl_DeleteBy(t *testing.T) {
	orm := &ORMImpl{
		dialect:         &MockDialectForRepository{},
		metadataManager: NewMetadataManager(),
		models:          make(map[reflect.Type]*ModelMetadata),
	}

	metadata := &ModelMetadata{
		TableName:  "users",
		Type:       reflect.TypeOf(TestUserForRepository{}),
		PrimaryKey: "ID",
		Columns: []Column{
			{Name: "id", Type: "INT", PrimaryKey: true, AutoIncrement: true},
			{Name: "name", Type: "VARCHAR(255)"},
			{Name: "age", Type: "INT"},
			{Name: "email", Type: "VARCHAR(255)"},
			{Name: "is_active", Type: "BOOLEAN"},
		},
	}

	repo := &RepositoryImpl{
		orm:      orm,
		metadata: metadata,
	}

	criteria := map[string]interface{}{
		"name": "John Doe",
	}

	err := repo.DeleteBy(criteria)
	if err != nil {
		t.Errorf("DeleteBy should not return error: %v", err)
	}
}

func TestRepositoryImpl_ErrorHandling(t *testing.T) {
	repo := &RepositoryImpl{
		err: fmt.Errorf("test error"),
	}

	// Test that methods return early when there's an error
	_, err := repo.Find(1)
	if err == nil {
		t.Error("Find should return error when repo has error")
	}

	_, err = repo.FindAll()
	if err == nil {
		t.Error("FindAll should return error when repo has error")
	}

	err = repo.Save(&TestUserForRepository{})
	if err == nil {
		t.Error("Save should return error when repo has error")
	}

	err = repo.Update(&TestUserForRepository{})
	if err == nil {
		t.Error("Update should return error when repo has error")
	}

	err = repo.Delete(&TestUserForRepository{})
	if err == nil {
		t.Error("Delete should return error when repo has error")
	}

	_, err = repo.Count()
	if err == nil {
		t.Error("Count should return error when repo has error")
	}

	_, err = repo.Exists(1)
	if err == nil {
		t.Error("Exists should return error when repo has error")
	}
}

func TestRepositoryImpl_NilMetadata(t *testing.T) {
	repo := &RepositoryImpl{
		metadata: nil,
	}

	// Test that methods return error when metadata is nil
	_, err := repo.Find(1)
	if err == nil {
		t.Error("Find should return error when metadata is nil")
	}

	_, err = repo.FindAll()
	if err == nil {
		t.Error("FindAll should return error when metadata is nil")
	}

	err = repo.Save(&TestUserForRepository{})
	if err == nil {
		t.Error("Save should return error when metadata is nil")
	}

	err = repo.Update(&TestUserForRepository{})
	if err == nil {
		t.Error("Update should return error when metadata is nil")
	}

	err = repo.Delete(&TestUserForRepository{})
	if err == nil {
		t.Error("Delete should return error when metadata is nil")
	}

	_, err = repo.Count()
	if err == nil {
		t.Error("Count should return error when metadata is nil")
	}

	_, err = repo.Exists(1)
	if err == nil {
		t.Error("Exists should return error when metadata is nil")
	}
}
