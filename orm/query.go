package orm

import (
	"fmt"
	"strings"
)

func Select(table string, fields []string, where string, limit int) string {
	selectFields := "*"
	if len(fields) > 0 {
		selectFields = strings.Join(fields, ", ")
	}

	whereClause := ""
	if where != "" {
		whereClause = "WHERE " + where
	}

	limitClause := ""
	if limit > 0 {
		limitClause = fmt.Sprintf("LIMIT %d", limit)
	}

	query := fmt.Sprintf("SELECT %s FROM %s %s %s", selectFields, table, whereClause, limitClause)
	query = strings.TrimSpace(query)
	return query
}
