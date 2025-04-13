package orm

import (
	"fmt"
	"project/orm/utils"
	"testing"
)

func TestQuery(t *testing.T) {
	query := &Query{}
	query.Select("*").
		From("users").
		Where("age > 30").
		Limit(2)

	table := utils.FilterData(query.Table)

	filteredUsers := query.Apply(table)
	fmt.Print("\n Filtre: ")
	fmt.Print(filteredUsers)
}

func TestInsert(t *testing.T) {
	insert := &InsertQuery{}
	insert.Into("users").
	Set([]string{"Name", "Age"}, []interface{}{"Alice", 25})

	table := utils.FilterData(insert.Table)

	result := insert.Build(table)
	fmt.Println("\n Résultat insert :" + result)
}
