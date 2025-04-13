package components

type InsertQuery struct {
	Table   string
	Columns []string
	Values  []interface{}
}