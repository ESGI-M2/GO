package orm

import (
	"fmt"
	"reflect"
	"strings"
)

func Insert(tableName string, model interface{}) (string, []interface{}) {
	v := reflect.ValueOf(model)
	t := reflect.TypeOf(model)

	var columns []string
	var placeholders []string
	var values []interface{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Tag.Get("autoincrement") == "true" {
			continue
		}

		dbTag := field.Tag.Get("db")
		if dbTag == "-" || dbTag == "" {
			continue
		}

		columns = append(columns, dbTag)
		placeholders = append(placeholders, "?")
		values = append(values, v.Field(i).Interface())
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	return query, values
}
