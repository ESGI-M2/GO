package components

import (
	"log"
	"regexp"
)

func (q *Query) Select(fields ...string) *Query {
	validField := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

	if len(fields) == 1 && fields[0] == "*" {
		q.Fields = fields
		return q
	}

	for _, field := range fields {
		if !validField.MatchString(field) || field == "*" {
			log.Fatalf("Invalid select")
			return nil
		}
	}

	q.Fields = fields
	return q
}

