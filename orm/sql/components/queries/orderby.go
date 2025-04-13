package components

func (q *Query) OrderBy(orderBy string) *Query {
	q.OrderByClause = orderBy
	return q
}