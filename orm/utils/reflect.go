package utils

import (
	"log"
	"project/memory"
	"project/models"
	"project/orm/sql/components/queries"
	"reflect"
	"regexp"
	"strings"
)

var EntityRegistry = map[string]interface{}{
	"users": []models.User{},
	"posts": []models.Post{},
}

type Query struct {
	components.Query
}

func GetEntityType(name string) reflect.Type {
	entity := EntityRegistry[name]

	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}

	return t
}

func GetFieldsFromType(t reflect.Type) []string {
	var fields []string

	if t.Kind() != reflect.Struct {
		return fields
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fields = append(fields, field.Name)
	}

	return fields
}


func SetFields(entityType reflect.Type, columns []string, values []interface{}) reflect.Value {
    v := reflect.New(entityType).Elem()
    valueIndex := 0

    for _, col := range columns {
        if col == "ID" {
            continue
        }
        
        field := v.FieldByName(col)

        if field.IsValid() && field.CanSet() {
            val := reflect.ValueOf(values[valueIndex])

            valueIndex++

            if val.Type().ConvertibleTo(field.Type()) {
                field.Set(val.Convert(field.Type()))
            } else {
                log.Printf("Le type de la valeur n'est pas convertible au champ: %s\n", col)
            }
        } else {
            log.Printf("Champ invalide: %s\n", col)
        }
    }
    
    return v
}

func AppendToData(table string, obj reflect.Value) {
    sliceVal := StructSliceToArray(table)

    if obj.Kind() == reflect.Struct {
        maxID := 0
        for i := 0; i < sliceVal.Len(); i++ {
            item := sliceVal.Index(i)
            idField := item.FieldByName("ID")
            if idField.IsValid() && idField.Kind() == reflect.Int {
                id := int(idField.Int())
                if id > maxID {
                    maxID = id
                }
            }
        }

        idField := obj.FieldByName("ID")
        if idField.IsValid() && idField.CanSet() && idField.Kind() == reflect.Int {
            idField.SetInt(int64(maxID + 1))
        }
    }

    sliceVal.Set(reflect.Append(sliceVal, obj))
}


func StructSliceToArray(table string) (reflect.Value) {
	listPtr, ok := data.Store[table]
    if !ok {
        log.Fatalf("Table non trouvée : %s", table)
    }

    slicePtrVal := reflect.ValueOf(listPtr)
    if slicePtrVal.Kind() != reflect.Ptr {
        log.Fatalf("Store[%s] n'est pas un pointeur", table)
    }

	return slicePtrVal.Elem()
}

func StructSliceToMapSlice(slice interface{}) []map[string]interface{} {
    result := []map[string]interface{}{}
    val := reflect.ValueOf(slice)

    if val.Kind() != reflect.Slice && val.Kind() != reflect.Ptr {
        return result
    }

    if val.Kind() == reflect.Ptr {
        val = val.Elem()
    }

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


func EvaluateCondition(q *components.Query) bool {
	condition := strings.TrimSpace(q.WhereClause)

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
	slice := data.Store[table]
	return StructSliceToMapSlice(slice)
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

func AsterixValue(q *components.Query, data map[string]interface{}) []string {
	if q.Fields[0] == "*" {
		q.Fields = []string{}
		for key := range data {
			q.Fields = append(q.Fields, key)
		}
	}
	return q.Fields
}