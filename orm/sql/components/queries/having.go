package components

func (q *Query) Having(having string) *Query {
	q.HavingClause = having
	return q
}