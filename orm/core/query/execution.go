package query

import (
	"database/sql"
	"fmt"
	"strings"
)

// Find executes the query and returns all results
func (qb *BuilderImpl) Find() ([]map[string]interface{}, error) {
	if qb.Err != nil {
		return nil, qb.Err
	}

	if qb.rawSQL != "" {
		return qb.executeRaw()
	}

	sql := qb.buildQuery()
	rows, err := qb.Orm.GetDialect().Query(sql, qb.args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	if rows != nil {
		defer rows.Close()
	}

	return qb.scanRows(rows)
}

// FindOne executes the query and returns one result
func (qb *BuilderImpl) FindOne() (map[string]interface{}, error) {
	if qb.Err != nil {
		return nil, qb.Err
	}

	if qb.rawSQL != "" {
		results, err := qb.executeRaw()
		if err != nil {
			return nil, err
		}
		if len(results) == 0 {
			return nil, nil
		}
		return results[0], nil
	}

	// Add LIMIT 1 for single result
	originalLimit := qb.limit
	qb.limit = 1

	sql := qb.buildQuery()
	rows, err := qb.Orm.GetDialect().Query(sql, qb.args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	if rows != nil {
		defer rows.Close()
	}

	results, err := qb.scanRows(rows)
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

// Count executes a COUNT query
func (qb *BuilderImpl) Count() (int64, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}

	if qb.rawSQL != "" {
		return 0, fmt.Errorf("count not supported for raw SQL")
	}

	// Save original fields and set to COUNT(*)
	originalFields := qb.fields
	qb.fields = []string{"COUNT(*)"}

	sql := qb.buildQuery()
	row := qb.Orm.GetDialect().QueryRow(sql, qb.args...)

	var count int64
	if row != nil {
		err := row.Scan(&count)
		if err != nil {
			return 0, fmt.Errorf("failed to scan count: %w", err)
		}
	} else {
		count = 0
	}

	// Restore original fields
	qb.fields = originalFields

	return count, nil
}

// Exists checks if any records exist
func (qb *BuilderImpl) Exists() (bool, error) {
	if qb.Err != nil {
		return false, qb.Err
	}

	if qb.rawSQL != "" {
		results, err := qb.executeRaw()
		if err != nil {
			return false, err
		}
		return len(results) > 0, nil
	}

	// Save original fields and set to SELECT 1
	originalFields := qb.fields
	qb.fields = []string{"1"}

	// Save original limit and set to LIMIT 1
	originalLimit := qb.limit
	qb.limit = 1

	sql := qb.buildQuery()
	rows, err := qb.Orm.GetDialect().Query(sql, qb.args...)
	if err != nil {
		return false, fmt.Errorf("failed to execute exists query: %w", err)
	}

	var exists bool
	if rows != nil {
		defer rows.Close()
		exists = rows.Next()
	} else {
		exists = false
	}

	// Restore original fields and limit
	qb.fields = originalFields
	qb.limit = originalLimit

	return exists, nil
}

// buildQuery builds the SQL query string
func (qb *BuilderImpl) buildQuery() string {
	var parts []string

	// SELECT clause
	parts = append(parts, "SELECT", strings.Join(qb.fields, ", "))

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

	return strings.Join(parts, " ")
}

// executeRaw executes a raw SQL query
func (qb *BuilderImpl) executeRaw() ([]map[string]interface{}, error) {
	rows, err := qb.Orm.GetDialect().Query(qb.rawSQL, qb.rawArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute raw query: %w", err)
	}
	if rows != nil {
		defer rows.Close()
	}

	return qb.scanRows(rows)
}

// scanRows scans database rows into a slice of maps
func (qb *BuilderImpl) scanRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	if rows == nil {
		return []map[string]interface{}{}, nil
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = values[i]
		}
		results = append(results, row)
	}

	return results, nil
}
