package orm

import (
	"fmt"
	"testing"
	"project/orm/utils"
)

func TestQuery(t *testing.T) {
	query := &Query{}
	query.Select("ID", "Name", "Age").
		From("users").
		Where("age > 29").
		Limit(4)

	table := utils.FilterData(query.Table)

	filteredUsers := query.Apply(table)
	fmt.Print("\n Filtre: ")
	fmt.Print(filteredUsers)
}
