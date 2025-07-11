package models

type Post struct {
	ID      int    `orm:"pk,auto"`
	Title   string `orm:"index"`
	Content string
	UserID  int `orm:"fk:users.id"`
}
