package query

import (
	"fmt"
	"strings"

	"github.com/ESGI-M2/GO/orm/core/interfaces"
)

// BuilderImpl implements the QueryBuilder interface
type BuilderImpl struct {
	Orm      interfaces.ORM
	Metadata *interfaces.ModelMetadata
	Err      error

	// Query components
	table      string
	fields     []string
	where      []interfaces.WhereCondition
	orderBy    []interfaces.OrderBy
	groupBy    []string
	having     string
	havingArgs []interface{}
	joins      []interfaces.Join
	limit      int
	offset     int
	args       []interface{}

	// Raw SQL
	rawSQL  string
	rawArgs []interface{}

	// New advanced features
	distinct      bool
	lockType      string
	withRelations map[string]func(interfaces.QueryBuilder) interfaces.QueryBuilder
	withCounts    []string
	withExists    map[string]func(interfaces.QueryBuilder) interfaces.QueryBuilder
	cacheTTL      int
	useCache      bool
	subQueries    map[string]interfaces.QueryBuilder
	unions        []interfaces.QueryBuilder
	unionAlls     []interfaces.QueryBuilder
	cursorField   string
	cursorValue   interface{}
	page          int
	perPage       int
}

// NewBuilder creates a new query builder
func NewBuilder(orm interfaces.ORM, metadata *interfaces.ModelMetadata) *BuilderImpl {
	return &BuilderImpl{
		Orm:           orm,
		Metadata:      metadata,
		table:         metadata.TableName,
		fields:        []string{"*"},
		where:         make([]interfaces.WhereCondition, 0),
		orderBy:       make([]interfaces.OrderBy, 0),
		joins:         make([]interfaces.Join, 0),
		limit:         0,
		offset:        0,
		args:          make([]interface{}, 0),
		withRelations: make(map[string]func(interfaces.QueryBuilder) interfaces.QueryBuilder),
		withExists:    make(map[string]func(interfaces.QueryBuilder) interfaces.QueryBuilder),
		subQueries:    make(map[string]interfaces.QueryBuilder),
		unions:        make([]interfaces.QueryBuilder, 0),
		unionAlls:     make([]interfaces.QueryBuilder, 0),
	}
}

// NewRawBuilder creates a new raw SQL query builder
func NewRawBuilder(orm interfaces.ORM, sql string, args ...interface{}) *BuilderImpl {
	return &BuilderImpl{
		Orm:      orm,
		rawSQL:   sql,
		rawArgs:  args,
		Metadata: nil,
	}
}

// Select sets the fields to select
func (qb *BuilderImpl) Select(fields ...string) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	if len(fields) == 0 {
		qb.fields = []string{"*"}
	} else {
		qb.fields = fields
	}
	return qb
}

// From sets the table name
func (qb *BuilderImpl) From(table string) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}
	qb.table = table
	return qb
}

// Where adds a WHERE condition
func (qb *BuilderImpl) Where(field, operator string, value interface{}) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	qb.where = append(qb.where, interfaces.WhereCondition{
		Field:    field,
		Operator: operator,
		Value:    value,
		Logical:  "AND",
	})

	if value != nil {
		qb.args = append(qb.args, value)
	}

	return qb
}

// WhereIn adds a WHERE IN condition
func (qb *BuilderImpl) WhereIn(field string, values []interface{}) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	if len(values) == 0 {
		return qb
	}

	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = qb.Orm.GetDialect().GetPlaceholder(i)
		qb.args = append(qb.args, values[i])
	}

	condition := fmt.Sprintf("%s IN (%s)", field, strings.Join(placeholders, ", "))
	qb.where = append(qb.where, interfaces.WhereCondition{
		Field:    condition,
		Operator: "",
		Value:    nil,
		Logical:  "AND",
	})

	return qb
}

// WhereNotIn adds a WHERE NOT IN condition
func (qb *BuilderImpl) WhereNotIn(field string, values []interface{}) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	if len(values) == 0 {
		return qb
	}

	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = qb.Orm.GetDialect().GetPlaceholder(i)
		qb.args = append(qb.args, values[i])
	}

	condition := fmt.Sprintf("%s NOT IN (%s)", field, strings.Join(placeholders, ", "))
	qb.where = append(qb.where, interfaces.WhereCondition{
		Field:    condition,
		Operator: "",
		Value:    nil,
		Logical:  "AND",
	})

	return qb
}

// WhereOr adds OR conditions
func (qb *BuilderImpl) WhereOr(conditions ...interfaces.WhereCondition) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	if len(conditions) == 0 {
		return qb
	}

	// Group OR conditions
	var orConditions []string
	for _, condition := range conditions {
		if condition.Field != "" {
			if condition.Operator != "" && condition.Value != nil {
				orConditions = append(orConditions, fmt.Sprintf("%s %s %s", condition.Field, condition.Operator, qb.Orm.GetDialect().GetPlaceholder(len(qb.args))))
				qb.args = append(qb.args, condition.Value)
			} else if condition.Field != "" {
				orConditions = append(orConditions, condition.Field)
			}
		}
	}

	if len(orConditions) > 0 {
		condition := fmt.Sprintf("(%s)", strings.Join(orConditions, " OR "))
		qb.where = append(qb.where, interfaces.WhereCondition{
			Field:    condition,
			Operator: "",
			Value:    nil,
			Logical:  "AND",
		})
	}

	return qb
}

// WhereRaw adds a raw WHERE condition
func (qb *BuilderImpl) WhereRaw(condition string, args ...interface{}) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	qb.where = append(qb.where, interfaces.WhereCondition{
		Field:    condition,
		Operator: "",
		Value:    nil,
		Logical:  "AND",
		Raw:      true,
	})

	qb.args = append(qb.args, args...)
	return qb
}

// WhereBetween adds a WHERE BETWEEN condition
func (qb *BuilderImpl) WhereBetween(field string, min, max interface{}) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	qb.where = append(qb.where, interfaces.WhereCondition{
		Field:    fmt.Sprintf("%s BETWEEN %s AND %s", field, qb.Orm.GetDialect().GetPlaceholder(len(qb.args)), qb.Orm.GetDialect().GetPlaceholder(len(qb.args)+1)),
		Operator: "",
		Value:    nil,
		Logical:  "AND",
	})

	qb.args = append(qb.args, min, max)
	return qb
}

// WhereNotBetween adds a WHERE NOT BETWEEN condition
func (qb *BuilderImpl) WhereNotBetween(field string, min, max interface{}) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	qb.where = append(qb.where, interfaces.WhereCondition{
		Field:    fmt.Sprintf("%s NOT BETWEEN %s AND %s", field, qb.Orm.GetDialect().GetPlaceholder(len(qb.args)), qb.Orm.GetDialect().GetPlaceholder(len(qb.args)+1)),
		Operator: "",
		Value:    nil,
		Logical:  "AND",
	})

	qb.args = append(qb.args, min, max)
	return qb
}

// WhereNull adds a WHERE IS NULL condition
func (qb *BuilderImpl) WhereNull(field string) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	qb.where = append(qb.where, interfaces.WhereCondition{
		Field:    fmt.Sprintf("%s IS NULL", field),
		Operator: "",
		Value:    nil,
		Logical:  "AND",
	})

	return qb
}

// WhereNotNull adds a WHERE IS NOT NULL condition
func (qb *BuilderImpl) WhereNotNull(field string) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	qb.where = append(qb.where, interfaces.WhereCondition{
		Field:    fmt.Sprintf("%s IS NOT NULL", field),
		Operator: "",
		Value:    nil,
		Logical:  "AND",
	})

	return qb
}

// WhereLike adds a WHERE LIKE condition
func (qb *BuilderImpl) WhereLike(field, pattern string) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	qb.where = append(qb.where, interfaces.WhereCondition{
		Field:    field,
		Operator: "LIKE",
		Value:    pattern,
		Logical:  "AND",
	})

	qb.args = append(qb.args, pattern)
	return qb
}

// WhereNotLike adds a WHERE NOT LIKE condition
func (qb *BuilderImpl) WhereNotLike(field, pattern string) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	qb.where = append(qb.where, interfaces.WhereCondition{
		Field:    field,
		Operator: "NOT LIKE",
		Value:    pattern,
		Logical:  "AND",
	})

	qb.args = append(qb.args, pattern)
	return qb
}

// WhereRegexp adds a WHERE REGEXP condition
func (qb *BuilderImpl) WhereRegexp(field, pattern string) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	qb.where = append(qb.where, interfaces.WhereCondition{
		Field:    field,
		Operator: "REGEXP",
		Value:    pattern,
		Logical:  "AND",
	})

	qb.args = append(qb.args, pattern)
	return qb
}

// WhereNotRegexp adds a WHERE NOT REGEXP condition
func (qb *BuilderImpl) WhereNotRegexp(field, pattern string) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	qb.where = append(qb.where, interfaces.WhereCondition{
		Field:    field,
		Operator: "NOT REGEXP",
		Value:    pattern,
		Logical:  "AND",
	})

	qb.args = append(qb.args, pattern)
	return qb
}

// FullTextSearch adds a full-text search condition
func (qb *BuilderImpl) FullTextSearch(fields []string, query string) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	fieldsStr := strings.Join(fields, ", ")
	condition := fmt.Sprintf("MATCH(%s) AGAINST(? IN BOOLEAN MODE)", fieldsStr)

	qb.where = append(qb.where, interfaces.WhereCondition{
		Field:    condition,
		Operator: "",
		Value:    nil,
		Logical:  "AND",
	})

	qb.args = append(qb.args, query)
	return qb
}

// SubQuery adds a subquery
func (qb *BuilderImpl) SubQuery(alias string, fn func(interfaces.QueryBuilder) interfaces.QueryBuilder) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	subQuery := NewBuilder(qb.Orm, qb.Metadata)
	subQuery = fn(subQuery).(*BuilderImpl)

	qb.subQueries[alias] = subQuery
	return qb
}

// With adds eager loading for relations
func (qb *BuilderImpl) With(relation string, fn func(interfaces.QueryBuilder) interfaces.QueryBuilder) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	qb.withRelations[relation] = fn
	return qb
}

// WithCount adds count for relations
func (qb *BuilderImpl) WithCount(relation string) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	qb.withCounts = append(qb.withCounts, relation)
	return qb
}

// WithExists adds exists condition for relations
func (qb *BuilderImpl) WithExists(relation string, fn func(interfaces.QueryBuilder) interfaces.QueryBuilder) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	qb.withExists[relation] = fn
	return qb
}

// CursorPaginate adds cursor-based pagination
func (qb *BuilderImpl) CursorPaginate(cursorField string, cursorValue interface{}, limit int) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	qb.cursorField = cursorField
	qb.cursorValue = cursorValue
	qb.limit = limit

	if cursorValue != nil {
		qb.Where(cursorField, ">", cursorValue)
	}

	return qb
}

// OffsetPaginate adds offset-based pagination
func (qb *BuilderImpl) OffsetPaginate(page, perPage int) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	qb.page = page
	qb.perPage = perPage
	qb.offset = (page - 1) * perPage
	qb.limit = perPage

	return qb
}

// ForUpdate adds FOR UPDATE lock
func (qb *BuilderImpl) ForUpdate() interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	qb.lockType = "FOR UPDATE"
	return qb
}

// ForShare adds FOR SHARE lock
func (qb *BuilderImpl) ForShare() interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	qb.lockType = "FOR SHARE"
	return qb
}

// Distinct adds DISTINCT clause
func (qb *BuilderImpl) Distinct() interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	qb.distinct = true
	return qb
}

// Union adds UNION clause
func (qb *BuilderImpl) Union(other interfaces.QueryBuilder) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	qb.unions = append(qb.unions, other)
	return qb
}

// UnionAll adds UNION ALL clause
func (qb *BuilderImpl) UnionAll(other interfaces.QueryBuilder) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	qb.unionAlls = append(qb.unionAlls, other)
	return qb
}

// Lock adds a lock clause
func (qb *BuilderImpl) Lock(lockType string) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	qb.lockType = lockType
	return qb
}

// Cache enables query caching
func (qb *BuilderImpl) Cache(ttl int) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	qb.useCache = true
	qb.cacheTTL = ttl
	return qb
}

// WithoutCache disables query caching
func (qb *BuilderImpl) WithoutCache() interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	qb.useCache = false
	return qb
}

// OrderBy adds an ORDER BY clause
func (qb *BuilderImpl) OrderBy(field, direction string) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	if direction == "" {
		direction = "ASC"
	}

	qb.orderBy = append(qb.orderBy, interfaces.OrderBy{
		Field:     field,
		Direction: strings.ToUpper(direction),
	})

	return qb
}

// GroupBy adds a GROUP BY clause
func (qb *BuilderImpl) GroupBy(fields ...string) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}
	qb.groupBy = append(qb.groupBy, fields...)
	return qb
}

// Having adds a HAVING clause
func (qb *BuilderImpl) Having(condition string, args ...interface{}) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}
	qb.having = condition
	qb.havingArgs = append(qb.havingArgs, args...)
	return qb
}

// Limit sets the LIMIT clause
func (qb *BuilderImpl) Limit(limit int) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}
	qb.limit = limit
	return qb
}

// Offset sets the OFFSET clause
func (qb *BuilderImpl) Offset(offset int) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}
	qb.offset = offset
	return qb
}

// Join adds a JOIN clause
func (qb *BuilderImpl) Join(table, condition string) interfaces.QueryBuilder {
	return qb.addJoin("INNER", table, condition)
}

// LeftJoin adds a LEFT JOIN clause
func (qb *BuilderImpl) LeftJoin(table, condition string) interfaces.QueryBuilder {
	return qb.addJoin("LEFT", table, condition)
}

// RightJoin adds a RIGHT JOIN clause
func (qb *BuilderImpl) RightJoin(table, condition string) interfaces.QueryBuilder {
	return qb.addJoin("RIGHT", table, condition)
}

// InnerJoin adds an INNER JOIN clause
func (qb *BuilderImpl) InnerJoin(table, condition string) interfaces.QueryBuilder {
	return qb.addJoin("INNER", table, condition)
}

// addJoin adds a join clause
func (qb *BuilderImpl) addJoin(joinType, table, condition string) interfaces.QueryBuilder {
	if qb.Err != nil {
		return qb
	}

	qb.joins = append(qb.joins, interfaces.Join{
		Type:      joinType,
		Table:     table,
		Condition: condition,
	})

	return qb
}

// Raw creates a raw SQL query builder
func (qb *BuilderImpl) Raw(sql string, args ...interface{}) interfaces.QueryBuilder {
	return NewRawBuilder(qb.Orm, sql, args...)
}

// GetSQL returns the generated SQL query
func (qb *BuilderImpl) GetSQL() string {
	if qb.rawSQL != "" {
		return qb.rawSQL
	}

	var parts []string

	// SELECT clause
	if qb.distinct {
		parts = append(parts, "SELECT DISTINCT", strings.Join(qb.fields, ", "))
	} else {
		parts = append(parts, "SELECT", strings.Join(qb.fields, ", "))
	}

	// FROM clause
	parts = append(parts, "FROM", qb.table)

	// JOIN clauses
	for _, join := range qb.joins {
		parts = append(parts, fmt.Sprintf("%s JOIN %s ON %s", join.Type, join.Table, join.Condition))
	}

	// WHERE clause
	if len(qb.where) > 0 {
		var conditions []string
		for _, condition := range qb.where {
			if condition.Field != "" {
				if condition.Operator != "" && condition.Value != nil {
					conditions = append(conditions, fmt.Sprintf("%s %s ?", condition.Field, condition.Operator))
				} else if condition.Field != "" {
					conditions = append(conditions, condition.Field)
				}
			}
		}
		if len(conditions) > 0 {
			parts = append(parts, "WHERE", strings.Join(conditions, " AND "))
		}
	}

	// GROUP BY clause
	if len(qb.groupBy) > 0 {
		parts = append(parts, "GROUP BY", strings.Join(qb.groupBy, ", "))
	}

	// HAVING clause
	if qb.having != "" {
		parts = append(parts, "HAVING", qb.having)
	}

	// ORDER BY clause
	if len(qb.orderBy) > 0 {
		var orders []string
		for _, order := range qb.orderBy {
			orders = append(orders, fmt.Sprintf("%s %s", order.Field, order.Direction))
		}
		parts = append(parts, "ORDER BY", strings.Join(orders, ", "))
	}

	// LIMIT clause
	if qb.limit > 0 {
		parts = append(parts, fmt.Sprintf("LIMIT %d", qb.limit))
	}

	// OFFSET clause
	if qb.offset > 0 {
		parts = append(parts, fmt.Sprintf("OFFSET %d", qb.offset))
	}

	// Lock clause
	if qb.lockType != "" {
		parts = append(parts, qb.lockType)
	}

	return strings.Join(parts, " ")
}

// GetArgs returns the query arguments
func (qb *BuilderImpl) GetArgs() []interface{} {
	if qb.rawSQL != "" {
		return qb.rawArgs
	}

	args := make([]interface{}, 0)
	args = append(args, qb.args...)
	args = append(args, qb.havingArgs...)
	return args
}
