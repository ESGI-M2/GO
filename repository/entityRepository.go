package repository

import (
	"project/memory"
	"reflect"
	"project/orm/utils"
)

func FindAll(table string) []map[string]interface{} {
	slice := data.Store[table]
	return utils.StructSliceToMapSlice(slice)
}

func FindBy(table string, field string, value interface{}) []map[string]interface{} {
	all := FindAll(table)
	result := []map[string]interface{}{}

	for _, row := range all {
		if rowVal, ok := row[field]; ok && reflect.DeepEqual(rowVal, value) {
			result = append(result, row)
		}
	}

	return result
}

func FindOneBy(table string, field string, value interface{}) map[string]interface{} {
	matches := FindBy(table, field, value)
	if len(matches) > 0 {
		return matches[0]
	}
	return nil
}

func Find(table string, criteria map[string]interface{}) []map[string]interface{} {
	all := FindAll(table)
	result := []map[string]interface{}{}

	for _, row := range all {
		match := true
		for field, val := range criteria {
			if rowVal, ok := row[field]; !ok || rowVal != val {
				match = false
				break
			}
		}
		if match {
			result = append(result, row)
		}
	}
	return result
}
