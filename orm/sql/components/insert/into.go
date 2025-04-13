package components

func (iq *InsertQuery) Into(table string) *InsertQuery {
	iq.Table = table
	return iq
}
