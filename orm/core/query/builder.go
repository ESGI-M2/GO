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
}

// NewBuilder creates a new query builder
func NewBuilder(orm interfaces.ORM, metadata *interfaces.ModelMetadata) *BuilderImpl {
	return &BuilderImpl{
		Orm:      orm,
		Metadata: metadata,
		table:    metadata.TableName,
		fields:   []string{"*"},
		where:    make([]interfaces.WhereCondition, 0),
		orderBy:  make([]interfaces.OrderBy, 0),
		joins:    make([]interfaces.Join, 0),
		limit:    0,
		offset:   0,
		args:     make([]interface{}, 0),
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
		placeholders[i] = "?"
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
		placeholders[i] = "?"
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
	return qb.buildQuery()
}

// GetArgs returns the query arguments
func (qb *BuilderImpl) GetArgs() []interface{} {
	if qb.rawSQL != "" {
		return qb.rawArgs
	}
	return qb.args
}
