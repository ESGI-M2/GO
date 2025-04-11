package orm

import (
	"fmt"
	"reflect"
	"strings"
)

func CreateTableSQL(tableName string, model interface{}) string {
	t := reflect.TypeOf(model)

	var fields []string
	var foreigns []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" || dbTag == "-" {
			continue
		}

		sqlType := goTypeToSQLType(field.Type.Name())

		definition := fmt.Sprintf("%s %s", dbTag, sqlType)

		if field.Tag.Get("primary") == "true" {
			definition += " PRIMARY KEY"
		}
		if field.Tag.Get("autoincrement") == "true" {
			definition += " AUTO_INCREMENT"
		}
		if fk := field.Tag.Get("foreign"); fk != "" {
			foreigns = append(foreigns, fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s", dbTag, fk))
		}

		fields = append(fields, definition)
	}

	allFields := append(fields, foreigns...)
	return fmt.Sprintf("CREATE TABLE %s (\n  %s\n);", tableName, strings.Join(allFields, ",\n  "))
}

func goTypeToSQLType(goType string) string {
	switch goType {
	case "int":
		return "INT"
	case "string":
		return "VARCHAR(255)"
	case "bool":
		return "BOOLEAN"
	default:
		return "TEXT"
	}
}
