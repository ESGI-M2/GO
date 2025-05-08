package models

type User struct {
	ID   int    `db:"id" primary:"true" autoincrement:"true"`
	Name string `db:"name"`
	Age  int    `db:"age"`
	IsActive bool
}
