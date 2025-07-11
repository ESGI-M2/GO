package core

import (
	"database/sql"
	"fmt"
	"strings"
)

// WhereCondition represents a WHERE clause condition
type WhereCondition struct {
	Field    string
	Operator string
	Value    interface{}
	Logical  string // AND, OR
}

// OrderBy represents an ORDER BY clause
type OrderBy struct {
	Field     string
	Direction string // ASC, DESC
}

// Join represents a JOIN clause
type Join struct {
	Type      string // INNER, LEFT, RIGHT, FULL
	Table     string
	Condition string
}

// QueryBuilderImpl implements the QueryBuilder interface
type QueryBuilderImpl struct {
	orm      *ORMImpl
	metadata *ModelMetadata
	err      error

	// Query components
	table      string
	fields     []string
	where      []WhereCondition
	orderBy    []OrderBy
	groupBy    []string
	having     string
	havingArgs []interface{}
	joins      []Join
	limit      int
	offset     int
	args       []interface{}

	// Raw SQL
	rawSQL  string
	rawArgs []interface{}
}

// Select sets the fields to select
func (qb *QueryBuilderImpl) Select(fields ...string) QueryBuilder {
	if qb.err != nil {
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
func (qb *QueryBuilderImpl) From(table string) QueryBuilder {
	if qb.err != nil {
		return qb
	}
	qb.table = table
	return qb
}

// Where adds a WHERE condition
func (qb *QueryBuilderImpl) Where(field, operator string, value interface{}) QueryBuilder {
	if qb.err != nil {
		return qb
	}

	qb.where = append(qb.where, WhereCondition{
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
func (qb *QueryBuilderImpl) WhereIn(field string, values []interface{}) QueryBuilder {
	if qb.err != nil {
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
	qb.where = append(qb.where, WhereCondition{
		Field:    condition,
		Operator: "",
		Value:    nil,
		Logical:  "AND",
	})

	return qb
}

// WhereNotIn adds a WHERE NOT IN condition
func (qb *QueryBuilderImpl) WhereNotIn(field string, values []interface{}) QueryBuilder {
	if qb.err != nil {
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
	qb.where = append(qb.where, WhereCondition{
		Field:    condition,
		Operator: "",
		Value:    nil,
		Logical:  "AND",
	})

	return qb
}

// OrderBy adds an ORDER BY clause
func (qb *QueryBuilderImpl) OrderBy(field, direction string) QueryBuilder {
	if qb.err != nil {
		return qb
	}

	if direction == "" {
		direction = "ASC"
	}

	qb.orderBy = append(qb.orderBy, OrderBy{
		Field:     field,
		Direction: strings.ToUpper(direction),
	})

	return qb
}

// GroupBy adds a GROUP BY clause
func (qb *QueryBuilderImpl) GroupBy(fields ...string) QueryBuilder {
	if qb.err != nil {
		return qb
	}
	qb.groupBy = append(qb.groupBy, fields...)
	return qb
}

// Having adds a HAVING clause
func (qb *QueryBuilderImpl) Having(condition string, args ...interface{}) QueryBuilder {
	if qb.err != nil {
		return qb
	}
	qb.having = condition
	qb.havingArgs = append(qb.havingArgs, args...)
	return qb
}

// Limit sets the LIMIT clause
func (qb *QueryBuilderImpl) Limit(limit int) QueryBuilder {
	if qb.err != nil {
		return qb
	}
	qb.limit = limit
	return qb
}

// Offset sets the OFFSET clause
func (qb *QueryBuilderImpl) Offset(offset int) QueryBuilder {
	if qb.err != nil {
		return qb
	}
	qb.offset = offset
	return qb
}

// Join adds a JOIN clause
func (qb *QueryBuilderImpl) Join(table, condition string) QueryBuilder {
	return qb.addJoin("INNER", table, condition)
}

// LeftJoin adds a LEFT JOIN clause
func (qb *QueryBuilderImpl) LeftJoin(table, condition string) QueryBuilder {
	return qb.addJoin("LEFT", table, condition)
}

// RightJoin adds a RIGHT JOIN clause
func (qb *QueryBuilderImpl) RightJoin(table, condition string) QueryBuilder {
	return qb.addJoin("RIGHT", table, condition)
}

// InnerJoin adds an INNER JOIN clause
func (qb *QueryBuilderImpl) InnerJoin(table, condition string) QueryBuilder {
	return qb.addJoin("INNER", table, condition)
}

// addJoin adds a join clause
func (qb *QueryBuilderImpl) addJoin(joinType, table, condition string) QueryBuilder {
	if qb.err != nil {
		return qb
	}

	qb.joins = append(qb.joins, Join{
		Type:      joinType,
		Table:     table,
		Condition: condition,
	})

	return qb
}

// Find executes the query and returns all results
func (qb *QueryBuilderImpl) Find() ([]map[string]interface{}, error) {
	if qb.err != nil {
		return nil, qb.err
	}

	if qb.rawSQL != "" {
		return qb.executeRaw()
	}

	query := qb.buildQuery()
	rows, err := qb.orm.dialect.Query(query, qb.args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	return qb.scanRows(rows)
}

// FindOne executes the query and returns the first result
func (qb *QueryBuilderImpl) FindOne() (map[string]interface{}, error) {
	if qb.err != nil {
		return nil, qb.err
	}

	// Set limit to 1 for FindOne
	originalLimit := qb.limit
	qb.limit = 1

	results, err := qb.Find()
	if err != nil {
		return nil, err
	}

	// Restore original limit
	qb.limit = originalLimit

	if len(results) == 0 {
		return nil, nil
	}

	return results[0], nil
}

// Count executes the query and returns the count
func (qb *QueryBuilderImpl) Count() (int64, error) {
	if qb.err != nil {
		return 0, qb.err
	}

	// Save original fields and set to COUNT(*)
	originalFields := qb.fields
	qb.fields = []string{"COUNT(*)"}

	query := qb.buildQuery()
	var count int64
	row := qb.orm.dialect.QueryRow(query, qb.args...)
	if row == nil {
		qb.fields = originalFields
		return 0, fmt.Errorf("failed to count: QueryRow returned nil")
	}
	err := row.Scan(&count)

	// Restore original fields
	qb.fields = originalFields

	if err != nil {
		return 0, fmt.Errorf("failed to count: %w", err)
	}

	return count, nil
}

// Exists checks if any records exist
func (qb *QueryBuilderImpl) Exists() (bool, error) {
	if qb.err != nil {
		return false, qb.err
	}

	count, err := qb.Count()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Raw sets raw SQL
func (qb *QueryBuilderImpl) Raw(sql string, args ...interface{}) QueryBuilder {
	qb.rawSQL = sql
	qb.rawArgs = args
	return qb
}

// GetSQL returns the built SQL query
func (qb *QueryBuilderImpl) GetSQL() string {
	if qb.rawSQL != "" {
		return qb.rawSQL
	}
	return qb.buildQuery()
}

// GetArgs returns the query arguments
func (qb *QueryBuilderImpl) GetArgs() []interface{} {
	if qb.rawSQL != "" {
		return qb.rawArgs
	}
	return qb.args
}

// buildQuery builds the SQL query string
func (qb *QueryBuilderImpl) buildQuery() string {
	var parts []string

	// SELECT
	fields := strings.Join(qb.fields, ", ")
	parts = append(parts, fmt.Sprintf("SELECT %s", fields))

	// FROM
	parts = append(parts, fmt.Sprintf("FROM %s", qb.table))

	// JOINs
	for _, join := range qb.joins {
		parts = append(parts, fmt.Sprintf("%s JOIN %s ON %s",
			join.Type, join.Table, join.Condition))
	}

	// WHERE
	if len(qb.where) > 0 {
		whereClauses := make([]string, 0, len(qb.where))
		for i, condition := range qb.where {
			if i == 0 {
				whereClauses = append(whereClauses, condition.Field+" "+condition.Operator+" ?")
			} else {
				whereClauses = append(whereClauses,
					condition.Logical+" "+condition.Field+" "+condition.Operator+" ?")
			}
		}
		parts = append(parts, "WHERE "+strings.Join(whereClauses, " "))
	}

	// GROUP BY
	if len(qb.groupBy) > 0 {
		parts = append(parts, "GROUP BY "+strings.Join(qb.groupBy, ", "))
	}

	// HAVING
	if qb.having != "" {
		parts = append(parts, "HAVING "+qb.having)
	}

	// ORDER BY
	if len(qb.orderBy) > 0 {
		orderClauses := make([]string, 0, len(qb.orderBy))
		for _, order := range qb.orderBy {
			orderClauses = append(orderClauses,
				fmt.Sprintf("%s %s", order.Field, order.Direction))
		}
		parts = append(parts, "ORDER BY "+strings.Join(orderClauses, ", "))
	}

	// LIMIT
	if qb.limit > 0 {
		parts = append(parts, fmt.Sprintf("LIMIT %d", qb.limit))
	}

	// OFFSET
	if qb.offset > 0 {
		parts = append(parts, fmt.Sprintf("OFFSET %d", qb.offset))
	}

	return strings.Join(parts, " ")
}

// executeRaw executes a raw SQL query
func (qb *QueryBuilderImpl) executeRaw() ([]map[string]interface{}, error) {
	rows, err := qb.orm.dialect.Query(qb.rawSQL, qb.rawArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute raw query: %w", err)
	}
	defer rows.Close()

	return qb.scanRows(rows)
}

// scanRows scans database rows into a slice of maps
func (qb *QueryBuilderImpl) scanRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var results []map[string]interface{}

	for rows.Next() {
		// Create a slice of interface{} to hold the values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan the row
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Create a map for this row
		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if val != nil {
				row[col] = val
			}
		}

		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return results, nil
}
