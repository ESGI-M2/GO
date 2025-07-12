package integration

import "time"

type User struct {
	ID    int    `orm:"primary_key"`
	Name  string `orm:"column:name"`
	Age   int    `orm:"column:age"`
	Email string `orm:"column:email"`
}

type UserWithSoftDelete struct {
	ID        int       `orm:"primary_key"`
	Name      string    `orm:"column:name"`
	Age       int       `orm:"column:age"`
	Email     string    `orm:"column:email"`
	DeletedAt time.Time `orm:"column:deleted_at;soft_delete"`
}

type Post struct {
	ID      int    `orm:"primary_key"`
	Title   string `orm:"column:title"`
	Content string `orm:"column:content"`
	UserID  int    `orm:"column:user_id;foreign_key"`
}
