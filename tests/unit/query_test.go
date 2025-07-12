package unit

import (
	"testing"

	"github.com/ESGI-M2/GO/orm/builder"
	"github.com/ESGI-M2/GO/orm/core/interfaces"
	"github.com/ESGI-M2/GO/orm/factory"
)

type QueryTestModel struct {
	ID    int    `orm:"pk,auto"`
	Name  string `orm:"column:name"`
	Email string `orm:"column:email,unique"`
}

func setupQueryBuilder() interfaces.QueryBuilder {
	orm := builder.NewSimpleORM().
		WithDialect(factory.Mock).
		RegisterModel(&QueryTestModel{})

	// Connect first
	err := orm.Connect()
	if err != nil {
		panic(err)
	}

	return orm.Query(&QueryTestModel{})
}

func TestQueryBuilder_Select(t *testing.T) {
	qb := setupQueryBuilder()
	result := qb.Select("id", "name")
	if result == nil {
		t.Error("Select should return a query builder")
	}
}

func TestQueryBuilder_From(t *testing.T) {
	qb := setupQueryBuilder()
	result := qb.From("test_table")
	if result == nil {
		t.Error("From should return a query builder")
	}
}

func TestQueryBuilder_Where(t *testing.T) {
	qb := setupQueryBuilder()
	result := qb.Where("name", "=", "test")
	if result == nil {
		t.Error("Where should return a query builder")
	}
}

func TestQueryBuilder_WhereIn(t *testing.T) {
	qb := setupQueryBuilder()
	values := []interface{}{"test1", "test2"}
	result := qb.WhereIn("name", values)
	if result == nil {
		t.Error("WhereIn should return a query builder")
	}
}

func TestQueryBuilder_WhereNotIn(t *testing.T) {
	qb := setupQueryBuilder()
	values := []interface{}{"test1", "test2"}
	result := qb.WhereNotIn("name", values)
	if result == nil {
		t.Error("WhereNotIn should return a query builder")
	}
}

func TestQueryBuilder_OrderBy(t *testing.T) {
	qb := setupQueryBuilder()
	result := qb.OrderBy("name", "ASC")
	if result == nil {
		t.Error("OrderBy should return a query builder")
	}
}

func TestQueryBuilder_GroupBy(t *testing.T) {
	qb := setupQueryBuilder()
	result := qb.GroupBy("name", "email")
	if result == nil {
		t.Error("GroupBy should return a query builder")
	}
}

func TestQueryBuilder_Having(t *testing.T) {
	qb := setupQueryBuilder()
	result := qb.Having("COUNT(*) > 1")
	if result == nil {
		t.Error("Having should return a query builder")
	}
}

func TestQueryBuilder_Limit(t *testing.T) {
	qb := setupQueryBuilder()
	result := qb.Limit(10)
	if result == nil {
		t.Error("Limit should return a query builder")
	}
}

func TestQueryBuilder_Offset(t *testing.T) {
	qb := setupQueryBuilder()
	result := qb.Offset(5)
	if result == nil {
		t.Error("Offset should return a query builder")
	}
}

func TestQueryBuilder_Join(t *testing.T) {
	qb := setupQueryBuilder()
	result := qb.Join("users", "users.id = posts.user_id")
	if result == nil {
		t.Error("Join should return a query builder")
	}
}

func TestQueryBuilder_LeftJoin(t *testing.T) {
	qb := setupQueryBuilder()
	result := qb.LeftJoin("users", "users.id = posts.user_id")
	if result == nil {
		t.Error("LeftJoin should return a query builder")
	}
}

func TestQueryBuilder_RightJoin(t *testing.T) {
	qb := setupQueryBuilder()
	result := qb.RightJoin("users", "users.id = posts.user_id")
	if result == nil {
		t.Error("RightJoin should return a query builder")
	}
}

func TestQueryBuilder_InnerJoin(t *testing.T) {
	qb := setupQueryBuilder()
	result := qb.InnerJoin("users", "users.id = posts.user_id")
	if result == nil {
		t.Error("InnerJoin should return a query builder")
	}
}

func TestQueryBuilder_GetSQL_GetArgs(t *testing.T) {
	qb := setupQueryBuilder()
	qb.Where("name", "=", "test")

	sql := qb.GetSQL()
	if sql == "" {
		t.Error("GetSQL should return a non-empty string")
	}

	args := qb.GetArgs()
	if args == nil {
		t.Error("GetArgs should return a slice")
	}
}

func TestQueryBuilder_Find(t *testing.T) {
	qb := setupQueryBuilder()
	results, err := qb.Find()
	if err != nil {
		t.Errorf("Find failed: %v", err)
	}
	if results == nil {
		t.Error("Find should return a slice, even if empty")
	}
}

func TestQueryBuilder_FindOne(t *testing.T) {
	qb := setupQueryBuilder()
	_, err := qb.FindOne()
	if err != nil {
		t.Errorf("FindOne failed: %v", err)
	}
	// result can be nil for empty results
}

func TestQueryBuilder_Count(t *testing.T) {
	qb := setupQueryBuilder()
	count, err := qb.Count()
	if err != nil {
		t.Errorf("Count failed: %v", err)
	}
	if count < 0 {
		t.Error("Count should be non-negative")
	}
}

func TestQueryBuilder_Exists(t *testing.T) {
	qb := setupQueryBuilder()
	_, err := qb.Exists()
	if err != nil {
		t.Errorf("Exists failed: %v", err)
	}
	// exists can be true or false
}

func TestQueryBuilder_Raw(t *testing.T) {
	orm := builder.NewSimpleORM().WithDialect(factory.Mock)

	// Connect first
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	raw := orm.Raw("SELECT * FROM users WHERE id = ?", 1)
	if raw == nil {
		t.Error("Raw should return a query builder")
	}

	results, err := raw.Find()
	if err != nil {
		t.Errorf("Raw Find failed: %v", err)
	}
	if results == nil {
		t.Error("Raw Find should return a slice, even if empty")
	}
}

func TestQueryBuilder_ErrorCases(t *testing.T) {
	// Test with error query builder
	orm := builder.NewSimpleORM().WithDialect(factory.Mock)
	// Don't register model to trigger error case
	qb := orm.Query(&QueryTestModel{})

	_, err := qb.Find()
	if err == nil {
		t.Error("Expected error when ORM is not connected")
	}
}
