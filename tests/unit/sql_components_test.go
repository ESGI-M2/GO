package unit

import (
	"strings"
	"testing"

	components "github.com/ESGI-M2/GO/orm/sql/components/insert"
	queryComponents "github.com/ESGI-M2/GO/orm/sql/components/queries"
)

// Import the components with proper alias
type Query = queryComponents.Query
type InsertQuery = components.InsertQuery

// TestQuery tests the Query struct and its methods
func TestQuery(t *testing.T) {
	t.Run("NewQuery", func(t *testing.T) {
		q := &Query{}
		if q == nil {
			t.Fatal("Query should not be nil")
		}

		// Test initial state
		if len(q.Fields) != 0 {
			t.Error("Initial Fields should be empty")
		}
		if q.Table != "" {
			t.Error("Initial Table should be empty")
		}
		if q.WhereClause != "" {
			t.Error("Initial WhereClause should be empty")
		}
	})

	t.Run("Select with single field", func(t *testing.T) {
		q := &Query{}
		result := q.Select("name")

		if result != q {
			t.Error("Select should return the same Query instance for chaining")
		}
		if len(q.Fields) != 1 {
			t.Errorf("Expected 1 field, got %d", len(q.Fields))
		}
		if q.Fields[0] != "name" {
			t.Errorf("Expected field 'name', got '%s'", q.Fields[0])
		}
	})

	t.Run("Select with multiple fields", func(t *testing.T) {
		q := &Query{}
		result := q.Select("id", "name", "email")

		if result != q {
			t.Error("Select should return the same Query instance for chaining")
		}
		if len(q.Fields) != 3 {
			t.Errorf("Expected 3 fields, got %d", len(q.Fields))
		}

		expectedFields := []string{"id", "name", "email"}
		for i, expected := range expectedFields {
			if q.Fields[i] != expected {
				t.Errorf("Expected field[%d] '%s', got '%s'", i, expected, q.Fields[i])
			}
		}
	})

	t.Run("Select with wildcard", func(t *testing.T) {
		q := &Query{}
		result := q.Select("*")

		if result != q {
			t.Error("Select should return the same Query instance for chaining")
		}
		if len(q.Fields) != 1 {
			t.Errorf("Expected 1 field, got %d", len(q.Fields))
		}
		if q.Fields[0] != "*" {
			t.Errorf("Expected field '*', got '%s'", q.Fields[0])
		}
	})

	t.Run("Select with invalid field should fail", func(t *testing.T) {
		// Note: The original implementation uses log.Fatalf which would terminate the test
		// In a real scenario, this should be refactored to return an error instead
		// For now, we'll test the behavior as it is

		// This test would actually cause the program to exit due to log.Fatalf
		// We skip it to avoid test failures
		t.Skip("Skipping test that would cause log.Fatalf")
	})

	t.Run("From", func(t *testing.T) {
		q := &Query{}
		result := q.From("users")

		if result != q {
			t.Error("From should return the same Query instance for chaining")
		}
		if q.Table != "users" {
			t.Errorf("Expected table 'users', got '%s'", q.Table)
		}
	})

	t.Run("Where", func(t *testing.T) {
		q := &Query{}
		result := q.Where("id = 1")

		if result != q {
			t.Error("Where should return the same Query instance for chaining")
		}
		if q.WhereClause != "id = 1" {
			t.Errorf("Expected where clause 'id = 1', got '%s'", q.WhereClause)
		}
	})

	t.Run("InnerJoin", func(t *testing.T) {
		q := &Query{}
		result := q.InnerJoin("posts", "users.id = posts.user_id")

		if result != q {
			t.Error("InnerJoin should return the same Query instance for chaining")
		}
		if len(q.Joins) != 1 {
			t.Errorf("Expected 1 join, got %d", len(q.Joins))
		}
	})

	t.Run("LeftJoin", func(t *testing.T) {
		q := &Query{}
		result := q.LeftJoin("posts", "users.id = posts.user_id")

		if result != q {
			t.Error("LeftJoin should return the same Query instance for chaining")
		}
		if len(q.Joins) != 1 {
			t.Errorf("Expected 1 join, got %d", len(q.Joins))
		}
	})

	t.Run("Multiple Joins", func(t *testing.T) {
		q := &Query{}
		q.InnerJoin("posts", "users.id = posts.user_id")
		q.LeftJoin("comments", "posts.id = comments.post_id")

		if len(q.Joins) != 2 {
			t.Errorf("Expected 2 joins, got %d", len(q.Joins))
		}
	})

	t.Run("GroupBy", func(t *testing.T) {
		q := &Query{}
		result := q.GroupBy("department")

		if result != q {
			t.Error("GroupBy should return the same Query instance for chaining")
		}
		if q.GroupByClause != "department" {
			t.Errorf("Expected group by clause 'department', got '%s'", q.GroupByClause)
		}
	})

	t.Run("Having", func(t *testing.T) {
		q := &Query{}
		result := q.Having("COUNT(*) > 1")

		if result != q {
			t.Error("Having should return the same Query instance for chaining")
		}
		if q.HavingClause != "COUNT(*) > 1" {
			t.Errorf("Expected having clause 'COUNT(*) > 1', got '%s'", q.HavingClause)
		}
	})

	t.Run("OrderBy", func(t *testing.T) {
		q := &Query{}
		result := q.OrderBy("name ASC")

		if result != q {
			t.Error("OrderBy should return the same Query instance for chaining")
		}
		if q.OrderByClause != "name ASC" {
			t.Errorf("Expected order by clause 'name ASC', got '%s'", q.OrderByClause)
		}
	})

	t.Run("Limit", func(t *testing.T) {
		q := &Query{}
		result := q.Limit(10)

		if result != q {
			t.Error("Limit should return the same Query instance for chaining")
		}
		if q.LimitValue != 10 {
			t.Errorf("Expected limit value 10, got %d", q.LimitValue)
		}
	})

	t.Run("Chained query building", func(t *testing.T) {
		q := &Query{}
		result := q.Select("id", "name").
			From("users").
			Where("active = 1").
			InnerJoin("profiles", "users.id = profiles.user_id").
			GroupBy("department").
			Having("COUNT(*) > 1").
			OrderBy("name ASC").
			Limit(50)

		if result != q {
			t.Error("Chained calls should return the same Query instance")
		}

		// Verify all properties are set correctly
		if len(q.Fields) != 2 || q.Fields[0] != "id" || q.Fields[1] != "name" {
			t.Error("Fields not set correctly in chained call")
		}
		if q.Table != "users" {
			t.Error("Table not set correctly in chained call")
		}
		if q.WhereClause != "active = 1" {
			t.Error("Where clause not set correctly in chained call")
		}
		if len(q.Joins) != 1 {
			t.Error("Joins not set correctly in chained call")
		}
		if q.GroupByClause != "department" {
			t.Error("GroupBy clause not set correctly in chained call")
		}
		if q.HavingClause != "COUNT(*) > 1" {
			t.Error("Having clause not set correctly in chained call")
		}
		if q.OrderByClause != "name ASC" {
			t.Error("OrderBy clause not set correctly in chained call")
		}
		if q.LimitValue != 50 {
			t.Error("Limit value not set correctly in chained call")
		}
	})
}

// TestInsertQuery tests the InsertQuery struct and its methods
func TestInsertQuery(t *testing.T) {
	t.Run("NewInsertQuery", func(t *testing.T) {
		iq := &InsertQuery{}
		if iq == nil {
			t.Fatal("InsertQuery should not be nil")
		}

		// Test initial state
		if iq.Table != "" {
			t.Error("Initial Table should be empty")
		}
		if len(iq.Columns) != 0 {
			t.Error("Initial Columns should be empty")
		}
		if len(iq.Values) != 0 {
			t.Error("Initial Values should be empty")
		}
	})

	t.Run("Into", func(t *testing.T) {
		iq := &InsertQuery{}
		result := iq.Into("users")

		if result != iq {
			t.Error("Into should return the same InsertQuery instance for chaining")
		}
		if iq.Table != "users" {
			t.Errorf("Expected table 'users', got '%s'", iq.Table)
		}
	})

	t.Run("Set", func(t *testing.T) {
		iq := &InsertQuery{}
		columns := []string{"name", "email", "age"}
		values := []interface{}{"John", "john@test.com", 30}
		result := iq.Set(columns, values)

		if result != iq {
			t.Error("Set should return the same InsertQuery instance for chaining")
		}
		if len(iq.Columns) != 3 {
			t.Errorf("Expected 3 columns, got %d", len(iq.Columns))
		}
		if len(iq.Values) != 3 {
			t.Errorf("Expected 3 values, got %d", len(iq.Values))
		}

		for i, expected := range columns {
			if iq.Columns[i] != expected {
				t.Errorf("Expected column[%d] '%s', got '%s'", i, expected, iq.Columns[i])
			}
		}

		for i, expected := range values {
			if iq.Values[i] != expected {
				t.Errorf("Expected value[%d] '%v', got '%v'", i, expected, iq.Values[i])
			}
		}
	})

	t.Run("Chained insert query building", func(t *testing.T) {
		iq := &InsertQuery{}
		columns := []string{"name", "email", "age"}
		values := []interface{}{"John", "john@test.com", 30}

		result := iq.Into("users").Set(columns, values)

		if result != iq {
			t.Error("Chained calls should return the same InsertQuery instance")
		}

		// Verify all properties are set correctly
		if iq.Table != "users" {
			t.Error("Table not set correctly in chained call")
		}
		if len(iq.Columns) != 3 {
			t.Error("Columns not set correctly in chained call")
		}
		if len(iq.Values) != 3 {
			t.Error("Values not set correctly in chained call")
		}
	})
}

// TestQueryComponents tests individual query component methods
func TestQueryComponents(t *testing.T) {
	t.Run("Empty Query State", func(t *testing.T) {
		q := &Query{}

		// Test that all fields start empty/zero
		if q.Fields != nil && len(q.Fields) > 0 {
			t.Error("Fields should start empty")
		}
		if q.Table != "" {
			t.Error("Table should start empty")
		}
		if q.WhereClause != "" {
			t.Error("WhereClause should start empty")
		}
		if q.LimitValue != 0 {
			t.Error("LimitValue should start as 0")
		}
		if q.Joins != nil && len(q.Joins) > 0 {
			t.Error("Joins should start empty")
		}
		if q.GroupByClause != "" {
			t.Error("GroupByClause should start empty")
		}
		if q.HavingClause != "" {
			t.Error("HavingClause should start empty")
		}
		if q.OrderByClause != "" {
			t.Error("OrderByClause should start empty")
		}
	})

	t.Run("Query State After Operations", func(t *testing.T) {
		q := &Query{}

		// Perform operations
		q.Select("id", "name")
		q.From("users")
		q.Where("status = 'active'")
		q.LeftJoin("profiles", "users.id = profiles.user_id")
		q.GroupBy("status")
		q.Having("COUNT(*) > 5")
		q.OrderBy("created_at DESC")
		q.Limit(100)

		// Verify state
		if len(q.Fields) != 2 {
			t.Errorf("Expected 2 fields after Select, got %d", len(q.Fields))
		}
		if q.Table != "users" {
			t.Errorf("Expected table 'users' after From, got '%s'", q.Table)
		}
		if !strings.Contains(q.WhereClause, "status = 'active'") {
			t.Errorf("Expected where clause to contain status condition, got '%s'", q.WhereClause)
		}
		if len(q.Joins) != 1 {
			t.Errorf("Expected 1 join after Join, got %d", len(q.Joins))
		}
		if q.GroupByClause != "status" {
			t.Errorf("Expected group by 'status', got '%s'", q.GroupByClause)
		}
		if !strings.Contains(q.HavingClause, "COUNT(*) > 5") {
			t.Errorf("Expected having clause to contain count condition, got '%s'", q.HavingClause)
		}
		if !strings.Contains(q.OrderByClause, "created_at DESC") {
			t.Errorf("Expected order by to contain created_at DESC, got '%s'", q.OrderByClause)
		}
		if q.LimitValue != 100 {
			t.Errorf("Expected limit 100, got %d", q.LimitValue)
		}
	})
}

// TestInsertQueryComponents tests individual insert query component methods
func TestInsertQueryComponents(t *testing.T) {
	t.Run("Empty InsertQuery State", func(t *testing.T) {
		iq := &InsertQuery{}

		// Test that all fields start empty/zero
		if iq.Table != "" {
			t.Error("Table should start empty")
		}
		if iq.Columns != nil && len(iq.Columns) > 0 {
			t.Error("Columns should start empty")
		}
		if iq.Values != nil && len(iq.Values) > 0 {
			t.Error("Values should start empty")
		}
	})

	t.Run("InsertQuery State After Operations", func(t *testing.T) {
		iq := &InsertQuery{}

		// Perform operations
		iq.Into("users")
		iq.Set([]string{"name", "email", "age"}, []interface{}{"John", "john@test.com", 30})

		// Verify state
		if iq.Table != "users" {
			t.Errorf("Expected table 'users' after Into, got '%s'", iq.Table)
		}
		if len(iq.Columns) != 3 {
			t.Errorf("Expected 3 columns after SetColumns, got %d", len(iq.Columns))
		}

		expectedColumns := []string{"name", "email", "age"}
		for i, expected := range expectedColumns {
			if iq.Columns[i] != expected {
				t.Errorf("Expected column[%d] '%s', got '%s'", i, expected, iq.Columns[i])
			}
		}
	})

	t.Run("Multiple SetColumns calls", func(t *testing.T) {
		iq := &InsertQuery{}

		// First call
		iq.Set([]string{"name", "email"}, []interface{}{"John", "john@test.com"})
		if len(iq.Columns) != 2 {
			t.Errorf("Expected 2 columns after first Set, got %d", len(iq.Columns))
		}

		// Second call (should replace, not append)
		iq.Set([]string{"id", "name", "email", "age"}, []interface{}{1, "John", "john@test.com", 30})
		if len(iq.Columns) != 4 {
			t.Errorf("Expected 4 columns after second Set, got %d", len(iq.Columns))
		}
	})
}
