package data

import "project/models"

var Store = map[string]interface{}{
	"users": &Users,
	"posts": &Posts,
}

var Users = []models.User{
	{ID: 1, Name: "Alice", Age: 25},
	{ID: 2, Name: "Bob", Age: 35},
	{ID: 3, Name: "Charlie", Age: 30},
	{ID: 4, Name: "David", Age: 40},
}

var Posts = []models.Post{
	{ID: 1, Title: "Titre", Content: "Contenu"},
	{ID: 2, Title: "Titre2", Content: "Contenu2"},
}
