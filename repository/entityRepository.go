package repository

import (
	data "project/memory"
	"reflect"
)

func FindAll(table string) []map[string]interface{} {
	slice := data.Store[table]
	return structSliceToMapSlice(slice)
}

// structSliceToMapSlice converts a slice of structs to a slice of maps
func structSliceToMapSlice(slice interface{}) []map[string]interface{} {
	if slice == nil {
		return []map[string]interface{}{}
	}

	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice {
		return []map[string]interface{}{}
	}

	result := make([]map[string]interface{}, sliceValue.Len())
	for i := 0; i < sliceValue.Len(); i++ {
		item := sliceValue.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}

		if item.Kind() == reflect.Struct {
			result[i] = structToMap(item)
		}
	}

	return result
}

// structToMap converts a struct to a map
func structToMap(v reflect.Value) map[string]interface{} {
	result := make(map[string]interface{})

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		value := v.Field(i)

		if value.CanInterface() {
			result[field.Name] = value.Interface()
		}
	}

	return result
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
