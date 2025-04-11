package utils

import (
	"reflect"
	"project/orm/sql/components"
	"strings"
	"regexp"
	"project/memory"
	"log"
)

type Query struct {
	components.Query
}

func StructSliceToMapSlice(slice interface{}) []map[string]interface{} {
	result := []map[string]interface{}{}
	val := reflect.ValueOf(slice)

	if val.Kind() != reflect.Slice {
		return result
	}

	for i := 0; i < val.Len(); i++ {
		item := val.Index(i)
		itemMap := make(map[string]interface{})

		for j := 0; j < item.NumField(); j++ {
			field := item.Type().Field(j)
			fieldName := field.Name
			fieldValue := item.Field(j).Interface()
			itemMap[fieldName] = fieldValue
		}

		result = append(result, itemMap)
	}

	return result
}

func EvaluateCondition(condition string) (bool) {
	condition = strings.TrimSpace(condition)

	if len(strings.Fields(condition)) < 3 || condition == "" {
		return false
	}

	validConditionRegex := `^(?i)([\w\d_]+)\s*(=|!=|<>|<|<=|>|>=|BETWEEN|IN)\s*(.*)$`

	re := regexp.MustCompile(validConditionRegex)

	if !re.MatchString(condition) {
		return false
	}

	return true
}

func FilterData(table string) ([]map[string]interface{}) {
	// TODO recuperer dynamiquement le nom de la table
	switch table {
	case "users":
		return StructSliceToMapSlice(data.Users)
	case "posts":
		return StructSliceToMapSlice(data.Posts)
	default:
		log.Fatalf("Invalid from")
		return nil
	}
}

func FieldInsensitive(data map[string]interface{}, field string) string {
	field = strings.ToLower(field)

	for key := range data {
		if strings.ToLower(key) == field {
			return key
		}
	}

	return ""
}