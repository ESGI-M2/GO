package data

import "github.com/ESGI-M2/GO/models"

var Store = map[string]interface{}{
	"users": &Users,
	"posts": &Posts,
}

var Users = []models.User{
	{ID: 1, Name: "Alice", Age: 25, IsActive: true},
	{ID: 2, Name: "Bob", Age: 35, IsActive: true},
	{ID: 3, Name: "Charlie", Age: 30, IsActive: false},
	{ID: 4, Name: "David", Age: 40, IsActive: true},
}

var Posts = []models.Post{
	{ID: 1, Title: "Titre", Content: "Contenu"},
	{ID: 2, Title: "Titre2", Content: "Contenu2"},
}
