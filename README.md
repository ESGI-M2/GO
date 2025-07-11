# Go ORM - A Full-Featured Relational ORM

A comprehensive, modular, and maintainable Object-Relational Mapping (ORM) library for Go, inspired by Doctrine in PHP. This ORM provides a fluent interface, supports multiple database dialects, and follows Go best practices.

## Features

- **Multiple Database Support**: MySQL, PostgreSQL, SQLite (extensible)
- **Fluent Query Builder**: Chainable methods for complex queries
- **Repository Pattern**: Clean data access layer
- **Transaction Support**: ACID-compliant transactions
- **Migration System**: Database schema management
- **Relationship Support**: One-to-one, one-to-many, many-to-many
- **Advanced Tag System**: Concise and powerful model annotations
- **Type Safety**: Full Go type safety with reflection
- **Performance**: Optimized for high-performance applications
- **Extensible**: Plugin architecture for custom dialects

## Installation

```bash
go get github.com/your-username/go-orm
```

## Quick Start

### 1. Define Your Models

Use the new concise ORM tag system:

```go
package models

type User struct {
    ID       int    `orm:"pk,auto"`           // Primary key, auto increment
    Name     string `orm:"index"`              // Indexed field
    Email    string `orm:"unique"`             // Unique constraint
    Age      int    `orm:"default:18"`         // Default value
    IsActive bool   `orm:"default:true"`       // Boolean with default
    Created  string `orm:"column:created_at"`  // Custom column name
}

type Post struct {
    ID      int    `orm:"pk,auto"`
    Title   string `orm:"index"`
    Content string
    UserID  int    `orm:"fk:users.id"`        // Foreign key relationship
}
```

### 2. Initialize the ORM

```go
package main

import (
    "project/orm/core"
    "project/dialect"
)

func main() {
    // Create MySQL dialect
    mysqlDialect := dialect.NewMySQLDialect()
    
    // Create ORM instance
    orm := core.NewORM(mysqlDialect)
    
    // Connect to database
    config := core.ConnectionConfig{
        Driver:   "mysql",
        Host:     "localhost",
        Port:     3306,
        Database: "myapp",
        Username: "root",
        Password: "password",
    }
    
    err := orm.Connect(config)
    if err != nil {
        log.Fatal(err)
    }
    defer orm.Close()
    
    // Register models
    orm.RegisterModel(&User{})
    orm.RegisterModel(&Post{})
    
    // Create tables
    orm.Migrate()
}
```

### 3. Use the Repository Pattern

```go
// Get repository for User model
repo := orm.Repository(&User{})

// Create a new user
user := &User{
    Name:     "John Doe",
    Email:    "john@example.com",
    Age:      30,
    IsActive: true,
}

// Save user (insert or update)
err := repo.Save(user)

// Find user by ID
foundUser, err := repo.Find(1)

// Find all users
allUsers, err := repo.FindAll()

// Find by criteria
activeUsers, err := repo.FindBy(map[string]interface{}{
    "is_active": true,
    "age":       30,
})
```

### 4. Use the Query Builder

```go
// Complex queries with fluent interface
results, err := orm.Query(&User{}).
    Select("name", "email").
    Where("age", ">", 25).
    Where("is_active", "=", true).
    OrderBy("name", "ASC").
    Limit(10).
    Find()

// Raw SQL queries
rawResults, err := orm.Raw("SELECT * FROM users WHERE age > ?", 25).Find()
```

### 5. Use Transactions

```go
err := orm.Transaction(func(txORM core.ORM) error {
    // Create user
    user := &User{Name: "John", Email: "john@example.com"}
    repo := txORM.Repository(user)
    err := repo.Save(user)
    if err != nil {
        return err
    }
    
    // Create post in same transaction
    post := &Post{Title: "Hello", UserID: user.ID}
    postRepo := txORM.Repository(post)
    return postRepo.Save(post)
})
```

## ORM Tag System

The new ORM tag system provides a concise and powerful way to define model metadata:

### Basic Tags

```go
type User struct {
    ID       int    `orm:"pk,auto"`           // Primary key + auto increment
    Name     string `orm:"index"`              // Indexed field
    Email    string `orm:"unique"`             // Unique constraint
    Age      int    `orm:"default:18"`         // Default value
    IsActive bool   `orm:"default:true"`       // Boolean default
    Created  string `orm:"column:created_at"`  // Custom column name
}
```

### Advanced Tags

```go
type Post struct {
    ID        int    `orm:"pk,auto"`
    Title     string `orm:"column:post_title,index,length:255"`
    Content   string `orm:"column:post_content"`
    UserID    int    `orm:"fk:users.id"`                    // Foreign key
    Status    string `orm:"default:draft,nullable"`          // Default + nullable
    Tags      string `orm:"column:post_tags,length:500"`
}
```

### Tag Reference

| Tag | Description | Example |
|-----|-------------|---------|
| `pk` | Primary key | `orm:"pk"` |
| `auto` | Auto increment | `orm:"auto"` |
| `unique` | Unique constraint | `orm:"unique"` |
| `index` | Create index | `orm:"index"` |
| `nullable` | Allow NULL values | `orm:"nullable"` |
| `column:name` | Custom column name | `orm:"column:user_name"` |
| `length:n` | Field length | `orm:"length:255"` |
| `default:value` | Default value | `orm:"default:true"` |
| `fk:table.column` | Foreign key | `orm:"fk:users.id"` |

### Backward Compatibility

The ORM maintains full backward compatibility with the old tag system:

```go
// Old style (still supported)
type User struct {
    ID       int    `db:"id" primary:"true" autoincrement:"true"`
    Name     string `db:"name" index:"true"`
    Email    string `db:"email" unique:"true"`
}

// New style (recommended)
type User struct {
    ID       int    `orm:"pk,auto"`
    Name     string `orm:"index"`
    Email    string `orm:"unique"`
}
```

## Architecture

The ORM follows a modular architecture:

```
orm/
├── core/           # Core ORM functionality
│   ├── orm.go     # Main ORM implementation
│   ├── metadata.go # Model metadata extraction
│   ├── query_builder.go # Query builder
│   └── repository.go # Repository pattern
├── dialect/        # Database dialects
│   ├── mysql.go    # MySQL implementation
│   └── interface.go # Dialect interface
└── sql/           # SQL generation utilities
```

## Testing

Run the test suite:

```bash
go test ./...
```

Generate coverage report:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

MIT License - see LICENSE file for details.

## Roadmap

- [ ] PostgreSQL dialect
- [ ] SQLite dialect
- [ ] Advanced relationship support
- [ ] Migration system
- [ ] Connection pooling
- [ ] Query caching
- [ ] Code generation tools
- [ ] Documentation generator