package components

func (q *Query) Where(condition string) *Query {
	q.WhereClause = condition
	return q
}