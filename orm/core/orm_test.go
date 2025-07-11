package core

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"
)

// TestUserForORM for testing
type TestUserForORM struct {
	ID       int    `orm:"pk,auto"`
	Name     string `orm:"index"`
	Email    string `orm:"unique"`
	Age      int
	IsActive bool `orm:"default:true"`
}

// MockDialectForORM for testing
type MockDialectForORM struct {
	queryResults []map[string]interface{}
	queryError   error
	countResult  int64
	countError   error
	execResult   sql.Result
	execError    error
	connected    bool
}

func (m *MockDialectForORM) Connect(config ConnectionConfig) error {
	m.connected = true
	return nil
}

func (m *MockDialectForORM) Close() error {
	m.connected = false
	return nil
}

func (m *MockDialectForORM) Ping() error {
	if !m.connected {
		return fmt.Errorf("not connected")
	}
	return nil
}

func (m *MockDialectForORM) Exec(query string, args ...interface{}) (sql.Result, error) {
	return m.execResult, m.execError
}

func (m *MockDialectForORM) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, m.queryError
}

func (m *MockDialectForORM) QueryRow(query string, args ...interface{}) *sql.Row {
	return nil
}

func (m *MockDialectForORM) Begin() (Transaction, error) {
	return &MockTransactionForORM{}, nil
}

func (m *MockDialectForORM) BeginTx(ctx context.Context, opts *sql.TxOptions) (Transaction, error) {
	return &MockTransactionForORM{}, nil
}

func (m *MockDialectForORM) CreateTable(tableName string, columns []Column) error {
	return nil
}

func (m *MockDialectForORM) DropTable(tableName string) error {
	return nil
}

func (m *MockDialectForORM) TableExists(tableName string) (bool, error) {
	return false, nil
}

func (m *MockDialectForORM) GetSQLType(goType reflect.Type) string {
	return "TEXT"
}

func (m *MockDialectForORM) GetPlaceholder(index int) string {
	return "?"
}

type MockTransactionForORM struct{}

func (mt *MockTransactionForORM) Commit() error {
	return nil
}

func (mt *MockTransactionForORM) Rollback() error {
	return nil
}

func (mt *MockTransactionForORM) Exec(query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (mt *MockTransactionForORM) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (mt *MockTransactionForORM) QueryRow(query string, args ...interface{}) *sql.Row {
	return nil
}

func TestORMImpl_Connect(t *testing.T) {
	orm := &ORMImpl{
		dialect: &MockDialectForORM{},
	}

	config := ConnectionConfig{
		Driver:   "mysql",
		Host:     "localhost",
		Port:     3306,
		Database: "test",
		Username: "root",
		Password: "password",
	}

	err := orm.Connect(config)
	if err != nil {
		t.Errorf("Connect should not return error: %v", err)
	}

	if !orm.IsConnected() {
		t.Error("ORM should be connected after Connect")
	}
}

func TestORMImpl_Close(t *testing.T) {
	orm := &ORMImpl{
		dialect: &MockDialectForORM{},
	}

	err := orm.Close()
	if err != nil {
		t.Errorf("Close should not return error: %v", err)
	}

	if orm.IsConnected() {
		t.Error("ORM should not be connected after Close")
	}
}

func TestORMImpl_RegisterModel(t *testing.T) {
	orm := &ORMImpl{
		dialect:         &MockDialectForORM{},
		metadataManager: NewMetadataManager(),
		models:          make(map[reflect.Type]*ModelMetadata),
	}

	user := &TestUserForORM{}
	err := orm.RegisterModel(user)
	if err != nil {
		t.Errorf("RegisterModel should not return error: %v", err)
	}

	metadata, err := orm.GetMetadata(user)
	if err != nil {
		t.Errorf("GetMetadata should not return error: %v", err)
	}

	if metadata == nil {
		t.Error("GetMetadata should return metadata")
	}

	if metadata.TableName != "testuserfororm" {
		t.Errorf("Expected table name 'testuserfororm', got '%s'", metadata.TableName)
	}
}

func TestORMImpl_Query(t *testing.T) {
	orm := &ORMImpl{
		dialect:         &MockDialectForORM{},
		metadataManager: NewMetadataManager(),
		models:          make(map[reflect.Type]*ModelMetadata),
	}

	user := &TestUserForORM{}
	qb := orm.Query(user)
	if qb == nil {
		t.Error("Query should return QueryBuilder")
	}
}

func TestORMImpl_Raw(t *testing.T) {
	orm := &ORMImpl{
		dialect: &MockDialectForORM{},
	}

	qb := orm.Raw("SELECT * FROM users WHERE age > ?", 25)
	if qb == nil {
		t.Error("Raw should return QueryBuilder")
	}
}

func TestORMImpl_Repository(t *testing.T) {
	orm := &ORMImpl{
		dialect:         &MockDialectForORM{},
		metadataManager: NewMetadataManager(),
		models:          make(map[reflect.Type]*ModelMetadata),
	}

	user := &TestUserForORM{}
	repo := orm.Repository(user)
	if repo == nil {
		t.Error("Repository should return Repository")
	}
}

func TestORMImpl_Transaction(t *testing.T) {
	orm := &ORMImpl{
		dialect: &MockDialectForORM{},
	}

	err := orm.Transaction(func(txORM ORM) error {
		// Test transaction function
		return nil
	})
	if err != nil {
		t.Errorf("Transaction should not return error: %v", err)
	}
}

func TestORMImpl_TransactionWithContext(t *testing.T) {
	orm := &ORMImpl{
		dialect: &MockDialectForORM{},
	}

	ctx := context.Background()
	err := orm.TransactionWithContext(ctx, func(txORM ORM) error {
		// Test transaction function with context
		return nil
	})
	if err != nil {
		t.Errorf("TransactionWithContext should not return error: %v", err)
	}
}

func TestORMImpl_CreateTable(t *testing.T) {
	orm := &ORMImpl{
		dialect:         &MockDialectForORM{},
		metadataManager: NewMetadataManager(),
		models:          make(map[reflect.Type]*ModelMetadata),
	}

	user := &TestUserForORM{}
	err := orm.CreateTable(user)
	if err != nil {
		t.Errorf("CreateTable should not return error: %v", err)
	}
}

func TestORMImpl_DropTable(t *testing.T) {
	orm := &ORMImpl{
		dialect:         &MockDialectForORM{},
		metadataManager: NewMetadataManager(),
		models:          make(map[reflect.Type]*ModelMetadata),
	}

	user := &TestUserForORM{}
	err := orm.DropTable(user)
	if err != nil {
		t.Errorf("DropTable should not return error: %v", err)
	}
}

func TestORMImpl_Migrate(t *testing.T) {
	orm := &ORMImpl{
		dialect:         &MockDialectForORM{},
		metadataManager: NewMetadataManager(),
		models:          make(map[reflect.Type]*ModelMetadata),
	}

	err := orm.Migrate()
	if err != nil {
		t.Errorf("Migrate should not return error: %v", err)
	}
}

func TestORMImpl_GetDialect(t *testing.T) {
	dialect := &MockDialectForORM{}
	orm := &ORMImpl{
		dialect: dialect,
	}

	result := orm.GetDialect()
	if result != dialect {
		t.Error("GetDialect should return the dialect")
	}
}

func TestORMImpl_IsConnected(t *testing.T) {
	dialect := &MockDialectForORM{}
	orm := &ORMImpl{
		dialect: dialect,
		models:  make(map[reflect.Type]*ModelMetadata),
	}

	// Test when not connected
	if orm.IsConnected() {
		t.Error("ORM should not be connected initially")
	}

	// Test when connected
	orm.connected = true
	if !orm.IsConnected() {
		t.Error("ORM should be connected when ORM is connected")
	}
}

func TestORMImpl_GetMetadata(t *testing.T) {
	orm := &ORMImpl{
		dialect:         &MockDialectForORM{},
		metadataManager: NewMetadataManager(),
		models:          make(map[reflect.Type]*ModelMetadata),
	}

	user := &TestUserForORM{}

	// Test with unregistered model
	metadata, err := orm.GetMetadata(user)
	if err != nil {
		t.Errorf("GetMetadata should not return error: %v", err)
	}

	if metadata == nil {
		t.Error("GetMetadata should return metadata even for unregistered model")
	}

	// Test with registered model
	err = orm.RegisterModel(user)
	if err != nil {
		t.Errorf("RegisterModel should not return error: %v", err)
	}

	metadata, err = orm.GetMetadata(user)
	if err != nil {
		t.Errorf("GetMetadata should not return error: %v", err)
	}

	if metadata == nil {
		t.Error("GetMetadata should return metadata for registered model")
	}
}

func TestORMImpl_ErrorHandling(t *testing.T) {
	// Test with nil dialect
	orm := &ORMImpl{
		dialect: nil,
		models:  make(map[reflect.Type]*ModelMetadata),
	}

	err := orm.Connect(ConnectionConfig{})
	if err == nil {
		t.Error("Connect should return error when dialect is nil")
	}

	err = orm.Close()
	if err != nil {
		t.Errorf("Close should not return error when dialect is nil: %v", err)
	}

	if orm.IsConnected() {
		t.Error("IsConnected should return false when dialect is nil")
	}

	// Test with error-prone dialect
	mockDialect := &MockDialectForORM{
		execError:  fmt.Errorf("exec error"),
		queryError: fmt.Errorf("query error"),
	}

	orm = &ORMImpl{
		dialect: mockDialect,
	}

	err = orm.Connect(ConnectionConfig{})
	if err != nil {
		t.Errorf("Connect should not return error: %v", err)
	}
}

func TestORMImpl_TransactionErrorHandling(t *testing.T) {
	orm := &ORMImpl{
		dialect: &MockDialectForORM{},
	}

	// Test transaction function that returns error
	err := orm.Transaction(func(txORM ORM) error {
		return fmt.Errorf("transaction error")
	})
	if err == nil {
		t.Error("Transaction should return error when function returns error")
	}

	// Test transaction with context that returns error
	ctx := context.Background()
	err = orm.TransactionWithContext(ctx, func(txORM ORM) error {
		return fmt.Errorf("transaction context error")
	})
	if err == nil {
		t.Error("TransactionWithContext should return error when function returns error")
	}
}

func TestORMImpl_ModelRegistration(t *testing.T) {
	orm := &ORMImpl{
		dialect:         &MockDialectForORM{},
		metadataManager: NewMetadataManager(),
		models:          make(map[reflect.Type]*ModelMetadata),
	}

	// Test registering multiple models
	user1 := &TestUserForORM{}
	user2 := &TestUserForORM{}

	err := orm.RegisterModel(user1)
	if err != nil {
		t.Errorf("RegisterModel should not return error: %v", err)
	}

	err = orm.RegisterModel(user2)
	if err != nil {
		t.Errorf("RegisterModel should not return error: %v", err)
	}

	// Test getting metadata for both models
	metadata1, err := orm.GetMetadata(user1)
	if err != nil {
		t.Errorf("GetMetadata should not return error: %v", err)
	}

	metadata2, err := orm.GetMetadata(user2)
	if err != nil {
		t.Errorf("GetMetadata should not return error: %v", err)
	}

	if metadata1 == nil || metadata2 == nil {
		t.Error("GetMetadata should return metadata for both models")
	}
}

func TestORMImpl_QueryBuilderIntegration(t *testing.T) {
	orm := &ORMImpl{
		dialect:         &MockDialectForORM{},
		metadataManager: NewMetadataManager(),
		models:          make(map[reflect.Type]*ModelMetadata),
	}

	user := &TestUserForORM{}
	qb := orm.Query(user)

	// Test query builder methods
	result := qb.Select("name", "age").Where("age", ">", 25).OrderBy("name", "ASC")
	if result == nil {
		t.Error("QueryBuilder methods should return QueryBuilder")
	}

	sql := result.GetSQL()
	if sql == "" {
		t.Error("GetSQL should return SQL string")
	}

	args := result.GetArgs()
	if args == nil {
		t.Error("GetArgs should return arguments")
	}
}

func TestORMImpl_RepositoryIntegration(t *testing.T) {
	orm := &ORMImpl{
		dialect:         &MockDialectForORM{},
		metadataManager: NewMetadataManager(),
		models:          make(map[reflect.Type]*ModelMetadata),
	}

	user := &TestUserForORM{}
	repo := orm.Repository(user)

	// Test repository creation
	if repo == nil {
		t.Error("Repository should not be nil")
	}

	// Test that repository methods exist (without executing queries)
	// The actual query execution is tested in separate repository tests
	repoType := reflect.TypeOf(repo)

	// Check that repository has required methods
	_, hasFind := repoType.MethodByName("Find")
	if !hasFind {
		t.Error("Repository should have Find method")
	}

	_, hasFindAll := repoType.MethodByName("FindAll")
	if !hasFindAll {
		t.Error("Repository should have FindAll method")
	}

	_, hasSave := repoType.MethodByName("Save")
	if !hasSave {
		t.Error("Repository should have Save method")
	}
}
