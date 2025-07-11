package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"project/dialect"
	"project/orm"
)

// BlogUser represents a user in the blog system
type BlogUser struct {
	Id        int       `orm:"pk,auto"`
	Name      string    `orm:"column:name"`
	Email     string    `orm:"column:email,unique"`
	Age       int       `orm:"column:age"`
	IsActive  bool      `orm:"column:is_active,default:true"`
	CreatedAt time.Time `orm:"column:created_at"`
	UpdatedAt time.Time `orm:"column:updated_at"`
}

// BlogPost represents a blog post
type BlogPost struct {
	Id        int       `orm:"pk,auto"`
	Title     string    `orm:"column:title"`
	Content   string    `orm:"column:content"`
	UserId    int       `orm:"column:user_id"`
	Published bool      `orm:"column:published,default:false"`
	CreatedAt time.Time `orm:"column:created_at"`
	UpdatedAt time.Time `orm:"column:updated_at"`
}

// BlogCategory represents a post category
type BlogCategory struct {
	Id          int    `orm:"pk,auto"`
	Name        string `orm:"column:name,unique"`
	Description string `orm:"column:description"`
}

// BlogComment represents a comment on a post
type BlogComment struct {
	Id        int       `orm:"pk,auto"`
	PostId    int       `orm:"column:post_id"`
	UserId    int       `orm:"column:user_id"`
	Content   string    `orm:"column:content"`
	CreatedAt time.Time `orm:"column:created_at"`
}

func main() {
	// Initialize the ORM with MySQL dialect
	mysqlDialect := dialect.NewMySQLDialect()
	ormInstance := orm.New(mysqlDialect)

	// Configure database connection
	config := orm.ConnectionConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvAsInt("DB_PORT", 3306),
		Database: getEnv("DB_NAME", "blog_db"),
		Username: getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "password"),
	}

	// Connect to database
	if err := ormInstance.Connect(config); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer ormInstance.Close()

	fmt.Println("✅ Connected to database successfully")

	// Register all models
	models := []interface{}{
		&BlogUser{},
		&BlogPost{},
		&BlogCategory{},
		&BlogComment{},
	}

	for _, model := range models {
		if err := ormInstance.RegisterModel(model); err != nil {
			log.Fatalf("Failed to register model: %v", err)
		}
	}

	fmt.Println("✅ Models registered successfully")

	// Create tables
	if err := ormInstance.Migrate(); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	fmt.Println("✅ Database tables created successfully")

	// Run different examples
	runUserManagement(ormInstance)
	runPostManagement(ormInstance)
	runCategoryManagement(ormInstance)
	runCommentManagement(ormInstance)
	runAdvancedQueries(ormInstance)
	runTransactions(ormInstance)
	runRawQueries(ormInstance)

	fmt.Println("\n🎉 All examples completed successfully!")
}

func runUserManagement(orm orm.ORM) {
	fmt.Println("\n👥 User Management Example")

	userRepo := orm.Repository(&BlogUser{})

	// Create users
	users := []*BlogUser{
		{
			Name:      "John Doe",
			Email:     "john@example.com",
			Age:       30,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Jane Smith",
			Email:     "jane@example.com",
			Age:       25,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Bob Johnson",
			Email:     "bob@example.com",
			Age:       35,
			IsActive:  false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Save users
	for _, user := range users {
		if err := userRepo.Save(user); err != nil {
			log.Printf("Failed to save user %s: %v", user.Name, err)
			continue
		}
		fmt.Printf("✅ Created user: %s (ID: %d)\n", user.Name, user.Id)
	}

	// Find all active users
	activeUsers, err := userRepo.FindBy(map[string]interface{}{
		"is_active": true,
	})
	if err != nil {
		log.Printf("Failed to find active users: %v", err)
		return
	}

	fmt.Printf("✅ Found %d active users\n", len(activeUsers))

	// Find user by email
	user, err := userRepo.FindOneBy(map[string]interface{}{
		"email": "john@example.com",
	})
	if err != nil {
		log.Printf("Failed to find user by email: %v", err)
		return
	}

	if user != nil {
		userPtr := user.(*BlogUser)
		fmt.Printf("✅ Found user by email: %s (Age: %d)\n", userPtr.Name, userPtr.Age)
	}

	// Update user
	if user != nil {
		userPtr := user.(*BlogUser)
		userPtr.Age = 31
		userPtr.UpdatedAt = time.Now()

		if err := userRepo.Update(userPtr); err != nil {
			log.Printf("Failed to update user: %v", err)
			return
		}
		fmt.Printf("✅ Updated user: %s (New age: %d)\n", userPtr.Name, userPtr.Age)
	}

	// Count total users
	count, err := userRepo.Count()
	if err != nil {
		log.Printf("Failed to count users: %v", err)
		return
	}
	fmt.Printf("✅ Total users: %d\n", count)
}

func runPostManagement(orm orm.ORM) {
	fmt.Println("\n📝 Post Management Example")

	postRepo := orm.Repository(&BlogPost{})

	// Create posts
	posts := []*BlogPost{
		{
			Title:     "Getting Started with Go",
			Content:   "Go is a powerful programming language...",
			UserId:    1, // Assuming user with ID 1 exists
			Published: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Title:     "Advanced ORM Patterns",
			Content:   "This post covers advanced ORM usage...",
			UserId:    1,
			Published: false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Title:     "Database Design Best Practices",
			Content:   "Learn about database design principles...",
			UserId:    2, // Assuming user with ID 2 exists
			Published: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Save posts
	for _, post := range posts {
		if err := postRepo.Save(post); err != nil {
			log.Printf("Failed to save post %s: %v", post.Title, err)
			continue
		}
		fmt.Printf("✅ Created post: %s (ID: %d)\n", post.Title, post.Id)
	}

	// Find published posts
	publishedPosts, err := postRepo.FindBy(map[string]interface{}{
		"published": true,
	})
	if err != nil {
		log.Printf("Failed to find published posts: %v", err)
		return
	}

	fmt.Printf("✅ Found %d published posts\n", len(publishedPosts))

	// Find posts by user
	userPosts, err := postRepo.FindBy(map[string]interface{}{
		"user_id": 1,
	})
	if err != nil {
		log.Printf("Failed to find user posts: %v", err)
		return
	}

	fmt.Printf("✅ User 1 has %d posts\n", len(userPosts))
}

func runCategoryManagement(orm orm.ORM) {
	fmt.Println("\n📂 Category Management Example")

	categoryRepo := orm.Repository(&BlogCategory{})

	// Create categories
	categories := []*BlogCategory{
		{
			Name:        "Programming",
			Description: "Articles about programming languages and techniques",
		},
		{
			Name:        "Database",
			Description: "Database design and optimization articles",
		},
		{
			Name:        "DevOps",
			Description: "DevOps practices and tools",
		},
	}

	// Save categories
	for _, category := range categories {
		if err := categoryRepo.Save(category); err != nil {
			log.Printf("Failed to save category %s: %v", category.Name, err)
			continue
		}
		fmt.Printf("✅ Created category: %s (ID: %d)\n", category.Name, category.Id)
	}

	// Find category by name
	category, err := categoryRepo.FindOneBy(map[string]interface{}{
		"name": "Programming",
	})
	if err != nil {
		log.Printf("Failed to find category: %v", err)
		return
	}

	if category != nil {
		categoryPtr := category.(*BlogCategory)
		fmt.Printf("✅ Found category: %s - %s\n", categoryPtr.Name, categoryPtr.Description)
	}
}

func runCommentManagement(orm orm.ORM) {
	fmt.Println("\n💬 Comment Management Example")

	commentRepo := orm.Repository(&BlogComment{})

	// Create comments
	comments := []*BlogComment{
		{
			PostId:    1, // Assuming post with ID 1 exists
			UserId:    1, // Assuming user with ID 1 exists
			Content:   "Great article! Very helpful.",
			CreatedAt: time.Now(),
		},
		{
			PostId:    1,
			UserId:    2, // Assuming user with ID 2 exists
			Content:   "I learned a lot from this post.",
			CreatedAt: time.Now(),
		},
		{
			PostId:    3, // Assuming post with ID 3 exists
			UserId:    1,
			Content:   "Excellent database design tips!",
			CreatedAt: time.Now(),
		},
	}

	// Save comments
	for _, comment := range comments {
		if err := commentRepo.Save(comment); err != nil {
			log.Printf("Failed to save comment: %v", err)
			continue
		}
		fmt.Printf("✅ Created comment (ID: %d) on post %d\n", comment.Id, comment.PostId)
	}

	// Find comments for a specific post
	postComments, err := commentRepo.FindBy(map[string]interface{}{
		"post_id": 1,
	})
	if err != nil {
		log.Printf("Failed to find post comments: %v", err)
		return
	}

	fmt.Printf("✅ Post 1 has %d comments\n", len(postComments))
}

func runAdvancedQueries(orm orm.ORM) {
	fmt.Println("\n🔍 Advanced Query Examples")

	// Query builder with complex conditions
	query := orm.Query(&BlogUser{}).
		Where("age", ">", 25).
		Where("is_active", "=", true).
		OrderBy("name", "ASC").
		Limit(10)

	results, err := query.Find()
	if err != nil {
		log.Printf("Failed to execute advanced query: %v", err)
		return
	}

	fmt.Printf("✅ Advanced query returned %d users over 25\n", len(results))

	// Count query
	count, err := orm.Query(&BlogUser{}).
		Where("is_active", "=", true).
		Count()
	if err != nil {
		log.Printf("Failed to count active users: %v", err)
		return
	}

	fmt.Printf("✅ Active users count: %d\n", count)

	// Exists query
	exists, err := orm.Query(&BlogUser{}).
		Where("email", "=", "john@example.com").
		Exists()
	if err != nil {
		log.Printf("Failed to check user existence: %v", err)
		return
	}

	fmt.Printf("✅ User with email 'john@example.com' exists: %t\n", exists)
}

func runTransactions(orm orm.ORM) {
	fmt.Println("\n💼 Transaction Examples")

	// Transaction with multiple operations
	err := orm.Transaction(func(txORM orm.ORM) error {
		// Create a new user within transaction
		user := &BlogUser{
			Name:      "Transaction User",
			Email:     "tx@example.com",
			Age:       28,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		userRepo := txORM.Repository(&BlogUser{})
		if err := userRepo.Save(user); err != nil {
			return fmt.Errorf("failed to save user in transaction: %w", err)
		}

		// Create a post for this user within the same transaction
		post := &BlogPost{
			Title:     "Transaction Post",
			Content:   "This post was created in a transaction",
			UserId:    user.Id,
			Published: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		postRepo := txORM.Repository(&BlogPost{})
		if err := postRepo.Save(post); err != nil {
			return fmt.Errorf("failed to save post in transaction: %w", err)
		}

		fmt.Printf("✅ Transaction: Created user %s and post '%s'\n", user.Name, post.Title)
		return nil
	})

	if err != nil {
		log.Printf("Transaction failed: %v", err)
		return
	}

	fmt.Println("✅ Transaction completed successfully")
}

func runRawQueries(orm orm.ORM) {
	fmt.Println("\n🔧 Raw SQL Examples")

	// Raw SQL query to get user statistics
	results, err := orm.Raw(`
		SELECT 
			COUNT(*) as total_users,
			COUNT(CASE WHEN is_active = 1 THEN 1 END) as active_users,
			AVG(age) as avg_age
		FROM blog_users
	`).Find()

	if err != nil {
		log.Printf("Failed to execute raw SQL: %v", err)
		return
	}

	if len(results) > 0 {
		stats := results[0]
		fmt.Printf("✅ User Statistics:\n")
		fmt.Printf("   Total Users: %v\n", stats["total_users"])
		fmt.Printf("   Active Users: %v\n", stats["active_users"])
		fmt.Printf("   Average Age: %.2f\n", stats["avg_age"])
	}

	// Raw SQL with complex join
	complexResults, err := orm.Raw(`
		SELECT 
			u.name as user_name,
			COUNT(p.id) as post_count,
			COUNT(c.id) as comment_count
		FROM blog_users u
		LEFT JOIN blog_posts p ON u.id = p.user_id
		LEFT JOIN blog_comments c ON u.id = c.user_id
		WHERE u.is_active = 1
		GROUP BY u.id, u.name
		ORDER BY post_count DESC
	`).Find()

	if err != nil {
		log.Printf("Failed to execute complex raw SQL: %v", err)
		return
	}

	fmt.Printf("✅ User Activity Report:\n")
	for _, result := range complexResults {
		fmt.Printf("   %s: %v posts, %v comments\n",
			result["user_name"],
			result["post_count"],
			result["comment_count"])
	}
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := fmt.Sscanf(value, "%d", &defaultValue); err == nil && intValue > 0 {
			return defaultValue
		}
	}
	return defaultValue
}
