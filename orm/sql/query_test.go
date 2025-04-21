package orm

import (
	"fmt"
	"testing"
	"project/orm/utils"
	"project/memory"
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
		Set([]string{"Name", "Age"}, []interface{}{"Test", 20})

	insert.Apply()

	for _, user := range data.Users {
		fmt.Printf("User: %+v\n", user)
	}
}