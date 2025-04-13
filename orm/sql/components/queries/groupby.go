package components

func (q *Query) GroupBy(groupBy string) *Query {
	q.GroupByClause = groupBy
	return q
}