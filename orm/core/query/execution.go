package query

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/ESGI-M2/GO/orm/core/interfaces"
)

// Find executes the query and returns all results
func (qb *BuilderImpl) Find() ([]map[string]interface{}, error) {
	if qb.Err != nil {
		return nil, qb.Err
	}

	// Check cache first
	if qb.useCache {
		if cached, found := qb.getFromCache(); found {
			return cached, nil
		}
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

	results, err := qb.scanRows(rows)
	if err != nil {
		return nil, err
	}

	// Load relations if specified
	if len(qb.withRelations) > 0 {
		results, err = qb.loadRelations(results)
		if err != nil {
			return nil, err
		}
	}

	// Cache results if enabled
	if qb.useCache {
		qb.setCache(results)
	}

	return results, nil
}

// FindOne executes the query and returns one result
func (qb *BuilderImpl) FindOne() (map[string]interface{}, error) {
	if qb.Err != nil {
		return nil, qb.Err
	}

	// Check cache first
	if qb.useCache {
		if cached, found := qb.getFromCache(); found {
			if len(cached) > 0 {
				return cached[0], nil
			}
			return nil, nil
		}
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

	result := results[0]

	// Load relations if specified
	if len(qb.withRelations) > 0 {
		results, err := qb.loadRelations([]map[string]interface{}{result})
		if err != nil {
			return nil, err
		}
		if len(results) > 0 {
			result = results[0]
		}
	}

	// Cache result if enabled
	if qb.useCache {
		qb.setCache([]map[string]interface{}{result})
	}

	return result, nil
}

// Count executes a COUNT query
func (qb *BuilderImpl) Count() (int64, error) {
	if qb.Err != nil {
		return 0, qb.Err
	}

	if qb.rawSQL != "" {
		return 0, fmt.Errorf("count not supported for raw SQL")
	}

	// Check cache first
	if qb.useCache {
		if cached, found := qb.getCountFromCache(); found {
			return cached, nil
		}
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

	// Cache count if enabled
	if qb.useCache {
		qb.setCountCache(count)
	}

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

	// Check cache first
	if qb.useCache {
		if cached, found := qb.getExistsFromCache(); found {
			return cached, nil
		}
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

	// Cache exists result if enabled
	if qb.useCache {
		qb.setExistsCache(exists)
	}

	return exists, nil
}

// Paginate executes the query with pagination
func (qb *BuilderImpl) Paginate(page, perPage int) (*interfaces.PaginationResult, error) {
	if qb.Err != nil {
		return nil, qb.Err
	}

	// Get total count
	total, err := qb.Count()
	if err != nil {
		return nil, err
	}

	// Set pagination
	qb.OffsetPaginate(page, perPage)

	// Get data
	data, err := qb.Find()
	if err != nil {
		return nil, err
	}

	// Calculate pagination info
	lastPage := int((total + int64(perPage) - 1) / int64(perPage))
	from := (page-1)*perPage + 1
	to := from + len(data) - 1
	hasMore := page < lastPage

	return &interfaces.PaginationResult{
		Data:        qb.convertToInterface(data),
		Total:       total,
		PerPage:     perPage,
		CurrentPage: page,
		LastPage:    lastPage,
		From:        from,
		To:          to,
		HasMore:     hasMore,
	}, nil
}

// buildQuery builds the SQL query string
func (qb *BuilderImpl) buildQuery() string {
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
		argIndex := 0
		for _, condition := range qb.where {
			if condition.Field != "" {
				if condition.Raw {
					// Raw WHERE conditions are used as-is
					conditions = append(conditions, condition.Field)
				} else if condition.Operator != "" && condition.Value != nil {
					conditions = append(conditions, fmt.Sprintf("%s %s %s", condition.Field, condition.Operator, qb.Orm.GetDialect().GetPlaceholder(argIndex)))
					argIndex++
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
		for i, column := range columns {
			val := values[i]
			if val != nil {
				row[column] = val
			}
		}
		results = append(results, row)
	}

	return results, nil
}

// loadRelations loads related data for eager loading
func (qb *BuilderImpl) loadRelations(results []map[string]interface{}) ([]map[string]interface{}, error) {
	if len(results) == 0 {
		return results, nil
	}

	// Get primary keys
	var ids []interface{}
	for _, result := range results {
		if id, ok := result["id"]; ok {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		return results, nil
	}

	// Load each relation using metadata information
	for relationName, relationFn := range qb.withRelations {
		relInfo, ok := qb.Metadata.Relations[relationName]
		if !ok {
			continue // unknown relation
		}

		// Determine the actual model type (handle slice for has_many)
		modelType := relInfo.TargetModel
		if modelType.Kind() == reflect.Slice {
			modelType = modelType.Elem()
		}

		relModelPtr := reflect.New(modelType).Interface()

		relationQuery := qb.Orm.Query(relModelPtr)
		relationQuery = relationFn(relationQuery).(interfaces.QueryBuilder)

		// Use the foreign key defined in relation metadata
		fkField := relInfo.ForeignKey
		if fkField == "" {
			fkField = "user_id" // fallback
		}

		relationQuery.WhereIn(fkField, ids)

		relationResults, err := relationQuery.Find()
		if err != nil {
			return nil, fmt.Errorf("failed to load relation %s: %w", relationName, err)
		}

		relationMap := make(map[interface{}][]map[string]interface{})
		for _, relationResult := range relationResults {
			if fk, ok := relationResult[fkField]; ok {
				relationMap[fk] = append(relationMap[fk], relationResult)
			}
		}

		for i, result := range results {
			if id, ok := result["id"]; ok {
				if rels, exists := relationMap[id]; exists {
					results[i][relationName] = rels
				} else {
					results[i][relationName] = []map[string]interface{}{}
				}
			}
		}
	}

	return results, nil
}

// getCacheKey generates a cache key for the current query
func (qb *BuilderImpl) getCacheKey() string {
	data := map[string]interface{}{
		"sql":  qb.GetSQL(),
		"args": qb.GetArgs(),
	}

	jsonData, _ := json.Marshal(data)
	hash := md5.Sum(jsonData)
	return hex.EncodeToString(hash[:])
}

// getFromCache retrieves results from cache
func (qb *BuilderImpl) getFromCache() ([]map[string]interface{}, bool) {
	// This is a simplified cache implementation
	// In a real implementation, you would use a proper cache like Redis or in-memory cache
	return nil, false
}

// setCache stores results in cache
func (qb *BuilderImpl) setCache(results []map[string]interface{}) {
	// This is a simplified cache implementation
	// In a real implementation, you would use a proper cache like Redis or in-memory cache
}

// getCountFromCache retrieves count from cache
func (qb *BuilderImpl) getCountFromCache() (int64, bool) {
	return 0, false
}

// setCountCache stores count in cache
func (qb *BuilderImpl) setCountCache(count int64) {
	// Simplified cache implementation
}

// getExistsFromCache retrieves exists result from cache
func (qb *BuilderImpl) getExistsFromCache() (bool, bool) {
	return false, false
}

// setExistsCache stores exists result in cache
func (qb *BuilderImpl) setExistsCache(exists bool) {
	// Simplified cache implementation
}

// convertToInterface converts []map[string]interface{} to []interface{}
func (qb *BuilderImpl) convertToInterface(data []map[string]interface{}) []interface{} {
	result := make([]interface{}, len(data))
	for i, item := range data {
		result[i] = item
	}
	return result
}
