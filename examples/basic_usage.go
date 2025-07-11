package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"project/dialect"
	"project/orm/core"
)

// Example models with various features
type User struct {
	ID        int       `db:"id" primary:"true" autoincrement:"true"`
	Name      string    `db:"name"`
	Email     string    `db:"email" unique:"true"`
	Age       int       `db:"age"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Post struct {
	ID        int       `db:"id" primary:"true" autoincrement:"true"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	UserID    int       `db:"user_id" foreign:"users.id"`
	Published bool      `db:"published"`
	CreatedAt time.Time `db:"created_at"`
}

type Category struct {
	ID          int    `db:"id" primary:"true" autoincrement:"true"`
	Name        string `db:"name" unique:"true"`
	Description string `db:"description"`
}

func main() {
	// Initialize MySQL dialect
	mysqlDialect := dialect.NewMySQLDialect()

	// Create ORM instance
	orm := core.NewORM(mysqlDialect)

	// Configure database connection
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

	// Connect to database
	if err := orm.Connect(config); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer orm.Close()

	fmt.Println("✅ Connected to database successfully")

	// Register models
	models := []interface{}{
		&User{},
		&Post{},
		&Category{},
	}

	for _, model := range models {
		if err := orm.RegisterModel(model); err != nil {
			log.Fatalf("Failed to register model: %v", err)
		}
	}

	fmt.Println("✅ Models registered successfully")

	// Create tables
	if err := orm.Migrate(); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	fmt.Println("✅ Database tables created successfully")

	// Example 1: Basic CRUD operations
	exampleBasicCRUD(orm)

	// Example 2: Query Builder
	exampleQueryBuilder(orm)

	// Example 3: Repository Pattern
	exampleRepositoryPattern(orm)

	// Example 4: Transactions
	exampleTransactions(orm)

	// Example 5: Raw SQL
	exampleRawSQL(orm)

	fmt.Println("\n🎉 All examples completed successfully!")
}

func exampleBasicCRUD(orm core.ORM) {
	fmt.Println("\n📝 Example 1: Basic CRUD Operations")

	// Create a new user
	user := &User{
		Name:      "John Doe",
		Email:     "john@example.com",
		Age:       30,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save user (insert)
	if err := orm.Repository(&User{}).Save(user); err != nil {
		log.Printf("Failed to save user: %v", err)
		return
	}

	fmt.Printf("✅ User created with ID: %d\n", user.ID)

	// Find user by ID
	foundUser, err := orm.Repository(&User{}).Find(user.ID)
	if err != nil {
		log.Printf("Failed to find user: %v", err)
		return
	}

	if foundUser != nil {
		userPtr := foundUser.(*User)
		fmt.Printf("✅ Found user: %s (%s)\n", userPtr.Name, userPtr.Email)
	}

	// Update user
	user.Age = 31
	user.UpdatedAt = time.Now()

	if err := orm.Repository(&User{}).Update(user); err != nil {
		log.Printf("Failed to update user: %v", err)
		return
	}

	fmt.Println("✅ User updated successfully")

	// Find all users
	allUsers, err := orm.Repository(&User{}).FindAll()
	if err != nil {
		log.Printf("Failed to find all users: %v", err)
		return
	}

	fmt.Printf("✅ Found %d users\n", len(allUsers))
}

func exampleQueryBuilder(orm core.ORM) {
	fmt.Println("\n🔍 Example 2: Query Builder")

	// Create some test data
	users := []*User{
		{Name: "Alice", Email: "alice@example.com", Age: 25, IsActive: true},
		{Name: "Bob", Email: "bob@example.com", Age: 35, IsActive: true},
		{Name: "Charlie", Email: "charlie@example.com", Age: 28, IsActive: false},
	}

	for _, user := range users {
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		orm.Repository(&User{}).Save(user)
	}

	// Query with conditions
	results, err := orm.Query(&User{}).
		Where("age", ">", 25).
		Where("is_active", "=", true).
		OrderBy("name", "ASC").
		Limit(10).
		Find()

	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return
	}

	fmt.Printf("✅ Query returned %d active users over 25\n", len(results))

	// Count query
	count, err := orm.Query(&User{}).
		Where("is_active", "=", true).
		Count()

	if err != nil {
		log.Printf("Failed to count users: %v", err)
		return
	}

	fmt.Printf("✅ Found %d active users\n", count)
}

func exampleRepositoryPattern(orm core.ORM) {
	fmt.Println("\n🏪 Example 3: Repository Pattern")

	repo := orm.Repository(&User{})

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

	// Find one by criteria
	user, err := repo.FindOneBy(map[string]interface{}{
		"email": "john@example.com",
	})

	if err != nil {
		log.Printf("Failed to find user by email: %v", err)
		return
	}

	if user != nil {
		userPtr := user.(*User)
		fmt.Printf("✅ Found user by email: %s\n", userPtr.Name)
	}

	// Count all users
	count, err := repo.Count()
	if err != nil {
		log.Printf("Failed to count users: %v", err)
		return
	}

	fmt.Printf("✅ Total users: %d\n", count)
}

func exampleTransactions(orm core.ORM) {
	fmt.Println("\n💼 Example 4: Transactions")

	err := orm.Transaction(func(txORM core.ORM) error {
		// Create a user within transaction
		user := &User{
			Name:      "Transaction User",
			Email:     "tx@example.com",
			Age:       25,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := txORM.Repository(&User{}).Save(user); err != nil {
			return fmt.Errorf("failed to save user in transaction: %w", err)
		}

		// Create a post within the same transaction
		post := &Post{
			Title:     "Transaction Post",
			Content:   "This post was created in a transaction",
			UserID:    user.ID,
			Published: true,
			CreatedAt: time.Now(),
		}

		if err := txORM.Repository(&Post{}).Save(post); err != nil {
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

func exampleRawSQL(orm core.ORM) {
	fmt.Println("\n🔧 Example 5: Raw SQL")

	// Raw SQL query
	results, err := orm.Raw("SELECT COUNT(*) as count FROM users WHERE age > ?", 25).Find()
	if err != nil {
		log.Printf("Failed to execute raw query: %v", err)
		return
	}

	if len(results) > 0 {
		fmt.Printf("✅ Raw query count: %v\n", results[0]["count"])
	}

	// Raw SQL with complex query
	complexResults, err := orm.Raw(`
		SELECT u.name, COUNT(p.id) as post_count 
		FROM users u 
		LEFT JOIN posts p ON u.id = p.user_id 
		WHERE u.is_active = ? 
		GROUP BY u.id, u.name
	`, true).Find()

	if err != nil {
		log.Printf("Failed to execute complex raw query: %v", err)
		return
	}

	fmt.Printf("✅ Complex raw query returned %d results\n", len(complexResults))
}
