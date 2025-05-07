package orm

import (
	"fmt"
	"strings"
	"log"
	"project/orm/utils"
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
	entityType := utils.GetEntityType(iq.Table)
    object := utils.SetFields(entityType, iq.Columns, iq.Values)
	utils.AppendToData(iq.Table, object)
}

