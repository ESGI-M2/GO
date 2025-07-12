# Simple ORM Usage Guide

This guide shows you how to use the enhanced ORM library with convenient features for easy dialect selection and automatic database creation.

## 🚀 Quick Start

### 1. One-liner Setup

```go
// Simplest possible setup
orm, err := builder.QuickSetup("mysql", "localhost", "myapp", "user", "password", &User{}, &Post{})
if err != nil {
    log.Fatal(err)
}
defer orm.Close()
```

### 2. Environment-based Setup

```go
// Set environment variables: MYSQL_HOST, MYSQL_DATABASE, MYSQL_USER, MYSQL_PASSWORD
orm, err := builder.QuickSetupFromEnv("mysql", &User{}, &Post{})
if err != nil {
    log.Fatal(err)
}
defer orm.Close()
```

### 3. Fluent API Setup

```go
orm := builder.NewSimpleORM().
    WithMySQL().
    WithQuickConfig("localhost", "myapp", "user", "password").
    WithAutoCreateDatabase().
    RegisterModels(&User{}, &Post{})

if err := orm.Connect(); err != nil {
    log.Fatal(err)
}
defer orm.Close()
```

## 📋 Features

### ✅ Easy Dialect Selection

Choose database dialects using strings or constants:

```go
// Method 1: Using strings
orm := builder.NewSimpleORM().WithDialect("mysql")
orm := builder.NewSimpleORM().WithDialect("postgresql")

// Method 2: Using constants
orm := builder.NewSimpleORM().WithDialect(factory.MySQL)
orm := builder.NewSimpleORM().WithDialect(factory.PostgreSQL)

// Method 3: Convenience functions
orm := builder.NewMySQL()
orm := builder.NewPostgreSQL()

// Method 4: Factory pattern
dialect, err := factory.CreateDialect(factory.MySQL)
dialect, err := factory.CreateDialectFromString("postgresql")
```

### ✅ Automatic Database Creation

The ORM can automatically create databases if they don't exist:

```go
orm := builder.NewSimpleORM().
    WithMySQL().
    WithQuickConfig("localhost", "new_database", "user", "password").
    WithAutoCreateDatabase()  // Enable auto-creation

if err := orm.Connect(); err != nil {
    log.Fatal(err)
}
```

### ✅ Configuration Builder

Build configurations fluently:

```go
// MySQL configuration
config := builder.MySQL().
    WithHost("localhost").
    WithDatabase("myapp").
    WithCredentials("user", "password").
    WithConnectionPool(25, 5).
    WithAutoCreateDatabase()

orm := builder.NewSimpleORM().
    WithConfigBuilder(config).
    RegisterModels(&User{}, &Post{})
```

### ✅ Environment Variables

Load configuration from environment variables:

```go
// MySQL environment variables
// MYSQL_HOST, MYSQL_PORT, MYSQL_DATABASE, MYSQL_USER, MYSQL_PASSWORD
orm := builder.NewMySQLFromEnv().
    WithAutoCreateDatabase().
    RegisterModels(&User{}, &Post{})

// PostgreSQL environment variables
// POSTGRES_HOST, POSTGRES_PORT, POSTGRES_DB, POSTGRES_USER, POSTGRES_PASSWORD
orm := builder.NewPostgreSQLFromEnv().
    WithAutoCreateDatabase().
    RegisterModels(&User{}, &Post{})
```

## 🏗️ Usage Patterns

### Pattern 1: Development Setup

```go
// Quick development setup with auto-creation
orm := builder.NewSimpleORM().
    WithMySQL().
    WithQuickConfig("localhost", "dev_db", "root", "password").
    WithAutoCreateDatabase().
    RegisterModels(&User{}, &Post{}, &Category{})

if err := orm.Connect(); err != nil {
    log.Fatal(err)
}
defer orm.Close()
```

### Pattern 2: Production Setup

```go
// Production setup with environment variables
orm := builder.NewSimpleORM().
    WithDialect(os.Getenv("DB_DIALECT")).
    WithEnvConfig().
    RegisterModels(&User{}, &Post{})

if err := orm.Connect(); err != nil {
    log.Fatal(err)
}
defer orm.Close()
```

### Pattern 3: Testing Setup

```go
// Mock setup for testing
orm := builder.NewSimpleORM().
    WithDialect(factory.Mock).
    RegisterModels(&User{}, &Post{})

if err := orm.Connect(); err != nil {
    log.Fatal(err)
}
defer orm.Close()
```

## 📝 Model Definition

Define your models with ORM tags:

```go
type User struct {
    ID        int       `db:"id" primary:"true" autoincrement:"true"`
    Name      string    `db:"name"`
    Email     string    `db:"email" unique:"true"`
    Age       int       `db:"age"`
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
```

## 🔧 Advanced Usage

### Custom Configuration

```go
orm := builder.NewSimpleORM().
    WithDialect("mysql").
    WithConfig(interfaces.ConnectionConfig{
        Host:            "localhost",
        Port:            3306,
        Database:        "myapp",
        Username:        "user",
        Password:        "password",
        MaxOpenConns:    25,
        MaxIdleConns:    5,
        ConnMaxLifetime: 300,
    }).
    RegisterModels(&User{}, &Post{})
```

### Using the ORM

```go
// Get the underlying ORM instance
underlyingORM := orm.GetORM()

// Use repositories
userRepo := orm.Repository(&User{})
users, err := userRepo.FindAll()

// Use query builder
query := orm.Query(&User{}).Where("age", ">", 18).OrderBy("name", "ASC")
results, err := query.Find()

// Use raw SQL
rawQuery := orm.Raw("SELECT COUNT(*) FROM users WHERE age > ?", 18)
count, err := rawQuery.Count()

// Use transactions
err = orm.Transaction(func(tx interfaces.ORM) error {
    // Your transaction logic here
    return nil
})
```

## 🗄️ Database Creation

### Manual Database Creation

```go
config := interfaces.ConnectionConfig{
    Host:     "localhost",
    Port:     3306,
    Database: "new_database",
    Username: "user",
    Password: "password",
}

// Create database if it doesn't exist
err := factory.CreateDatabaseIfNotExists(config, factory.MySQL)
if err != nil {
    log.Fatal(err)
}
```

### Automatic Database Creation

```go
// Enable auto-creation during ORM setup
orm := builder.NewSimpleORM().
    WithMySQL().
    WithQuickConfig("localhost", "auto_db", "user", "password").
    WithAutoCreateDatabase()  // This will create the database automatically

if err := orm.Connect(); err != nil {
    log.Fatal(err)
}
```

## 🔍 Dialect Information

```go
// Check available dialects
factory := factory.NewDialectFactory()
dialects := factory.GetAvailableDialects()
fmt.Printf("Available dialects: %v\n", dialects)

// Check if a dialect is supported
supported := factory.IsSupported(factory.MySQL)
fmt.Printf("MySQL supported: %v\n", supported)
```

## 📖 Environment Variables

### MySQL Environment Variables
```bash
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DATABASE=myapp
MYSQL_USER=user
MYSQL_PASSWORD=password
```

### PostgreSQL Environment Variables
```bash
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=myapp
POSTGRES_USER=user
POSTGRES_PASSWORD=password
```

## 🎯 Examples

### Complete Example

```go
package main

import (
    "log"
    "time"
    
    "github.com/ESGI-M2/GO/orm/builder"
)

type User struct {
    ID        int       `db:"id" primary:"true" autoincrement:"true"`
    Name      string    `db:"name"`
    Email     string    `db:"email" unique:"true"`
    CreatedAt time.Time `db:"created_at"`
}

func main() {
    // Quick setup with all features
    orm, err := builder.QuickSetup("mysql", "localhost", "myapp", "user", "password", &User{})
    if err != nil {
        log.Fatal(err)
    }
    defer orm.Close()
    
    // Use the ORM
    userRepo := orm.Repository(&User{})
    
    // Create a user
    user := &User{
        Name:      "John Doe",
        Email:     "john@example.com",
        CreatedAt: time.Now(),
    }
    
    if err := userRepo.Save(user); err != nil {
        log.Printf("Failed to save user: %v", err)
    }
    
    // Query users
    users, err := userRepo.FindAll()
    if err != nil {
        log.Printf("Failed to find users: %v", err)
    }
    
    log.Printf("Found %d users", len(users))
}
```

## 🔄 Migration from Old API

### Before (Old API)
```go
// Old way - complex setup
mysqlDialect := dialect.NewMySQLDialect()
orm := connection.NewORM(mysqlDialect)

config := interfaces.ConnectionConfig{
    Host:     "localhost",
    Port:     3306,
    Database: "myapp",
    Username: "user",
    Password: "password",
}

if err := orm.Connect(config); err != nil {
    log.Fatal(err)
}

if err := orm.RegisterModel(&User{}); err != nil {
    log.Fatal(err)
}

if err := orm.Migrate(); err != nil {
    log.Fatal(err)
}
```

### After (New Simple API)
```go
// New way - simple setup
orm, err := builder.QuickSetup("mysql", "localhost", "myapp", "user", "password", &User{})
if err != nil {
    log.Fatal(err)
}
defer orm.Close()
```

## 📚 Additional Resources

- [Complete API Documentation](API_DOCUMENTATION.md)
- [Examples Directory](examples/)
- [Migration Guide](MIGRATION.md)

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch
3. Add tests for new features
4. Run `go test ./...`
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details. 