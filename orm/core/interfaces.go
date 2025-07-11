package core

import (
	"context"
	"database/sql"
	"reflect"
)

// Dialect interface defines the contract for database-specific implementations
type Dialect interface {
	// Connection management
	Connect(config ConnectionConfig) error
	Close() error
	Ping() error
	
	// Query execution
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	
	// Transaction support
	Begin() (Transaction, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Transaction, error)
	
	// Schema operations
	CreateTable(tableName string, columns []Column) error
	DropTable(tableName string) error
	TableExists(tableName string) (bool, error)
	
	// Type mapping
	GetSQLType(goType reflect.Type) string
	GetPlaceholder(index int) string
}

// Transaction interface for database transactions
type Transaction interface {
	Commit() error
	Rollback() error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// ConnectionConfig holds database connection parameters
type ConnectionConfig struct {
	Driver   string
	Host     string
	Port     int
	Database string
	Username string
	Password string
	SSLMode  string
	MaxOpenConns int
	MaxIdleConns int
	ConnMaxLifetime int // in seconds
}

// Column represents a database column definition
type Column struct {
	Name         string
	Type         string
	Length       int
	Nullable     bool
	PrimaryKey   bool
	AutoIncrement bool
	Default      interface{}
	Unique       bool
	Index        bool
	ForeignKey   *ForeignKey
}

// ForeignKey represents a foreign key constraint
type ForeignKey struct {
	ReferencedTable  string
	ReferencedColumn string
	OnDelete         string // CASCADE, SET NULL, RESTRICT
	OnUpdate         string // CASCADE, SET NULL, RESTRICT
}

// ModelMetadata holds reflection information about a model
type ModelMetadata struct {
	Type           reflect.Type
	TableName      string
	Columns        []Column
	PrimaryKey     string
	AutoIncrement  string
	Relations      map[string]*Relation
	Indexes        []Index
}

// Relation represents a relationship between models
type Relation struct {
	Type          RelationType
	TargetModel   reflect.Type
	ForeignKey    string
	ReferencedKey string
	JoinTable     string // for many-to-many
	Lazy          bool
}

// RelationType defines the type of relationship
type RelationType int

const (
	OneToOne RelationType = iota
	OneToMany
	ManyToOne
	ManyToMany
)

// Index represents a database index
type Index struct {
	Name    string
	Columns []string
	Unique  bool
}

// QueryBuilder interface for building SQL queries
type QueryBuilder interface {
	Select(fields ...string) QueryBuilder
	From(table string) QueryBuilder
	Where(field, operator string, value interface{}) QueryBuilder
	WhereIn(field string, values []interface{}) QueryBuilder
	WhereNotIn(field string, values []interface{}) QueryBuilder
	OrderBy(field, direction string) QueryBuilder
	GroupBy(fields ...string) QueryBuilder
	Having(condition string, args ...interface{}) QueryBuilder
	Limit(limit int) QueryBuilder
	Offset(offset int) QueryBuilder
	Join(table, condition string) QueryBuilder
	LeftJoin(table, condition string) QueryBuilder
	RightJoin(table, condition string) QueryBuilder
	InnerJoin(table, condition string) QueryBuilder
	
	// Execution methods
	Find() ([]map[string]interface{}, error)
	FindOne() (map[string]interface{}, error)
	Count() (int64, error)
	Exists() (bool, error)
	
	// Raw SQL
	Raw(sql string, args ...interface{}) QueryBuilder
	GetSQL() string
	GetArgs() []interface{}
}

// Repository interface for model-specific operations
type Repository interface {
	Find(id interface{}) (interface{}, error)
	FindAll() ([]interface{}, error)
	FindBy(criteria map[string]interface{}) ([]interface{}, error)
	FindOneBy(criteria map[string]interface{}) (interface{}, error)
	Save(entity interface{}) error
	Update(entity interface{}) error
	Delete(entity interface{}) error
	DeleteBy(criteria map[string]interface{}) error
	Count() (int64, error)
	Exists(id interface{}) (bool, error)
}

// ORM is the main interface for the ORM system
type ORM interface {
	// Connection management
	Connect(config ConnectionConfig) error
	Close() error
	
	// Model registration
	RegisterModel(model interface{}) error
	GetMetadata(model interface{}) (*ModelMetadata, error)
	
	// Query building
	Query(model interface{}) QueryBuilder
	Raw(sql string, args ...interface{}) QueryBuilder
	
	// Repository pattern
	Repository(model interface{}) Repository
	
	// Transaction support
	Transaction(fn func(ORM) error) error
	TransactionWithContext(ctx context.Context, fn func(ORM) error) error
	
	// Schema operations
	CreateTable(model interface{}) error
	DropTable(model interface{}) error
	Migrate() error
	
	// Utility methods
	GetDialect() Dialect
	IsConnected() bool
} 