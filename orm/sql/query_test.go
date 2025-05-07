package orm

import (
	"fmt"
	"testing"
	"project/memory"
)

func TestQuery(t *testing.T) {
	query := &Query{}
	query.Select("*").
		From("users").
		Where("age > 0").
		Limit(2)

	filteredUsers := query.Apply(query.Table)
	fmt.Print("\n Filtre: ")
	fmt.Print(filteredUsers, "\n")
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