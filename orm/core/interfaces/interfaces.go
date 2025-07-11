package interfaces

import (
	"context"
	"database/sql"
	"reflect"
)

// ORM defines the main ORM interface
type ORM interface {
	Connect(config ConnectionConfig) error
	Close() error
	IsConnected() bool
	GetDialect() Dialect
	RegisterModel(model interface{}) error
	GetMetadata(model interface{}) (*ModelMetadata, error)
	Query(model interface{}) QueryBuilder
	Raw(sql string, args ...interface{}) QueryBuilder
	Repository(model interface{}) Repository
	Transaction(fn func(ORM) error) error
	TransactionWithContext(ctx context.Context, fn func(ORM) error) error
	CreateTable(model interface{}) error
	DropTable(model interface{}) error
	Migrate() error
}

// Dialect defines the database dialect interface
type Dialect interface {
	Connect(config ConnectionConfig) error
	Close() error
	Ping() error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Begin() (Transaction, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Transaction, error)
	CreateTable(tableName string, columns []Column) error
	DropTable(tableName string) error
	TableExists(tableName string) (bool, error)
	GetSQLType(goType reflect.Type) string
	GetPlaceholder(index int) string
}

// Transaction defines the transaction interface
type Transaction interface {
	Commit() error
	Rollback() error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// QueryBuilder defines the query builder interface
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
	Find() ([]map[string]interface{}, error)
	FindOne() (map[string]interface{}, error)
	Count() (int64, error)
	Exists() (bool, error)
	Raw(sql string, args ...interface{}) QueryBuilder
	GetSQL() string
	GetArgs() []interface{}
}

// Repository defines the repository interface
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

// ConnectionConfig defines database connection configuration
type ConnectionConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	SSLMode  string
}

// Column represents a database column
type Column struct {
	Name          string
	Type          string
	Length        int
	PrimaryKey    bool
	AutoIncrement bool
	Unique        bool
	Index         bool
	Nullable      bool
	Default       interface{}
	ForeignKey    *ForeignKey
}

// ForeignKey represents a foreign key constraint
type ForeignKey struct {
	ReferencedTable  string
	ReferencedColumn string
	OnDelete         string
	OnUpdate         string
}

// Index represents a database index
type Index struct {
	Name    string
	Columns []string
	Unique  bool
}

// Relation represents a relationship between models
type Relation struct {
	Type          RelationType
	TargetModel   reflect.Type
	ForeignKey    string
	ReferencedKey string
	JoinTable     string
	Lazy          bool
}

// RelationType represents the type of relationship
type RelationType int

const (
	OneToOne RelationType = iota
	OneToMany
	ManyToOne
	ManyToMany
)

// ModelMetadata represents metadata for a model
type ModelMetadata struct {
	Type          reflect.Type
	TableName     string
	Columns       []Column
	PrimaryKey    string
	AutoIncrement string
	Relations     map[string]*Relation
	Indexes       []Index
}

// WhereCondition represents a WHERE clause condition
type WhereCondition struct {
	Field    string
	Operator string
	Value    interface{}
	Logical  string // AND, OR
}

// OrderBy represents an ORDER BY clause
type OrderBy struct {
	Field     string
	Direction string // ASC, DESC
}

// Join represents a JOIN clause
type Join struct {
	Type      string // INNER, LEFT, RIGHT, FULL
	Table     string
	Condition string
}
