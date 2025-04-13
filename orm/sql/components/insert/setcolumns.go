package components

import (
	"log"
)

func (iq *InsertQuery) Set(columns []string, values []interface{}) *InsertQuery {
	if len(columns) != len(values) {
		log.Fatalf("Number of columns and values must match")
	}
	iq.Columns = columns
	iq.Values = values
	return iq
}
