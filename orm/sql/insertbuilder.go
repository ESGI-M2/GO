package orm

import (
	"fmt"
	"strings"
	"log"
	"project/models"
	"project/memory"
)

func (iq *InsertQuery) Build() string {
	if iq.Table == "" || len(iq.Columns) == 0 || len(iq.Values) == 0 {
		log.Fatalf("Incomplete insert query")
	}

	columns := strings.Join(iq.Columns, ", ")
	placeholders := make([]string, len(iq.Values))
	for i := range iq.Values {
		placeholders[i] = "?"
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", iq.Table, columns, strings.Join(placeholders, ", "))
}

func (iq *InsertQuery) Apply() {
	switch iq.Table {
	case "users":
		var newUser models.User
		for i, col := range iq.Columns {
			switch strings.ToLower(col) {
			case "name":
				newUser.Name = iq.Values[i].(string)
			case "age":
				newUser.Age = iq.Values[i].(int)
			}
		}

		maxID := 0
		for _, u := range data.Users {
			if u.ID > maxID {
				maxID = u.ID
			}
		}
		newUser.ID = maxID + 1
		data.Users = append(data.Users, newUser)
	default:
		log.Fatalf("Unknown table: %s", iq.Table)
	}
}

