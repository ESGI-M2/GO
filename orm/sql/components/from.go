package components

func (q *Query) From(table string) *Query {
	q.Table = table
	return q
}