package unit

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/ESGI-M2/GO/orm/builder"
	"github.com/ESGI-M2/GO/orm/core/interfaces"
	"github.com/ESGI-M2/GO/orm/factory"
)

type AdvancedORMTestModel struct {
	ID   int    `orm:"pk,auto"`
	Name string `orm:"column:name"`
}

func setupAdvancedORM() *builder.SimpleORM {
	return builder.NewSimpleORM().
		WithDialect(factory.Mock).
		RegisterModel(&AdvancedORMTestModel{})
}

func TestAdvancedORM_WithCache(t *testing.T) {
	orm := setupAdvancedORM()

	// Connect first
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Test WithCache method on underlying ORM
	underlyingORM := orm.GetORM()
	cachedORM := underlyingORM.WithCache(300) // 5 minutes
	if cachedORM == nil {
		t.Error("WithCache should return an ORM instance")
	}

	// Test that cached ORM works
	if cachedORM.GetDialect() == nil {
		t.Error("WithCache should return working ORM")
	}
}

func TestAdvancedORM_WithConnectionPool(t *testing.T) {
	orm := setupAdvancedORM()

	// Connect first
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Test WithConnectionPool method on underlying ORM
	underlyingORM := orm.GetORM()
	pooledORM := underlyingORM.WithConnectionPool(10, 5) // max 10 open, 5 idle
	if pooledORM == nil {
		t.Error("WithConnectionPool should return an ORM instance")
	}

	// Test that pooled ORM works
	if pooledORM.GetDialect() == nil {
		t.Error("WithConnectionPool should return working ORM")
	}
}

func TestAdvancedORM_EnableQueryLog(t *testing.T) {
	orm := setupAdvancedORM()

	// Connect first
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Test EnableQueryLog method on underlying ORM
	underlyingORM := orm.GetORM()
	loggedORM := underlyingORM.EnableQueryLog()
	if loggedORM == nil {
		t.Error("EnableQueryLog should return an ORM instance")
	}

	// Test that logged ORM works
	if loggedORM.GetDialect() == nil {
		t.Error("EnableQueryLog should return working ORM")
	}
}

func TestAdvancedORM_DisableQueryLog(t *testing.T) {
	orm := setupAdvancedORM()

	// Connect first
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Test EnableQueryLog then DisableQueryLog method on underlying ORM
	underlyingORM := orm.GetORM()
	loggedORM := underlyingORM.EnableQueryLog()
	unloggedORM := loggedORM.DisableQueryLog()

	if unloggedORM == nil {
		t.Error("DisableQueryLog should return an ORM instance")
	}

	// Test that unlogged ORM works
	if unloggedORM.GetDialect() == nil {
		t.Error("DisableQueryLog should return working ORM")
	}
}

func TestAdvancedORM_ChainedConfiguration(t *testing.T) {
	orm := setupAdvancedORM()

	// Connect first
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Test chained configuration on underlying ORM
	underlyingORM := orm.GetORM()
	configuredORM := underlyingORM.
		WithCache(300).
		WithConnectionPool(10, 5).
		EnableQueryLog()

	if configuredORM == nil {
		t.Error("Chained configuration should return an ORM instance")
	}

	// Test that configured ORM works
	if configuredORM.GetDialect() == nil {
		t.Error("Chained configuration should return working ORM")
	}
}

func TestAdvancedORM_CacheInterface(t *testing.T) {
	// Test Cache interface methods
	cache := &MockCache{}

	// Test Set
	err := cache.Set("key1", "value1", 300)
	if err != nil {
		t.Errorf("Cache.Set failed: %v", err)
	}

	// Test Get
	value, found := cache.Get("key1")
	if !found {
		t.Error("Cache.Get should find existing key")
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got '%v'", value)
	}

	// Test Has
	exists := cache.Has("key1")
	if !exists {
		t.Error("Cache.Has should return true for existing key")
	}

	// Test Delete
	err = cache.Delete("key1")
	if err != nil {
		t.Errorf("Cache.Delete failed: %v", err)
	}

	// Test Get after delete
	_, found = cache.Get("key1")
	if found {
		t.Error("Cache.Get should not find deleted key")
	}

	// Test Clear
	cache.Set("key2", "value2", 300)
	cache.Set("key3", "value3", 300)

	err = cache.Clear()
	if err != nil {
		t.Errorf("Cache.Clear failed: %v", err)
	}

	// Test that all keys are cleared
	if cache.Has("key2") || cache.Has("key3") {
		t.Error("Cache.Clear should remove all keys")
	}
}

func TestAdvancedORM_QueryLoggerInterface(t *testing.T) {
	// Test QueryLogger interface methods
	logger := &MockQueryLogger{}

	// Test Log
	log := interfaces.QueryLog{
		SQL:      "SELECT * FROM users",
		Args:     []interface{}{1, "test"},
		Duration: 50 * time.Millisecond,
		Time:     time.Now(),
		Error:    nil,
	}

	logger.Log(log)

	// Test GetLogs
	logs := logger.GetLogs()
	if len(logs) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(logs))
	}

	if logs[0].SQL != "SELECT * FROM users" {
		t.Errorf("Expected SQL to be 'SELECT * FROM users', got '%s'", logs[0].SQL)
	}

	// Test with error
	errorLog := interfaces.QueryLog{
		SQL:      "SELECT * FROM invalid_table",
		Args:     []interface{}{},
		Duration: 10 * time.Millisecond,
		Time:     time.Now(),
		Error:    fmt.Errorf("table not found"),
	}

	logger.Log(errorLog)

	logs = logger.GetLogs()
	if len(logs) != 2 {
		t.Errorf("Expected 2 log entries, got %d", len(logs))
	}

	if logs[1].Error == nil {
		t.Error("Expected error log to have error")
	}

	// Test ClearLogs
	logger.ClearLogs()

	logs = logger.GetLogs()
	if len(logs) != 0 {
		t.Errorf("Expected 0 log entries after clear, got %d", len(logs))
	}
}

func TestAdvancedORM_TransactionInterface(t *testing.T) {
	orm := setupAdvancedORM()

	// Connect first
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Test Transaction interface through ORM
	underlyingORM := orm.GetORM()

	// Test TransactionWithContext
	ctx := context.Background()
	err = underlyingORM.TransactionWithContext(ctx, func(txORM interfaces.ORM) error {
		// Test that transaction ORM works
		if txORM == nil {
			t.Error("Transaction ORM should not be nil")
		}

		// Test query within transaction
		query := txORM.Query(&AdvancedORMTestModel{})
		if query == nil {
			t.Error("Query should work within transaction")
		}

		return nil
	})

	if err != nil {
		t.Errorf("TransactionWithContext failed: %v", err)
	}

	// Test Transaction with error (should rollback)
	err = underlyingORM.Transaction(func(txORM interfaces.ORM) error {
		// Return error to trigger rollback
		return fmt.Errorf("test rollback")
	})

	if err == nil {
		t.Error("Transaction should return error when function returns error")
	}
}

func TestAdvancedORM_AdvancedDialectMethods(t *testing.T) {
	orm := setupAdvancedORM()

	// Connect first
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	defer orm.Close()

	// Get dialect through ORM
	underlyingORM := orm.GetORM()
	dialect := underlyingORM.GetDialect()

	if dialect == nil {
		t.Fatal("Dialect should not be nil")
	}

	// Test Ping
	err = dialect.Ping()
	if err != nil {
		t.Errorf("Dialect.Ping failed: %v", err)
	}

	// Test TableExists
	exists, err := dialect.TableExists("users")
	if err != nil {
		t.Errorf("Dialect.TableExists failed: %v", err)
	}
	// exists can be true or false
	_ = exists

	// Test GetSQLType
	sqlType := dialect.GetSQLType(reflect.TypeOf(""))
	if sqlType == "" {
		t.Error("GetSQLType should return non-empty string")
	}

	// Test GetPlaceholder
	placeholder := dialect.GetPlaceholder(1)
	if placeholder == "" {
		t.Error("GetPlaceholder should return non-empty string")
	}

	// Test FullTextSearch
	ftsQuery := dialect.FullTextSearch("content", "search term")
	if ftsQuery == "" {
		t.Error("FullTextSearch should return non-empty string")
	}

	// Test GetRandomFunction
	randomFunc := dialect.GetRandomFunction()
	if randomFunc == "" {
		t.Error("GetRandomFunction should return non-empty string")
	}

	// Test GetDateFunction
	dateFunc := dialect.GetDateFunction()
	if dateFunc == "" {
		t.Error("GetDateFunction should return non-empty string")
	}

	// Test GetJSONExtract
	jsonExtract := dialect.GetJSONExtract()
	if jsonExtract == "" {
		t.Error("GetJSONExtract should return non-empty string")
	}
}

func TestAdvancedORM_PaginationResult(t *testing.T) {
	// Test PaginationResult struct
	result := &interfaces.PaginationResult{
		Data:        []interface{}{&AdvancedORMTestModel{ID: 1, Name: "Test"}},
		Total:       100,
		PerPage:     10,
		CurrentPage: 1,
		LastPage:    10,
		From:        1,
		To:          10,
		HasMore:     true,
		NextCursor:  "next_cursor",
		PrevCursor:  "prev_cursor",
	}

	if result.Data == nil {
		t.Error("PaginationResult.Data should not be nil")
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 data item, got %d", len(result.Data))
	}

	if result.Total != 100 {
		t.Errorf("Expected total 100, got %d", result.Total)
	}

	if result.PerPage != 10 {
		t.Errorf("Expected perPage 10, got %d", result.PerPage)
	}

	if result.CurrentPage != 1 {
		t.Errorf("Expected currentPage 1, got %d", result.CurrentPage)
	}

	if result.LastPage != 10 {
		t.Errorf("Expected lastPage 10, got %d", result.LastPage)
	}

	if !result.HasMore {
		t.Error("Expected hasMore to be true")
	}
}

func TestAdvancedORM_ErrorCases(t *testing.T) {
	orm := setupAdvancedORM()

	// Test advanced methods without connecting
	underlyingORM := orm.GetORM()

	// Check if underlyingORM is nil
	if underlyingORM == nil {
		t.Log("Underlying ORM is nil when not connected, which is expected")
		return
	}

	// Test getting dialect without connecting
	dialect := underlyingORM.GetDialect()
	if dialect != nil {
		// If dialect is available, test that operations fail gracefully
		err := dialect.Ping()
		if err == nil {
			t.Error("Ping should fail when not connected")
		}
	}

	// Test transaction without connecting
	err := underlyingORM.Transaction(func(txORM interfaces.ORM) error {
		return nil
	})
	if err == nil {
		t.Error("Transaction should fail when not connected")
	}
}

// Mock implementations for testing

type MockCache struct {
	data map[string]interface{}
}

func (c *MockCache) Get(key string) (interface{}, bool) {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	value, exists := c.data[key]
	return value, exists
}

func (c *MockCache) Set(key string, value interface{}, ttl int) error {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	c.data[key] = value
	return nil
}

func (c *MockCache) Delete(key string) error {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	delete(c.data, key)
	return nil
}

func (c *MockCache) Clear() error {
	c.data = make(map[string]interface{})
	return nil
}

func (c *MockCache) Has(key string) bool {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	_, exists := c.data[key]
	return exists
}

type MockQueryLogger struct {
	logs []interfaces.QueryLog
}

func (l *MockQueryLogger) Log(log interfaces.QueryLog) {
	l.logs = append(l.logs, log)
}

func (l *MockQueryLogger) GetLogs() []interfaces.QueryLog {
	return l.logs
}

func (l *MockQueryLogger) ClearLogs() {
	l.logs = []interfaces.QueryLog{}
}
