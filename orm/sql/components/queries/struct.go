package components

type Query struct {
	Fields         []string
	Table          string
	WhereClause    string
	LimitValue     int
	Joins          []string
	GroupByClause  string
	HavingClause   string
	OrderByClause  string
}
