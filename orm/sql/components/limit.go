package components

func (q *Query) Limit(limit int) *Query {
	q.LimitValue = limit
	return q
}