package orm

import (
	"log"
	"project/orm/utils"
	"strconv"
	"strings"

	insert "project/orm/sql/components/insert"
	queries "project/orm/sql/components/queries"
)

type Query struct {
	queries.Query
}

type InsertQuery struct {
	insert.InsertQuery
}

func (q *Query) Apply(datas []map[string]interface{}) []map[string]interface{} {
	var result []map[string]interface{}

	if !utils.EvaluateCondition(&q.Query) {
		log.Fatalf("Invalid request")
		return nil
	}

	parts := strings.Fields(q.WhereClause)

	field := parts[0]
	operator := parts[1]
	value := parts[2]

	q.Fields = utils.AsterixValue(&q.Query, datas[0])

	for _, data := range datas {
		include := false
		field := utils.FieldInsensitive(data, field)

		if fieldValue, ok := data[field]; ok {
			parsedValue, err := strconv.Atoi(value)
			if err != nil {
				log.Fatalf("Fail convertion from %s to int: %v", value, err)
			}
			include = Operation(operator, fieldValue, parsedValue, value, parts)
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

func Operation(operator string, fieldValue interface{}, parsedValue int, value string, other []string) bool {
	operator = strings.ToLower(operator)
	switch operator {
	case ">":
		if fieldValue, ok := fieldValue.(int); ok {
			if fieldValue > parsedValue {
				return true
			}
		}

	case "<":
		if fieldValue, ok := fieldValue.(int); ok {
			if fieldValue < parsedValue {
				return true
			}
		}

	case "=":
		parsedVal, err := strconv.Atoi(value)
		if err != nil {
			return false
		}
		if fieldValue.(int) == parsedVal{
			return true
		}

	case ">=":
		if fieldValue, ok := fieldValue.(int); ok {
			if fieldValue >= parsedValue {
				return true
			}
		}

	case "<=":
		if fieldValue, ok := fieldValue.(int); ok {
			if fieldValue <= parsedValue {
				return true
			}
		}

	case "<>":
		if fieldValue, ok := fieldValue.(int); ok {
			if fieldValue != parsedValue {
				return true
			}
		}

	case "between":
		first, err := strconv.Atoi(other[2])
		second, err := strconv.Atoi(other[4])

		if err != nil {
			log.Fatalf("Fail convertion from %s to int: %v", other[4], err)
		}

		if fieldValue, ok := fieldValue.(int); ok {
			if fieldValue >= first && fieldValue <= second {
				return true
			}
		}

	case "like":
		// TODO
	case "in":
		// TODO	
	default:
		log.Fatalf("Operation not supported: %s", operator)
	}
	return false
}
