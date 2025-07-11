# Go ORM Library

A lightweight, flexible Object-Relational Mapping (ORM) library for Go that supports multiple database dialects and provides a clean, intuitive API for database operations.

## Features

- **Multi-dialect support**: MySQL, PostgreSQL, SQLite (extensible)
- **Simple API**: Easy-to-use repository pattern
- **Query Builder**: Fluent interface for complex queries
- **Transaction Support**: ACID-compliant transactions
- **Raw SQL**: Execute custom SQL when needed
- **Auto-migration**: Automatic table creation from structs
- **Type Safety**: Full Go type safety with reflection

## Installation

```bash
go get github.com/yourusername/go-orm
```

## Quick Start

### 1. Define Your Models

```go
package main

import (
    "time"
    "project/orm"
    "project/dialect"
)

// User represents a user in your system
type User struct {
    Id        int       `orm:"pk,auto"`
    Name      string    `orm:"column:name"`
    Email     string    `orm:"column:email,unique"`
    Age       int       `orm:"column:age"`
    IsActive  bool      `orm:"column:is_active,default:true"`
    CreatedAt time.Time `orm:"column:created_at"`
    UpdatedAt time.Time `orm:"column:updated_at"`
}

// Post represents a blog post
type Post struct {
    Id        int       `orm:"pk,auto"`
    Title     string    `orm:"column:title"`
    Content   string    `orm:"column:content"`
    UserId    int       `orm:"column:user_id"`
    Published bool      `orm:"column:published,default:false"`
    CreatedAt time.Time `orm:"column:created_at"`
    UpdatedAt time.Time `orm:"column:updated_at"`
}
```

### 2. Initialize and Connect

```go
func main() {
    // Initialize with MySQL dialect
    mysqlDialect := dialect.NewMySQLDialect()
    ormInstance := orm.New(mysqlDialect)

    // Configure connection
    config := orm.ConnectionConfig{
        Host:     "localhost",
        Port:     3306,
        Database: "myapp",
        Username: "root",
        Password: "password",
    }

    // Connect to database
    if err := ormInstance.Connect(config); err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer ormInstance.Close()

    // Register models
    ormInstance.RegisterModel(&User{})
    ormInstance.RegisterModel(&Post{})

    // Create tables
    if err := ormInstance.Migrate(); err != nil {
        log.Fatalf("Failed to create tables: %v", err)
    }
}
```

## Usage Examples

### Basic CRUD Operations

```go
// Get repository for User model
userRepo := ormInstance.Repository(&User{})

// Create a new user
user := &User{
    Name:      "John Doe",
    Email:     "john@example.com",
    Age:       30,
    IsActive:  true,
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
}

// Save user
if err := userRepo.Save(user); err != nil {
    log.Printf("Failed to save user: %v", err)
}

// Find user by ID
foundUser, err := userRepo.FindById(1)
if err != nil {
    log.Printf("Failed to find user: %v", err)
}

// Find users by criteria
activeUsers, err := userRepo.FindBy(map[string]interface{}{
    "is_active": true,
    "age": map[string]interface{}{
        ">": 25,
    },
})
if err != nil {
    log.Printf("Failed to find active users: %v", err)
}

// Update user
if foundUser != nil {
    userPtr := foundUser.(*User)
    userPtr.Age = 31
    userPtr.UpdatedAt = time.Now()
    
    if err := userRepo.Update(userPtr); err != nil {
        log.Printf("Failed to update user: %v", err)
    }
}

// Delete user
if err := userRepo.Delete(1); err != nil {
    log.Printf("Failed to delete user: %v", err)
}

// Count users
count, err := userRepo.Count()
if err != nil {
    log.Printf("Failed to count users: %v", err)
}
```

### Advanced Query Builder

```go
// Complex queries with query builder
query := ormInstance.Query(&User{}).
    Where("age", ">", 25).
    Where("is_active", "=", true).
    OrderBy("name", "ASC").
    Limit(10)

results, err := query.Find()
if err != nil {
    log.Printf("Failed to execute query: %v", err)
}

// Count with conditions
count, err := ormInstance.Query(&User{}).
    Where("is_active", "=", true).
    Count()
if err != nil {
    log.Printf("Failed to count: %v", err)
}

// Check existence
exists, err := ormInstance.Query(&User{}).
    Where("email", "=", "john@example.com").
    Exists()
if err != nil {
    log.Printf("Failed to check existence: %v", err)
}
```

### Transactions

```go
// Transaction with multiple operations
err := ormInstance.Transaction(func(txORM orm.ORM) error {
    // Create user within transaction
    user := &User{
        Name:      "Transaction User",
        Email:     "tx@example.com",
        Age:       28,
        IsActive:  true,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    userRepo := txORM.Repository(&User{})
    if err := userRepo.Save(user); err != nil {
        return fmt.Errorf("failed to save user: %w", err)
    }

    // Create post for this user
    post := &Post{
        Title:     "Transaction Post",
        Content:   "Created in transaction",
        UserId:    user.Id,
        Published: true,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    postRepo := txORM.Repository(&Post{})
    if err := postRepo.Save(post); err != nil {
        return fmt.Errorf("failed to save post: %w", err)
    }

    return nil
})

if err != nil {
    log.Printf("Transaction failed: %v", err)
}
```

### Raw SQL Queries

```go
// Execute raw SQL
results, err := ormInstance.Raw(`
    SELECT 
        u.name,
        COUNT(p.id) as post_count
    FROM users u
    LEFT JOIN posts p ON u.id = p.user_id
    WHERE u.is_active = 1
    GROUP BY u.id, u.name
    ORDER BY post_count DESC
`).Find()

if err != nil {
    log.Printf("Failed to execute raw SQL: %v", err)
}

for _, result := range results {
    fmt.Printf("User: %s, Posts: %v\n", 
        result["name"], 
        result["post_count"])
}
```

## Real-World Blog Application Example

Here's a complete example of a blog application using the ORM:

```go
package main

import (
    "fmt"
    "log"
    "time"
    "project/orm"
    "project/dialect"
)

// Blog Models
type BlogUser struct {
    Id        int       `orm:"pk,auto"`
    Name      string    `orm:"column:name"`
    Email     string    `orm:"column:email,unique"`
    Age       int       `orm:"column:age"`
    IsActive  bool      `orm:"column:is_active,default:true"`
    CreatedAt time.Time `orm:"column:created_at"`
    UpdatedAt time.Time `orm:"column:updated_at"`
}

type BlogPost struct {
    Id        int       `orm:"pk,auto"`
    Title     string    `orm:"column:title"`
    Content   string    `orm:"column:content"`
    UserId    int       `orm:"column:user_id"`
    Published bool      `orm:"column:published,default:false"`
    CreatedAt time.Time `orm:"column:created_at"`
    UpdatedAt time.Time `orm:"column:updated_at"`
}

type BlogComment struct {
    Id        int       `orm:"pk,auto"`
    PostId    int       `orm:"column:post_id"`
    UserId    int       `orm:"column:user_id"`
    Content   string    `orm:"column:content"`
    CreatedAt time.Time `orm:"column:created_at"`
}

func main() {
    // Initialize ORM
    mysqlDialect := dialect.NewMySQLDialect()
    ormInstance := orm.New(mysqlDialect)

    // Connect to database
    config := orm.ConnectionConfig{
        Host:     "localhost",
        Port:     3306,
        Database: "blog_db",
        Username: "root",
        Password: "password",
    }

    if err := ormInstance.Connect(config); err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer ormInstance.Close()

    // Register models
    models := []interface{}{
        &BlogUser{},
        &BlogPost{},
        &BlogComment{},
    }

    for _, model := range models {
        ormInstance.RegisterModel(model)
    }

    // Create tables
    if err := ormInstance.Migrate(); err != nil {
        log.Fatalf("Failed to create tables: %v", err)
    }

    // Run blog operations
    runBlogExample(ormInstance)
}

func runBlogExample(orm orm.ORM) {
    userRepo := orm.Repository(&BlogUser{})
    postRepo := orm.Repository(&BlogPost{})
    commentRepo := orm.Repository(&BlogComment{})

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
    }

    for _, user := range users {
        if err := userRepo.Save(user); err != nil {
            log.Printf("Failed to save user: %v", err)
            continue
        }
        fmt.Printf("Created user: %s (ID: %d)\n", user.Name, user.Id)
    }

    // Create posts
    posts := []*BlogPost{
        {
            Title:     "Getting Started with Go",
            Content:   "Go is a powerful programming language...",
            UserId:    1,
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
    }

    for _, post := range posts {
        if err := postRepo.Save(post); err != nil {
            log.Printf("Failed to save post: %v", err)
            continue
        }
        fmt.Printf("Created post: %s (ID: %d)\n", post.Title, post.Id)
    }

    // Create comments
    comments := []*BlogComment{
        {
            PostId:    1,
            UserId:    2,
            Content:   "Great article! Very helpful.",
            CreatedAt: time.Now(),
        },
        {
            PostId:    1,
            UserId:    1,
            Content:   "Thanks for the feedback!",
            CreatedAt: time.Now(),
        },
    }

    for _, comment := range comments {
        if err := commentRepo.Save(comment); err != nil {
            log.Printf("Failed to save comment: %v", err)
            continue
        }
        fmt.Printf("Created comment (ID: %d) on post %d\n", comment.Id, comment.PostId)
    }

    // Query examples
    publishedPosts, err := postRepo.FindBy(map[string]interface{}{
        "published": true,
    })
    if err != nil {
        log.Printf("Failed to find published posts: %v", err)
    } else {
        fmt.Printf("Found %d published posts\n", len(publishedPosts))
    }

    // Advanced query
    userStats, err := orm.Raw(`
        SELECT 
            u.name,
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
        log.Printf("Failed to get user stats: %v", err)
    } else {
        fmt.Println("User Activity Report:")
        for _, stat := range userStats {
            fmt.Printf("  %s: %v posts, %v comments\n",
                stat["name"],
                stat["post_count"],
                stat["comment_count"])
        }
    }
}
```

## Environment Configuration

You can use environment variables for database configuration:

```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_NAME=myapp
export DB_USER=root
export DB_PASSWORD=password
```

```go
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

config := orm.ConnectionConfig{
    Host:     getEnv("DB_HOST", "localhost"),
    Port:     getEnvAsInt("DB_PORT", 3306),
    Database: getEnv("DB_NAME", "myapp"),
    Username: getEnv("DB_USER", "root"),
    Password: getEnv("DB_PASSWORD", "password"),
}
```

## Testing Your Application

### Unit Tests

```go
package main

import (
    "testing"
    "project/orm"
    "project/dialect"
)

func TestUserOperations(t *testing.T) {
    // Use mock dialect for testing
    mockDialect := dialect.NewMockDialect()
    ormInstance := orm.New(mockDialect)

    // Test user creation
    user := &User{
        Name:      "Test User",
        Email:     "test@example.com",
        Age:       25,
        IsActive:  true,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    userRepo := ormInstance.Repository(&User{})
    if err := userRepo.Save(user); err != nil {
        t.Errorf("Failed to save user: %v", err)
    }

    // Test user retrieval
    foundUser, err := userRepo.FindById(user.Id)
    if err != nil {
        t.Errorf("Failed to find user: %v", err)
    }

    if foundUser == nil {
        t.Error("User not found")
    }
}
```

### Integration Tests

```go
func TestBlogIntegration(t *testing.T) {
    // Use real database for integration tests
    mysqlDialect := dialect.NewMySQLDialect()
    ormInstance := orm.New(mysqlDialect)

    config := orm.ConnectionConfig{
        Host:     "localhost",
        Port:     3306,
        Database: "test_blog_db",
        Username: "root",
        Password: "password",
    }

    if err := ormInstance.Connect(config); err != nil {
        t.Fatalf("Failed to connect: %v", err)
    }
    defer ormInstance.Close()

    // Register models
    ormInstance.RegisterModel(&BlogUser{})
    ormInstance.RegisterModel(&BlogPost{})

    // Create tables
    if err := ormInstance.Migrate(); err != nil {
        t.Fatalf("Failed to create tables: %v", err)
    }

    // Test complete blog workflow
    userRepo := ormInstance.Repository(&BlogUser{})
    postRepo := ormInstance.Repository(&BlogPost{})

    // Create user
    user := &BlogUser{
        Name:      "Test User",
        Email:     "test@example.com",
        Age:       25,
        IsActive:  true,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    if err := userRepo.Save(user); err != nil {
        t.Errorf("Failed to save user: %v", err)
    }

    // Create post
    post := &BlogPost{
        Title:     "Test Post",
        Content:   "Test content",
        UserId:    user.Id,
        Published: true,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    if err := postRepo.Save(post); err != nil {
        t.Errorf("Failed to save post: %v", err)
    }

    // Verify post was created
    foundPost, err := postRepo.FindById(post.Id)
    if err != nil {
        t.Errorf("Failed to find post: %v", err)
    }

    if foundPost == nil {
        t.Error("Post not found")
    }
}
```

## Best Practices

1. **Use Transactions**: Wrap related operations in transactions for data consistency
2. **Handle Errors**: Always check for errors and handle them appropriately
3. **Use Environment Variables**: Configure database connections via environment variables
4. **Test Thoroughly**: Write both unit and integration tests
5. **Close Connections**: Always defer `ormInstance.Close()` to clean up resources
6. **Use Proper Field Tags**: Define your struct tags carefully for proper mapping

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.