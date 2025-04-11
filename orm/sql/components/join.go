package components

import "fmt"

func Join(joinType, table, onClause string) string {
	if table != "" && onClause != "" {
		return fmt.Sprintf("%s JOIN %s ON %s", joinType, table, onClause)
	}
	return ""
}

// TODO faire la logique des jointures

func InnerJoin(table, onClause string) string {
	return Join("INNER", table, onClause)
}

func LeftJoin(table, onClause string) string {
	return Join("LEFT", table, onClause)
}

func RightJoin(table, onClause string) string {
	return Join("RIGHT", table, onClause)
}


func (q *Query) InnerJoin(table, onClause string) *Query {
	joinClause := InnerJoin(table, onClause)
	if joinClause != "" {
		q.Joins = append(q.Joins, joinClause)
	}
	return q
}

func (q *Query) LeftJoin(table, onClause string) *Query {
	joinClause := LeftJoin(table, onClause)
	if joinClause != "" {
		q.Joins = append(q.Joins, joinClause)
	}
	return q
}

func (q *Query) RightJoin(table, onClause string) *Query {
	joinClause := RightJoin(table, onClause)
	if joinClause != "" {
		q.Joins = append(q.Joins, joinClause)
	}
	return q
}
