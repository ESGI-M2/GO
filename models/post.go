package models

type Post struct {
	ID      int    `db:"id" primary:"true" autoincrement:"true"`
	Title   string `db:"title"`
	Content string `db:"content"`
}
