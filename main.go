package main

import (
	"fmt"
	"log"
	"os"

	"project/dialect"
	"project/models"
	"project/orm/core"
)

func main() {
	// Create MySQL dialect
	mysqlDialect := dialect.NewMySQLDialect()

	// Create ORM instance
	orm := core.NewORM(mysqlDialect)

	// Connect to database
	config := core.ConnectionConfig{
		Driver:          "mysql",
		Host:            os.Getenv("MYSQL_HOST"),
		Port:            3306,
		Database:        os.Getenv("MYSQL_DATABASE"),
		Username:        os.Getenv("MYSQL_USER"),
		Password:        os.Getenv("MYSQL_PASSWORD"),
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 300,
	}

	if err := orm.Connect(config); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer orm.Close()

	// Register models
	if err := orm.RegisterModel(&models.User{}); err != nil {
		log.Fatalf("Failed to register User model: %v", err)
	}

	if err := orm.RegisterModel(&models.Post{}); err != nil {
		log.Fatalf("Failed to register Post model: %v", err)
	}

	// Create tables
	if err := orm.Migrate(); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	fmt.Println("✅ Database tables created successfully!")

	// Test CRUD operations
	testUserOperations(orm)
	testQueryBuilder(orm)
	testRepositoryPattern(orm)
	testTransactions(orm)
}

func testUserOperations(orm core.ORM) {
	fmt.Println("\n🧪 Testing User CRUD operations...")

	// Create a new user
	user := &models.User{
		Name:     "John Doe",
		Age:      30,
		IsActive: true,
	}

	// Save user (insert)
	if err := orm.Repository(&models.User{}).Save(user); err != nil {
		log.Printf("Failed to save user: %v", err)
		return
	}

	fmt.Printf("✅ User created with ID: %d\n", user.ID)

	// Find user by ID
	foundUser, err := orm.Repository(&models.User{}).Find(user.ID)
	if err != nil {
		log.Printf("Failed to find user: %v", err)
		return
	}

	if foundUser != nil {
		userPtr := foundUser.(*models.User)
		fmt.Printf("✅ Found user: %s (Age: %d, Active: %t)\n",
			userPtr.Name, userPtr.Age, userPtr.IsActive)
	}

	// Update user
	user.Age = 31
	if err := orm.Repository(&models.User{}).Update(user); err != nil {
		log.Printf("Failed to update user: %v", err)
		return
	}

	fmt.Println("✅ User updated successfully")

	// Find all users
	allUsers, err := orm.Repository(&models.User{}).FindAll()
	if err != nil {
		log.Printf("Failed to find all users: %v", err)
		return
	}

	fmt.Printf("✅ Found %d users\n", len(allUsers))
}

func testQueryBuilder(orm core.ORM) {
	fmt.Println("\n🧪 Testing Query Builder...")

	// Query with conditions
	results, err := orm.Query(&models.User{}).
		Where("age", ">", 25).
		Where("is_active", "=", true).
		OrderBy("name", "ASC").
		Limit(10).
		Find()

	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return
	}

	fmt.Printf("✅ Query returned %d results\n", len(results))

	// Raw SQL query
	rawResults, err := orm.Raw("SELECT COUNT(*) as count FROM users WHERE age > ?", 25).Find()
	if err != nil {
		log.Printf("Failed to execute raw query: %v", err)
		return
	}

	if len(rawResults) > 0 {
		fmt.Printf("✅ Raw query count: %v\n", rawResults[0]["count"])
	}
}

func testRepositoryPattern(orm core.ORM) {
	fmt.Println("\n🧪 Testing Repository Pattern...")

	repo := orm.Repository(&models.User{})

	// Find by criteria
	users, err := repo.FindBy(map[string]interface{}{
		"is_active": true,
		"age":       30,
	})

	if err != nil {
		log.Printf("Failed to find users by criteria: %v", err)
		return
	}

	fmt.Printf("✅ Repository found %d users matching criteria\n", len(users))

	// Count all users
	count, err := repo.Count()
	if err != nil {
		log.Printf("Failed to count users: %v", err)
		return
	}

	fmt.Printf("✅ Total users: %d\n", count)
}

func testTransactions(orm core.ORM) {
	fmt.Println("\n🧪 Testing Transactions...")

	err := orm.Transaction(func(txORM core.ORM) error {
		// Create a user within transaction
		user := &models.User{
			Name:     "Transaction User",
			Age:      25,
			IsActive: true,
		}

		if err := txORM.Repository(&models.User{}).Save(user); err != nil {
			return fmt.Errorf("failed to save user in transaction: %w", err)
		}

		// Create a post within the same transaction
		post := &models.Post{
			Title:   "Transaction Post",
			Content: "This post was created in a transaction",
		}

		if err := txORM.Repository(&models.Post{}).Save(post); err != nil {
			return fmt.Errorf("failed to save post in transaction: %w", err)
		}

		fmt.Println("✅ Transaction completed successfully")
		return nil
	})

	if err != nil {
		log.Printf("Transaction failed: %v", err)
		return
	}

	fmt.Println("✅ Transaction test completed")
}
