package orm

import (
	"fmt"
	"reflect"
	"strings"
)

func Update(tableName string, model interface{}, where string) (string, []interface{}) {
	v := reflect.ValueOf(model)
	t := reflect.TypeOf(model)

	var setClauses []string
	var values []interface{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "id" {
			setClauses = append(setClauses, fmt.Sprintf("%s = ?", dbTag))
			values = append(values, v.Field(i).Interface())
		}
	}

	whereClause := ""
	if where != "" {
		whereClause = "WHERE " + where
	}

	query := fmt.Sprintf("UPDATE %s SET %s %s", tableName, strings.Join(setClauses, ", "), whereClause)

	return query, values
}
