package unit

import (
	"testing"

	"github.com/ESGI-M2/GO/orm/builder"
	"github.com/ESGI-M2/GO/orm/core/interfaces"
	"github.com/ESGI-M2/GO/orm/factory"
)

type AdvancedQueryTestModel struct {
	ID          int    `orm:"pk,auto"`
	Name        string `orm:"column:name"`
	Email       string `orm:"column:email,unique"`
	Age         int    `orm:"column:age"`
	IsActive    bool   `orm:"column:is_active"`
	Description string `orm:"column:description"`
}

func setupAdvancedQueryBuilder() interfaces.QueryBuilder {
	orm := builder.NewSimpleORM().
		WithDialect(factory.Mock).
		RegisterModel(&AdvancedQueryTestModel{})

	err := orm.Connect()
	if err != nil {
		panic(err)
	}

	return orm.Query(&AdvancedQueryTestModel{})
}

func TestAdvancedQueryBuilder_WhereOr(t *testing.T) {
	qb := setupAdvancedQueryBuilder()

	conditions := []interfaces.WhereCondition{
		{Field: "name", Operator: "=", Value: "John"},
		{Field: "age", Operator: ">", Value: 18},
	}

	result := qb.WhereOr(conditions...)
	if result == nil {
		t.Error("WhereOr should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("WhereOr should generate SQL")
	}
}

func TestAdvancedQueryBuilder_WhereRaw(t *testing.T) {
	qb := setupAdvancedQueryBuilder()

	result := qb.WhereRaw("name = ? AND age > ?", "John", 18)
	if result == nil {
		t.Error("WhereRaw should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("WhereRaw should generate SQL")
	}

	// Test args
	args := result.GetArgs()
	if len(args) != 2 {
		t.Errorf("WhereRaw should have 2 args, got %d", len(args))
	}
}

func TestAdvancedQueryBuilder_WhereBetween(t *testing.T) {
	qb := setupAdvancedQueryBuilder()

	result := qb.WhereBetween("age", 18, 65)
	if result == nil {
		t.Error("WhereBetween should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("WhereBetween should generate SQL")
	}
}

func TestAdvancedQueryBuilder_WhereNotBetween(t *testing.T) {
	qb := setupAdvancedQueryBuilder()

	result := qb.WhereNotBetween("age", 18, 65)
	if result == nil {
		t.Error("WhereNotBetween should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("WhereNotBetween should generate SQL")
	}
}

func TestAdvancedQueryBuilder_WhereNull(t *testing.T) {
	qb := setupAdvancedQueryBuilder()

	result := qb.WhereNull("description")
	if result == nil {
		t.Error("WhereNull should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("WhereNull should generate SQL")
	}
}

func TestAdvancedQueryBuilder_WhereNotNull(t *testing.T) {
	qb := setupAdvancedQueryBuilder()

	result := qb.WhereNotNull("description")
	if result == nil {
		t.Error("WhereNotNull should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("WhereNotNull should generate SQL")
	}
}

func TestAdvancedQueryBuilder_WhereLike(t *testing.T) {
	qb := setupAdvancedQueryBuilder()

	result := qb.WhereLike("name", "%John%")
	if result == nil {
		t.Error("WhereLike should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("WhereLike should generate SQL")
	}
}

func TestAdvancedQueryBuilder_WhereNotLike(t *testing.T) {
	qb := setupAdvancedQueryBuilder()

	result := qb.WhereNotLike("name", "%John%")
	if result == nil {
		t.Error("WhereNotLike should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("WhereNotLike should generate SQL")
	}
}

func TestAdvancedQueryBuilder_WhereRegexp(t *testing.T) {
	qb := setupAdvancedQueryBuilder()

	result := qb.WhereRegexp("email", "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
	if result == nil {
		t.Error("WhereRegexp should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("WhereRegexp should generate SQL")
	}
}

func TestAdvancedQueryBuilder_WhereNotRegexp(t *testing.T) {
	qb := setupAdvancedQueryBuilder()

	result := qb.WhereNotRegexp("email", "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
	if result == nil {
		t.Error("WhereNotRegexp should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("WhereNotRegexp should generate SQL")
	}
}

func TestAdvancedQueryBuilder_FullTextSearch(t *testing.T) {
	qb := setupAdvancedQueryBuilder()

	result := qb.FullTextSearch([]string{"name", "description"}, "search query")
	if result == nil {
		t.Error("FullTextSearch should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("FullTextSearch should generate SQL")
	}
}

func TestAdvancedQueryBuilder_SubQuery(t *testing.T) {
	qb := setupAdvancedQueryBuilder()

	result := qb.SubQuery("sub", func(subQb interfaces.QueryBuilder) interfaces.QueryBuilder {
		return subQb.Select("id").From("users").Where("active", "=", true)
	})
	if result == nil {
		t.Error("SubQuery should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("SubQuery should generate SQL")
	}
}

func TestAdvancedQueryBuilder_With(t *testing.T) {
	qb := setupAdvancedQueryBuilder()

	result := qb.With("profile", func(withQb interfaces.QueryBuilder) interfaces.QueryBuilder {
		return withQb.Select("*").From("user_profiles").Where("user_id", "=", "users.id")
	})
	if result == nil {
		t.Error("With should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("With should generate SQL")
	}
}

func TestAdvancedQueryBuilder_WithCount(t *testing.T) {
	qb := setupAdvancedQueryBuilder()

	result := qb.WithCount("posts")
	if result == nil {
		t.Error("WithCount should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("WithCount should generate SQL")
	}
}

func TestAdvancedQueryBuilder_WithExists(t *testing.T) {
	qb := setupAdvancedQueryBuilder()

	result := qb.WithExists("posts", func(existsQb interfaces.QueryBuilder) interfaces.QueryBuilder {
		return existsQb.Select("1").From("posts").Where("user_id", "=", "users.id")
	})
	if result == nil {
		t.Error("WithExists should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("WithExists should generate SQL")
	}
}

func TestAdvancedQueryBuilder_CursorPaginate(t *testing.T) {
	qb := setupAdvancedQueryBuilder()

	result := qb.CursorPaginate("id", 10, 20)
	if result == nil {
		t.Error("CursorPaginate should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("CursorPaginate should generate SQL")
	}
}

func TestAdvancedQueryBuilder_OffsetPaginate(t *testing.T) {
	qb := setupAdvancedQueryBuilder()

	result := qb.OffsetPaginate(1, 10)
	if result == nil {
		t.Error("OffsetPaginate should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("OffsetPaginate should generate SQL")
	}
}

func TestAdvancedQueryBuilder_ForUpdate(t *testing.T) {
	qb := setupAdvancedQueryBuilder()

	result := qb.ForUpdate()
	if result == nil {
		t.Error("ForUpdate should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("ForUpdate should generate SQL")
	}
}

func TestAdvancedQueryBuilder_ForShare(t *testing.T) {
	qb := setupAdvancedQueryBuilder()

	result := qb.ForShare()
	if result == nil {
		t.Error("ForShare should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("ForShare should generate SQL")
	}
}

func TestAdvancedQueryBuilder_Distinct(t *testing.T) {
	qb := setupAdvancedQueryBuilder()

	result := qb.Distinct()
	if result == nil {
		t.Error("Distinct should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("Distinct should generate SQL")
	}
}

func TestAdvancedQueryBuilder_Union(t *testing.T) {
	qb := setupAdvancedQueryBuilder()
	otherQb := setupAdvancedQueryBuilder()

	result := qb.Union(otherQb)
	if result == nil {
		t.Error("Union should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("Union should generate SQL")
	}
}

func TestAdvancedQueryBuilder_UnionAll(t *testing.T) {
	qb := setupAdvancedQueryBuilder()
	otherQb := setupAdvancedQueryBuilder()

	result := qb.UnionAll(otherQb)
	if result == nil {
		t.Error("UnionAll should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("UnionAll should generate SQL")
	}
}

func TestAdvancedQueryBuilder_Lock(t *testing.T) {
	qb := setupAdvancedQueryBuilder()

	result := qb.Lock("FOR UPDATE")
	if result == nil {
		t.Error("Lock should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("Lock should generate SQL")
	}
}

func TestAdvancedQueryBuilder_Cache(t *testing.T) {
	qb := setupAdvancedQueryBuilder()

	result := qb.Cache(300) // 5 minutes
	if result == nil {
		t.Error("Cache should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("Cache should generate SQL")
	}
}

func TestAdvancedQueryBuilder_WithoutCache(t *testing.T) {
	qb := setupAdvancedQueryBuilder()

	result := qb.WithoutCache()
	if result == nil {
		t.Error("WithoutCache should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("WithoutCache should generate SQL")
	}
}

func TestAdvancedQueryBuilder_ChainedOperations(t *testing.T) {
	qb := setupAdvancedQueryBuilder()

	result := qb.
		Select("id", "name", "email").
		From("users").
		Where("age", ">", 18).
		WhereNotNull("email").
		WhereLike("name", "%John%").
		OrderBy("name", "ASC").
		Distinct().
		Limit(10).
		Offset(5)

	if result == nil {
		t.Error("Chained operations should return a query builder")
	}

	// Test SQL generation
	sql := result.GetSQL()
	if sql == "" {
		t.Error("Chained operations should generate SQL")
	}

	// Test args
	args := result.GetArgs()
	if len(args) == 0 {
		t.Error("Chained operations should have args")
	}
}

func TestAdvancedQueryBuilder_ComplexQuery(t *testing.T) {
	qb := setupAdvancedQueryBuilder()

	result := qb.
		Select("u.id", "u.name", "u.email").
		From("users u").
		LeftJoin("user_profiles p", "p.user_id = u.id").
		Where("u.age", ">", 18).
		WhereOr(
			interfaces.WhereCondition{Field: "u.is_active", Operator: "=", Value: true},
			interfaces.WhereCondition{Field: "p.verified", Operator: "=", Value: true},
		).
		WhereBetween("u.age", 18, 65).
		WhereNotNull("u.email").
		GroupBy("u.id").
		Having("COUNT(p.id) > 0").
		OrderBy("u.name", "ASC").
		Limit(20)

	if result == nil {
		t.Error("Complex query should return a query builder")
	}

	// Test execution
	_, err := result.Find()
	if err != nil {
		t.Errorf("Complex query execution failed: %v", err)
	}
}

func TestAdvancedQueryBuilder_ErrorCases(t *testing.T) {
	// Test with error query builder
	orm := builder.NewSimpleORM().WithDialect(factory.Mock)
	qb := orm.Query(&AdvancedQueryTestModel{})

	// Test WhereOr with disconnected ORM
	result := qb.WhereOr(interfaces.WhereCondition{Field: "name", Operator: "=", Value: "test"})
	_, err := result.Find()
	if err == nil {
		t.Error("Expected error when ORM is not connected")
	}

	// Test WhereRaw with disconnected ORM
	result = qb.WhereRaw("name = ?", "test")
	_, err = result.Find()
	if err == nil {
		t.Error("Expected error when ORM is not connected")
	}
}
