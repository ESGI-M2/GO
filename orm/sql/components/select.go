package components

func (q *Query) Select(fields ...string) *Query {
	if len(fields) > 0 {
		q.Fields = fields
	} else {
		q.Fields = append(q.Fields, "*")
	}
	return q
}