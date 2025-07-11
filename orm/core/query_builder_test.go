package core

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"
)

// MockDialect for testing
type MockDialectForTest struct {
	queryResults []map[string]interface{}
	queryError   error
	countResult  int64
	countError   error
}

func (m *MockDialectForTest) Connect(config ConnectionConfig) error {
	return nil
}

func (m *MockDialectForTest) Close() error {
	return nil
}

func (m *MockDialectForTest) Ping() error {
	return nil
}

func (m *MockDialectForTest) Exec(query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (m *MockDialectForTest) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, m.queryError
}

func (m *MockDialectForTest) QueryRow(query string, args ...interface{}) *sql.Row {
	return nil
}

func (m *MockDialectForTest) Begin() (Transaction, error) {
	return &MockTransactionForTest{}, nil
}

func (m *MockDialectForTest) BeginTx(ctx context.Context, opts *sql.TxOptions) (Transaction, error) {
	return &MockTransactionForTest{}, nil
}

func (m *MockDialectForTest) CreateTable(tableName string, columns []Column) error {
	return nil
}

func (m *MockDialectForTest) DropTable(tableName string) error {
	return nil
}

func (m *MockDialectForTest) TableExists(tableName string) (bool, error) {
	return false, nil
}

func (m *MockDialectForTest) GetSQLType(goType reflect.Type) string {
	return "TEXT"
}

func (m *MockDialectForTest) GetPlaceholder(index int) string {
	return "?"
}

type MockTransactionForTest struct{}

func (mt *MockTransactionForTest) Commit() error {
	return nil
}

func (mt *MockTransactionForTest) Rollback() error {
	return nil
}

func (mt *MockTransactionForTest) Exec(query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (mt *MockTransactionForTest) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (mt *MockTransactionForTest) QueryRow(query string, args ...interface{}) *sql.Row {
	return nil
}

// TestUserForQueryBuilder for testing
type TestUserForQueryBuilder struct {
	ID       int    `orm:"pk,auto"`
	Name     string `orm:"index"`
	Email    string `orm:"unique"`
	Age      int
	IsActive bool `orm:"default:true"`
}

func TestQueryBuilderImpl_Select(t *testing.T) {
	orm := &ORMImpl{
		dialect: &MockDialectForTest{},
	}

	metadata := &ModelMetadata{
		TableName: "users",
		Type:      reflect.TypeOf(TestUserForQueryBuilder{}),
	}

	qb := &QueryBuilderImpl{
		orm:      orm,
		metadata: metadata,
		table:    "users",
		fields:   []string{"*"},
		where:    make([]WhereCondition, 0),
		orderBy:  make([]OrderBy, 0),
		joins:    make([]Join, 0),
		args:     make([]interface{}, 0),
	}

	// Test with fields
	result := qb.Select("name", "age")
	if result == nil {
		t.Error("Select should return QueryBuilder")
	}

	// Test with no fields
	result = qb.Select()
	if result == nil {
		t.Error("Select with no fields should return QueryBuilder")
	}
}

func TestQueryBuilderImpl_From(t *testing.T) {
	orm := &ORMImpl{
		dialect: &MockDialectForTest{},
	}

	qb := &QueryBuilderImpl{
		orm:     orm,
		table:   "",
		fields:  []string{"*"},
		where:   make([]WhereCondition, 0),
		orderBy: make([]OrderBy, 0),
		joins:   make([]Join, 0),
		args:    make([]interface{}, 0),
	}

	result := qb.From("users")
	if result == nil {
		t.Error("From should return QueryBuilder")
	}

	if qb.table != "users" {
		t.Error("Table should be set to 'users'")
	}
}

func TestQueryBuilderImpl_Where(t *testing.T) {
	orm := &ORMImpl{
		dialect: &MockDialectForTest{},
	}

	qb := &QueryBuilderImpl{
		orm:     orm,
		table:   "users",
		fields:  []string{"*"},
		where:   make([]WhereCondition, 0),
		orderBy: make([]OrderBy, 0),
		joins:   make([]Join, 0),
		args:    make([]interface{}, 0),
	}

	// Test with value
	result := qb.Where("age", ">", 25)
	if result == nil {
		t.Error("Where should return QueryBuilder")
	}

	if len(qb.where) != 1 {
		t.Error("Should have one where condition")
	}

	if len(qb.args) != 1 {
		t.Error("Should have one argument")
	}

	// Test with nil value
	result = qb.Where("name", "=", nil)
	if result == nil {
		t.Error("Where with nil should return QueryBuilder")
	}

	if len(qb.args) != 1 {
		t.Error("Should not add argument for nil value")
	}
}

func TestQueryBuilderImpl_WhereIn(t *testing.T) {
	orm := &ORMImpl{
		dialect: &MockDialectForTest{},
	}

	qb := &QueryBuilderImpl{
		orm:     orm,
		table:   "users",
		fields:  []string{"*"},
		where:   make([]WhereCondition, 0),
		orderBy: make([]OrderBy, 0),
		joins:   make([]Join, 0),
		args:    make([]interface{}, 0),
	}

	// Test with values
	result := qb.WhereIn("age", []interface{}{25, 30, 35})
	if result == nil {
		t.Error("WhereIn should return QueryBuilder")
	}

	if len(qb.where) != 1 {
		t.Error("Should have one where condition")
	}

	if len(qb.args) != 3 {
		t.Error("Should have three arguments")
	}

	// Test with empty values
	result = qb.WhereIn("age", []interface{}{})
	if result == nil {
		t.Error("WhereIn with empty values should return QueryBuilder")
	}

	if len(qb.where) != 1 {
		t.Error("Should not add condition for empty values")
	}
}

func TestQueryBuilderImpl_WhereNotIn(t *testing.T) {
	orm := &ORMImpl{
		dialect: &MockDialectForTest{},
	}

	qb := &QueryBuilderImpl{
		orm:     orm,
		table:   "users",
		fields:  []string{"*"},
		where:   make([]WhereCondition, 0),
		orderBy: make([]OrderBy, 0),
		joins:   make([]Join, 0),
		args:    make([]interface{}, 0),
	}

	// Test with values
	result := qb.WhereNotIn("age", []interface{}{25, 30, 35})
	if result == nil {
		t.Error("WhereNotIn should return QueryBuilder")
	}

	if len(qb.where) != 1 {
		t.Error("Should have one where condition")
	}

	if len(qb.args) != 3 {
		t.Error("Should have three arguments")
	}
}

func TestQueryBuilderImpl_OrderBy(t *testing.T) {
	orm := &ORMImpl{
		dialect: &MockDialectForTest{},
	}

	qb := &QueryBuilderImpl{
		orm:     orm,
		table:   "users",
		fields:  []string{"*"},
		where:   make([]WhereCondition, 0),
		orderBy: make([]OrderBy, 0),
		joins:   make([]Join, 0),
		args:    make([]interface{}, 0),
	}

	// Test with direction
	result := qb.OrderBy("name", "ASC")
	if result == nil {
		t.Error("OrderBy should return QueryBuilder")
	}

	if len(qb.orderBy) != 1 {
		t.Error("Should have one order by clause")
	}

	// Test without direction
	result = qb.OrderBy("age", "")
	if result == nil {
		t.Error("OrderBy without direction should return QueryBuilder")
	}

	if len(qb.orderBy) != 2 {
		t.Error("Should have two order by clauses")
	}

	if qb.orderBy[1].Direction != "ASC" {
		t.Error("Default direction should be ASC")
	}
}

func TestQueryBuilderImpl_GroupBy(t *testing.T) {
	orm := &ORMImpl{
		dialect: &MockDialectForTest{},
	}

	qb := &QueryBuilderImpl{
		orm:     orm,
		table:   "users",
		fields:  []string{"*"},
		where:   make([]WhereCondition, 0),
		orderBy: make([]OrderBy, 0),
		groupBy: make([]string, 0),
		joins:   make([]Join, 0),
		args:    make([]interface{}, 0),
	}

	result := qb.GroupBy("age", "name")
	if result == nil {
		t.Error("GroupBy should return QueryBuilder")
	}

	if len(qb.groupBy) != 2 {
		t.Error("Should have two group by fields")
	}
}

func TestQueryBuilderImpl_Having(t *testing.T) {
	orm := &ORMImpl{
		dialect: &MockDialectForTest{},
	}

	qb := &QueryBuilderImpl{
		orm:     orm,
		table:   "users",
		fields:  []string{"*"},
		where:   make([]WhereCondition, 0),
		orderBy: make([]OrderBy, 0),
		joins:   make([]Join, 0),
		args:    make([]interface{}, 0),
	}

	result := qb.Having("COUNT(*) > ?", 5)
	if result == nil {
		t.Error("Having should return QueryBuilder")
	}

	if qb.having != "COUNT(*) > ?" {
		t.Error("Having clause should be set")
	}

	if len(qb.havingArgs) != 1 {
		t.Error("Should have one having argument")
	}
}

func TestQueryBuilderImpl_Limit(t *testing.T) {
	orm := &ORMImpl{
		dialect: &MockDialectForTest{},
	}

	qb := &QueryBuilderImpl{
		orm:     orm,
		table:   "users",
		fields:  []string{"*"},
		where:   make([]WhereCondition, 0),
		orderBy: make([]OrderBy, 0),
		joins:   make([]Join, 0),
		args:    make([]interface{}, 0),
	}

	result := qb.Limit(10)
	if result == nil {
		t.Error("Limit should return QueryBuilder")
	}

	if qb.limit != 10 {
		t.Error("Limit should be set to 10")
	}
}

func TestQueryBuilderImpl_Offset(t *testing.T) {
	orm := &ORMImpl{
		dialect: &MockDialectForTest{},
	}

	qb := &QueryBuilderImpl{
		orm:     orm,
		table:   "users",
		fields:  []string{"*"},
		where:   make([]WhereCondition, 0),
		orderBy: make([]OrderBy, 0),
		joins:   make([]Join, 0),
		args:    make([]interface{}, 0),
	}

	result := qb.Offset(20)
	if result == nil {
		t.Error("Offset should return QueryBuilder")
	}

	if qb.offset != 20 {
		t.Error("Offset should be set to 20")
	}
}

func TestQueryBuilderImpl_Join(t *testing.T) {
	orm := &ORMImpl{
		dialect: &MockDialectForTest{},
	}

	qb := &QueryBuilderImpl{
		orm:     orm,
		table:   "users",
		fields:  []string{"*"},
		where:   make([]WhereCondition, 0),
		orderBy: make([]OrderBy, 0),
		joins:   make([]Join, 0),
		args:    make([]interface{}, 0),
	}

	result := qb.Join("posts", "users.id = posts.user_id")
	if result == nil {
		t.Error("Join should return QueryBuilder")
	}

	if len(qb.joins) != 1 {
		t.Error("Should have one join")
	}

	if qb.joins[0].Type != "INNER" {
		t.Error("Join type should be INNER")
	}
}

func TestQueryBuilderImpl_LeftJoin(t *testing.T) {
	orm := &ORMImpl{
		dialect: &MockDialectForTest{},
	}

	qb := &QueryBuilderImpl{
		orm:     orm,
		table:   "users",
		fields:  []string{"*"},
		where:   make([]WhereCondition, 0),
		orderBy: make([]OrderBy, 0),
		joins:   make([]Join, 0),
		args:    make([]interface{}, 0),
	}

	result := qb.LeftJoin("posts", "users.id = posts.user_id")
	if result == nil {
		t.Error("LeftJoin should return QueryBuilder")
	}

	if len(qb.joins) != 1 {
		t.Error("Should have one join")
	}

	if qb.joins[0].Type != "LEFT" {
		t.Error("Join type should be LEFT")
	}
}

func TestQueryBuilderImpl_RightJoin(t *testing.T) {
	orm := &ORMImpl{
		dialect: &MockDialectForTest{},
	}

	qb := &QueryBuilderImpl{
		orm:     orm,
		table:   "users",
		fields:  []string{"*"},
		where:   make([]WhereCondition, 0),
		orderBy: make([]OrderBy, 0),
		joins:   make([]Join, 0),
		args:    make([]interface{}, 0),
	}

	result := qb.RightJoin("posts", "users.id = posts.user_id")
	if result == nil {
		t.Error("RightJoin should return QueryBuilder")
	}

	if len(qb.joins) != 1 {
		t.Error("Should have one join")
	}

	if qb.joins[0].Type != "RIGHT" {
		t.Error("Join type should be RIGHT")
	}
}

func TestQueryBuilderImpl_InnerJoin(t *testing.T) {
	orm := &ORMImpl{
		dialect: &MockDialectForTest{},
	}

	qb := &QueryBuilderImpl{
		orm:     orm,
		table:   "users",
		fields:  []string{"*"},
		where:   make([]WhereCondition, 0),
		orderBy: make([]OrderBy, 0),
		joins:   make([]Join, 0),
		args:    make([]interface{}, 0),
	}

	result := qb.InnerJoin("posts", "users.id = posts.user_id")
	if result == nil {
		t.Error("InnerJoin should return QueryBuilder")
	}

	if len(qb.joins) != 1 {
		t.Error("Should have one join")
	}

	if qb.joins[0].Type != "INNER" {
		t.Error("Join type should be INNER")
	}
}

func TestQueryBuilderImpl_Raw(t *testing.T) {
	orm := &ORMImpl{
		dialect: &MockDialectForTest{},
	}

	qb := &QueryBuilderImpl{
		orm:     orm,
		table:   "users",
		fields:  []string{"*"},
		where:   make([]WhereCondition, 0),
		orderBy: make([]OrderBy, 0),
		joins:   make([]Join, 0),
		args:    make([]interface{}, 0),
	}

	result := qb.Raw("SELECT * FROM users WHERE age > ?", 25)
	if result == nil {
		t.Error("Raw should return QueryBuilder")
	}

	if qb.rawSQL != "SELECT * FROM users WHERE age > ?" {
		t.Error("Raw SQL should be set")
	}

	if len(qb.rawArgs) != 1 {
		t.Error("Should have one raw argument")
	}
}

func TestQueryBuilderImpl_GetSQL(t *testing.T) {
	orm := &ORMImpl{
		dialect: &MockDialectForTest{},
	}

	qb := &QueryBuilderImpl{
		orm:     orm,
		table:   "users",
		fields:  []string{"name", "age"},
		where:   make([]WhereCondition, 0),
		orderBy: make([]OrderBy, 0),
		joins:   make([]Join, 0),
		args:    make([]interface{}, 0),
	}

	sql := qb.GetSQL()
	expected := "SELECT name, age FROM users"
	if sql != expected {
		t.Errorf("Expected SQL '%s', got '%s'", expected, sql)
	}

	// Test with raw SQL
	qb.rawSQL = "SELECT COUNT(*) FROM users"
	sql = qb.GetSQL()
	if sql != "SELECT COUNT(*) FROM users" {
		t.Error("Should return raw SQL when set")
	}
}

func TestQueryBuilderImpl_GetArgs(t *testing.T) {
	orm := &ORMImpl{
		dialect: &MockDialectForTest{},
	}

	qb := &QueryBuilderImpl{
		orm:     orm,
		table:   "users",
		fields:  []string{"*"},
		where:   make([]WhereCondition, 0),
		orderBy: make([]OrderBy, 0),
		joins:   make([]Join, 0),
		args:    []interface{}{25, "John"},
	}

	args := qb.GetArgs()
	if len(args) != 2 {
		t.Error("Should return correct number of arguments")
	}

	// Test with raw args
	qb.rawSQL = "SELECT * FROM users WHERE age > ?"
	qb.rawArgs = []interface{}{30}
	args = qb.GetArgs()
	if len(args) != 1 {
		t.Error("Should return raw arguments when set")
	}
}

func TestQueryBuilderImpl_buildQuery(t *testing.T) {
	orm := &ORMImpl{
		dialect: &MockDialectForTest{},
	}

	qb := &QueryBuilderImpl{
		orm:     orm,
		table:   "users",
		fields:  []string{"name", "age"},
		where:   []WhereCondition{{Field: "age", Operator: ">", Value: 25}},
		orderBy: []OrderBy{{Field: "name", Direction: "ASC"}},
		groupBy: []string{"age"},
		having:  "COUNT(*) > 5",
		joins:   []Join{{Type: "INNER", Table: "posts", Condition: "users.id = posts.user_id"}},
		limit:   10,
		offset:  20,
		args:    []interface{}{25},
	}

	sql := qb.buildQuery()
	expected := "SELECT name, age FROM users INNER JOIN posts ON users.id = posts.user_id WHERE age > ? GROUP BY age HAVING COUNT(*) > 5 ORDER BY name ASC LIMIT 10 OFFSET 20"
	if sql != expected {
		t.Errorf("Expected SQL '%s', got '%s'", expected, sql)
	}
}

func TestQueryBuilderImpl_ErrorHandling(t *testing.T) {
	orm := &ORMImpl{
		dialect: &MockDialectForTest{},
	}

	qb := &QueryBuilderImpl{
		orm:     orm,
		table:   "users",
		fields:  []string{"*"},
		where:   make([]WhereCondition, 0),
		orderBy: make([]OrderBy, 0),
		joins:   make([]Join, 0),
		args:    make([]interface{}, 0),
		err:     fmt.Errorf("test error"),
	}

	// Test that methods return a non-nil QueryBuilder and preserve the error
	result := qb.Where("age", ">", 25)
	if result == nil {
		t.Error("Should return a non-nil QueryBuilder even with error")
	}
	if result.(*QueryBuilderImpl).err == nil {
		t.Error("Error should be preserved in QueryBuilder")
	}

	result = qb.Select("name")
	if result == nil {
		t.Error("Should return a non-nil QueryBuilder even with error")
	}
	if result.(*QueryBuilderImpl).err == nil {
		t.Error("Error should be preserved in QueryBuilder")
	}
}
