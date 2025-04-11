package orm

import (
	"log"
	"project/orm/sql/components"
	"project/orm/utils"
	"strings"
	"strconv"
)

type Query struct {
	components.Query
}

func (q *Query) Apply(datas []map[string]interface{}) []map[string]interface{} {
	var result []map[string]interface{}

	// TODO Essayer de passer la Query en parametres 
	if !utils.EvaluateCondition(q.WhereClause) {
		log.Fatalf("Invalid where")
		return nil
	}

	parts := strings.Fields(q.WhereClause)

	field := parts[0]
	operator := parts[1]
	value := parts[2]

	for _, data := range datas {
		include := false
		field := utils.FieldInsensitive(data, field)

		if fieldValue, ok := data[field]; ok {
			parsedValue, err := strconv.Atoi(value)
			if err != nil {
				log.Fatalf("Fail convertion from %s to int: %v", value, err)
			}
			switch operator {
			case ">":
				if fieldValue, ok := fieldValue.(int); ok {
					if fieldValue > parsedValue {
						include = true
					}
				}

			case "<":
				if fieldValue, ok := fieldValue.(int); ok {
					if fieldValue < parsedValue {
						include = true
					}
				}

			case "=":
				if fieldValue == value {
					include = true
				}

			case ">=":
				if fieldValue, ok := fieldValue.(int); ok {
					if fieldValue >= parsedValue {
						include = true
					}
				}

			case "<=":
				if fieldValue, ok := fieldValue.(int); ok {
					if fieldValue <= parsedValue {
						include = true
					}
				}
			default:
				log.Fatalf("Operation not supported: %s", operator)
			}
		}

		if include {
			allData := make(map[string]interface{})

			for _, field := range q.Fields {
				allData[field] = data[field]
			}

			result = append(result, allData)

			if q.LimitValue > 0 && len(result) >= q.LimitValue {
				break
			}
		}
	}

	return result
}
