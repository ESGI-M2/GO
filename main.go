package main

import (
	"fmt"
	"log"

	"project/dialect"
	"project/orm"
)

// User model with new ORM tags
type User struct {
	ID       int    `orm:"pk,auto"`
	Name     string `orm:"index"`
	Email    string `orm:"unique"`
	Age      int
	IsActive bool `orm:"default:true"`
}

func main() {
	// Create MySQL dialect
	mysqlDialect := dialect.NewMySQLDialect()

	// Create ORM instance
	ormInstance := orm.New(mysqlDialect)

	// Connect to database
	config := orm.ConnectionConfig{
		Host:     "localhost",
		Port:     3306,
		Database: "test_db",
		Username: "root",
		Password: "password",
	}

	err := ormInstance.Connect(config)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer ormInstance.Close()

	// Register model
	err = ormInstance.RegisterModel(&User{})
	if err != nil {
		log.Fatalf("Failed to register model: %v", err)
	}

	// Create table
	err = ormInstance.CreateTable(&User{})
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Create user
	user := &User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Age:      30,
		IsActive: true,
	}

	// Save user
	repo := ormInstance.Repository(&User{})
	err = repo.Save(user)
	if err != nil {
		log.Fatalf("Failed to save user: %v", err)
	}

	fmt.Printf("User saved with ID: %d\n", user.ID)

	// Find user by ID
	foundUser, err := repo.Find(user.ID)
	if err != nil {
		log.Fatalf("Failed to find user: %v", err)
	}

	if foundUser != nil {
		fmt.Printf("Found user: %+v\n", foundUser)
	}

	// Query builder example
	query := ormInstance.Query(&User{}).
		Where("age", ">", 25).
		OrderBy("name", "ASC").
		Limit(10)

	results, err := query.Find()
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}

	fmt.Printf("Found %d users over 25\n", len(results))

	// Transaction example
	err = ormInstance.Transaction(func(txORM orm.ORM) error {
		user2 := &User{
			Name:     "Jane Doe",
			Email:    "jane@example.com",
			Age:      25,
			IsActive: true,
		}

		repo := txORM.Repository(user2)
		return repo.Save(user2)
	})

	if err != nil {
		log.Fatalf("Transaction failed: %v", err)
	}

	fmt.Println("Transaction completed successfully")
}
