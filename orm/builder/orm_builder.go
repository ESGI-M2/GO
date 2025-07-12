package builder

import (
	"fmt"
	"log"

	"github.com/ESGI-M2/GO/orm/core/connection"
	"github.com/ESGI-M2/GO/orm/core/interfaces"
	"github.com/ESGI-M2/GO/orm/core/query"
	"github.com/ESGI-M2/GO/orm/factory"
)

// SimpleORM provides a simplified, user-friendly interface for ORM operations
type SimpleORM struct {
	orm         interfaces.ORM
	dialect     interfaces.Dialect
	config      interfaces.ConnectionConfig
	dialectType factory.DialectType
	autoCreate  bool
	models      []interface{}
	connected   bool
}

// NewSimpleORM creates a new simple ORM instance
func NewSimpleORM() *SimpleORM {
	return &SimpleORM{
		autoCreate: false,
		models:     make([]interface{}, 0),
		connected:  false,
	}
}

// WithDialect sets the database dialect by string or constant
func (s *SimpleORM) WithDialect(dialectType interface{}) *SimpleORM {
	switch dt := dialectType.(type) {
	case string:
		s.dialectType = factory.DialectType(dt)
	case factory.DialectType:
		s.dialectType = dt
	default:
		log.Printf("Warning: Unknown dialect type %T, using MySQL as default", dialectType)
		s.dialectType = factory.MySQL
	}
	return s
}

// WithMySQL sets the dialect to MySQL
func (s *SimpleORM) WithMySQL() *SimpleORM {
	s.dialectType = factory.MySQL
	return s
}

// WithPostgreSQL sets the dialect to PostgreSQL
func (s *SimpleORM) WithPostgreSQL() *SimpleORM {
	s.dialectType = factory.PostgreSQL
	return s
}

// WithConfig sets the connection configuration
func (s *SimpleORM) WithConfig(config interfaces.ConnectionConfig) *SimpleORM {
	s.config = config
	return s
}

// WithConfigBuilder uses a config builder
func (s *SimpleORM) WithConfigBuilder(builder *ConfigBuilder) *SimpleORM {
	config, dialectType, autoCreate, err := builder.Build()
	if err != nil {
		log.Printf("Config builder error: %v", err)
		return s
	}
	s.config = config
	s.dialectType = dialectType
	s.autoCreate = autoCreate
	return s
}

// WithQuickConfig provides a quick way to set common configuration
func (s *SimpleORM) WithQuickConfig(host, database, username, password string) *SimpleORM {
	port := 3306
	if s.dialectType == factory.PostgreSQL || s.dialectType == factory.Postgres {
		port = 5432
	}

	s.config = interfaces.ConnectionConfig{
		Host:     host,
		Port:     port,
		Database: database,
		Username: username,
		Password: password,
	}
	return s
}

// WithEnvConfig loads configuration from environment variables
func (s *SimpleORM) WithEnvConfig() *SimpleORM {
	builder := NewConfigBuilder().WithDialect(s.dialectType).FromEnv()
	config, _, autoCreate, err := builder.Build()
	if err != nil {
		log.Printf("Environment config error: %v", err)
		return s
	}
	s.config = config
	s.autoCreate = autoCreate
	return s
}

// WithAutoCreateDatabase enables automatic database creation
func (s *SimpleORM) WithAutoCreateDatabase() *SimpleORM {
	s.autoCreate = true
	return s
}

// RegisterModel registers a model with the ORM
func (s *SimpleORM) RegisterModel(model interface{}) *SimpleORM {
	s.models = append(s.models, model)
	return s
}

// RegisterModels registers multiple models with the ORM
func (s *SimpleORM) RegisterModels(models ...interface{}) *SimpleORM {
	s.models = append(s.models, models...)
	return s
}

// Connect establishes the database connection
func (s *SimpleORM) Connect() error {
	if s.connected {
		return nil
	}

	// Create dialect
	dialectFactory := factory.NewDialectFactory()
	var err error
	s.dialect, err = dialectFactory.Create(s.dialectType)
	if err != nil {
		return fmt.Errorf("failed to create dialect: %w", err)
	}

	// Create database if needed
	if s.autoCreate && s.dialectType != factory.Mock {
		dbCreator := factory.NewDatabaseCreator()
		if err := dbCreator.CreateDatabaseIfNotExists(s.config, s.dialectType); err != nil {
			log.Printf("Warning: Failed to create database: %v", err)
		}
	}

	// Create ORM instance
	s.orm = connection.NewORM(s.dialect)

	// Connect to database
	if err := s.orm.Connect(s.config); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Register models
	for _, model := range s.models {
		if err := s.orm.RegisterModel(model); err != nil {
			return fmt.Errorf("failed to register model %T: %w", model, err)
		}
	}

	// Migrate tables
	if err := s.orm.Migrate(); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}

	s.connected = true
	log.Printf("Successfully connected to %s database: %s", s.dialectType, s.config.Database)
	return nil
}

// GetORM returns the underlying ORM instance
func (s *SimpleORM) GetORM() interfaces.ORM {
	return s.orm
}

// Query creates a query builder for the model
func (s *SimpleORM) Query(model interface{}) interfaces.QueryBuilder {
	if !s.connected {
		// Return a QueryBuilder with an error instead of panicking
		return &query.BuilderImpl{
			Err: fmt.Errorf("SimpleORM not connected. Call Connect() first."),
		}
	}
	return s.orm.Query(model)
}

// ErrorRepository is a simple repository that returns errors on all methods
type ErrorRepository struct {
	err error
}

// NewErrorRepository creates a new error repository
func NewErrorRepository(err error) *ErrorRepository {
	return &ErrorRepository{err: err}
}

// All methods return the error
func (r *ErrorRepository) Find(id interface{}) (interface{}, error) { return nil, r.err }
func (r *ErrorRepository) FindAll() ([]interface{}, error)          { return nil, r.err }
func (r *ErrorRepository) FindBy(criteria map[string]interface{}) ([]interface{}, error) {
	return nil, r.err
}
func (r *ErrorRepository) FindOneBy(criteria map[string]interface{}) (interface{}, error) {
	return nil, r.err
}
func (r *ErrorRepository) Save(entity interface{}) error                  { return r.err }
func (r *ErrorRepository) Update(entity interface{}) error                { return r.err }
func (r *ErrorRepository) Delete(entity interface{}) error                { return r.err }
func (r *ErrorRepository) DeleteBy(criteria map[string]interface{}) error { return r.err }
func (r *ErrorRepository) Count() (int64, error)                          { return 0, r.err }
func (r *ErrorRepository) Exists(id interface{}) (bool, error)            { return false, r.err }
func (r *ErrorRepository) FindWithRelations(id interface{}, relations ...string) (interface{}, error) {
	return nil, r.err
}
func (r *ErrorRepository) FindAllWithRelations(relations ...string) ([]interface{}, error) {
	return nil, r.err
}
func (r *ErrorRepository) FindByWithRelations(criteria map[string]interface{}, relations ...string) ([]interface{}, error) {
	return nil, r.err
}
func (r *ErrorRepository) BatchCreate(entities []interface{}) error                     { return r.err }
func (r *ErrorRepository) BatchUpdate(entities []interface{}) error                     { return r.err }
func (r *ErrorRepository) BatchDelete(entities []interface{}) error                     { return r.err }
func (r *ErrorRepository) SoftDelete(entity interface{}) error                          { return r.err }
func (r *ErrorRepository) Restore(entity interface{}) error                             { return r.err }
func (r *ErrorRepository) ForceDelete(entity interface{}) error                         { return r.err }
func (r *ErrorRepository) FindTrashed() ([]interface{}, error)                          { return nil, r.err }
func (r *ErrorRepository) RestoreBy(criteria map[string]interface{}) error              { return r.err }
func (r *ErrorRepository) Scope(name string, args ...interface{}) interfaces.Repository { return r }
func (r *ErrorRepository) Chunk(size int, fn func([]interface{}) error) error           { return r.err }
func (r *ErrorRepository) Each(fn func(interface{}) error) error                        { return r.err }
func (r *ErrorRepository) Pluck(field string) ([]interface{}, error)                    { return nil, r.err }
func (r *ErrorRepository) Value(field string) (interface{}, error)                      { return nil, r.err }
func (r *ErrorRepository) Increment(field string, amount interface{}) error             { return r.err }
func (r *ErrorRepository) Decrement(field string, amount interface{}) error             { return r.err }

// Repository creates a repository for the model
func (s *SimpleORM) Repository(model interface{}) interfaces.Repository {
	if !s.connected {
		// Return an error repository instead of panicking
		return NewErrorRepository(fmt.Errorf("SimpleORM not connected. Call Connect() first."))
	}
	return s.orm.Repository(model)
}

// Raw creates a raw SQL query builder
func (s *SimpleORM) Raw(sql string, args ...interface{}) interfaces.QueryBuilder {
	if !s.connected {
		// Return a QueryBuilder with an error instead of panicking
		return &query.BuilderImpl{
			Err: fmt.Errorf("SimpleORM not connected. Call Connect() first."),
		}
	}
	return s.orm.Raw(sql, args...)
}

// Transaction executes a function within a transaction
func (s *SimpleORM) Transaction(fn func(interfaces.ORM) error) error {
	if !s.connected {
		return fmt.Errorf("SimpleORM not connected. Call Connect() first.")
	}
	return s.orm.Transaction(fn)
}

// Close closes the database connection
func (s *SimpleORM) Close() error {
	if s.orm != nil {
		s.connected = false
		return s.orm.Close()
	}
	return nil
}

// IsConnected returns whether the ORM is connected
func (s *SimpleORM) IsConnected() bool {
	return s.connected
}

// GetConfig returns the current configuration
func (s *SimpleORM) GetConfig() interfaces.ConnectionConfig {
	return s.config
}

// GetDialectType returns the current dialect type
func (s *SimpleORM) GetDialectType() factory.DialectType {
	return s.dialectType
}

// Convenience functions for common patterns

// NewMySQL creates a new MySQL SimpleORM instance
func NewMySQL() *SimpleORM {
	return NewSimpleORM().WithMySQL()
}

// NewPostgreSQL creates a new PostgreSQL SimpleORM instance
func NewPostgreSQL() *SimpleORM {
	return NewSimpleORM().WithPostgreSQL()
}

// NewMySQLFromEnv creates a MySQL SimpleORM instance from environment variables
func NewMySQLFromEnv() *SimpleORM {
	return NewMySQL().WithEnvConfig()
}

// NewPostgreSQLFromEnv creates a PostgreSQL SimpleORM instance from environment variables
func NewPostgreSQLFromEnv() *SimpleORM {
	return NewPostgreSQL().WithEnvConfig()
}

// QuickSetup provides a one-liner setup for common use cases
func QuickSetup(dialectType string, host, database, username, password string, models ...interface{}) (*SimpleORM, error) {
	orm := NewSimpleORM().
		WithDialect(dialectType).
		WithQuickConfig(host, database, username, password).
		WithAutoCreateDatabase().
		RegisterModels(models...)

	err := orm.Connect()
	return orm, err
}

// QuickSetupFromEnv provides a one-liner setup using environment variables
func QuickSetupFromEnv(dialectType string, models ...interface{}) (*SimpleORM, error) {
	orm := NewSimpleORM().
		WithDialect(dialectType).
		WithEnvConfig().
		WithAutoCreateDatabase().
		RegisterModels(models...)

	err := orm.Connect()
	return orm, err
}

// Example usage patterns (for documentation)

// ExampleUsage demonstrates typical usage patterns
func ExampleUsage() {
	// Method 1: Fluent API with manual configuration
	orm := NewSimpleORM().
		WithMySQL().
		WithQuickConfig("localhost", "myapp", "user", "password").
		WithAutoCreateDatabase().
		RegisterModel(&User{}).
		RegisterModel(&Post{})

	if err := orm.Connect(); err != nil {
		log.Fatal(err)
	}
	defer orm.Close()

	// Method 2: Using environment variables
	orm2 := NewMySQLFromEnv().
		WithAutoCreateDatabase().
		RegisterModels(&User{}, &Post{})

	if err := orm2.Connect(); err != nil {
		log.Fatal(err)
	}
	defer orm2.Close()

	// Method 3: Quick setup one-liner
	orm3, err := QuickSetup("mysql", "localhost", "myapp", "user", "password", &User{}, &Post{})
	if err != nil {
		log.Fatal(err)
	}
	defer orm3.Close()

	// Method 4: Config builder
	config := MySQL().
		WithHost("localhost").
		WithDatabase("myapp").
		WithCredentials("user", "password").
		WithAutoCreateDatabase()

	orm4 := NewSimpleORM().WithConfigBuilder(config).RegisterModels(&User{}, &Post{})
	if err := orm4.Connect(); err != nil {
		log.Fatal(err)
	}
	defer orm4.Close()
}

// Example models for documentation
type User struct {
	ID    int    `db:"id" primary:"true" autoincrement:"true"`
	Name  string `db:"name"`
	Email string `db:"email" unique:"true"`
}

type Post struct {
	ID      int    `db:"id" primary:"true" autoincrement:"true"`
	Title   string `db:"title"`
	Content string `db:"content"`
	UserID  int    `db:"user_id" foreign:"users.id"`
}
