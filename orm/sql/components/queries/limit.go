package components

import(
	"log"
)

func (q *Query) Limit(limit int) *Query {
	if limit < 0 {
		log.Fatalf("Invalid limit")
		return nil
	}
	q.LimitValue = limit
	return q
}