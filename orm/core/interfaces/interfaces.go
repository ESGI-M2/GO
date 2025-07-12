package interfaces

import (
	"context"
	"database/sql"
	"reflect"
	"time"
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
	// New advanced features
	WithCache(ttl int) ORM
	WithConnectionPool(maxOpen, maxIdle int) ORM
	EnableQueryLog() ORM
	DisableQueryLog() ORM
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
	// New advanced features
	FullTextSearch(field, query string) string
	GetRandomFunction() string
	GetDateFunction() string
	GetJSONExtract() string
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
	// New advanced features
	WhereOr(conditions ...WhereCondition) QueryBuilder
	WhereRaw(condition string, args ...interface{}) QueryBuilder
	WhereBetween(field string, min, max interface{}) QueryBuilder
	WhereNotBetween(field string, min, max interface{}) QueryBuilder
	WhereNull(field string) QueryBuilder
	WhereNotNull(field string) QueryBuilder
	WhereLike(field, pattern string) QueryBuilder
	WhereNotLike(field, pattern string) QueryBuilder
	WhereRegexp(field, pattern string) QueryBuilder
	WhereNotRegexp(field, pattern string) QueryBuilder
	FullTextSearch(fields []string, query string) QueryBuilder
	SubQuery(alias string, fn func(QueryBuilder) QueryBuilder) QueryBuilder
	With(relation string, fn func(QueryBuilder) QueryBuilder) QueryBuilder
	WithCount(relation string) QueryBuilder
	WithExists(relation string, fn func(QueryBuilder) QueryBuilder) QueryBuilder
	CursorPaginate(cursorField string, cursorValue interface{}, limit int) QueryBuilder
	OffsetPaginate(page, perPage int) QueryBuilder
	ForUpdate() QueryBuilder
	ForShare() QueryBuilder
	Distinct() QueryBuilder
	Union(other QueryBuilder) QueryBuilder
	UnionAll(other QueryBuilder) QueryBuilder
	Lock(lockType string) QueryBuilder
	Cache(ttl int) QueryBuilder
	WithoutCache() QueryBuilder
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
	// New advanced features
	FindWithRelations(id interface{}, relations ...string) (interface{}, error)
	FindAllWithRelations(relations ...string) ([]interface{}, error)
	FindByWithRelations(criteria map[string]interface{}, relations ...string) ([]interface{}, error)
	BatchCreate(entities []interface{}) error
	BatchUpdate(entities []interface{}) error
	BatchDelete(entities []interface{}) error
	SoftDelete(entity interface{}) error
	Restore(entity interface{}) error
	ForceDelete(entity interface{}) error
	FindTrashed() ([]interface{}, error)
	RestoreBy(criteria map[string]interface{}) error
	Scope(name string, args ...interface{}) Repository
	Chunk(size int, fn func([]interface{}) error) error
	Each(fn func(interface{}) error) error
	Pluck(field string) ([]interface{}, error)
	Value(field string) (interface{}, error)
	Increment(field string, amount interface{}) error
	Decrement(field string, amount interface{}) error
}

// ConnectionConfig defines database connection configuration
type ConnectionConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	SSLMode  string
	// New advanced features
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int // in seconds
	ConnMaxIdleTime int // in seconds
	QueryTimeout    int // in seconds
	EnableQueryLog  bool
	CacheTTL        int // in seconds
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
	// New advanced features
	SoftDelete bool
	Timestamp  bool
	JSON       bool
	FullText   bool
	Encrypted  bool
	Validation []ValidationRule
}

// ValidationRule represents a validation rule
type ValidationRule struct {
	Type    string
	Value   interface{}
	Message string
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
	// New advanced features
	Type    string // BTREE, HASH, GIN, etc.
	Partial string // partial index condition
}

// Relation represents a relationship between models
type Relation struct {
	Type          RelationType
	TargetModel   reflect.Type
	ForeignKey    string
	ReferencedKey string
	JoinTable     string
	Lazy          bool
	// New advanced features
	Eager          bool
	Cascade        bool
	Polymorphic    bool
	MorphType      string
	MorphID        string
	Through        string
	WithPivot      []string
	As             string
	WithTimestamps bool
}

// RelationType represents the type of relationship
type RelationType int

const (
	OneToOne RelationType = iota
	OneToMany
	ManyToOne
	ManyToMany
	HasOne
	HasMany
	BelongsTo
	BelongsToMany
	MorphOne
	MorphMany
	MorphTo
	MorphToMany
	MorphedByMany
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
	// New advanced features
	SoftDeletes bool
	Timestamps  bool
	CreatedAt   string
	UpdatedAt   string
	DeletedAt   string
	Hooks       *ModelHooks
	Scopes      map[string]func(QueryBuilder) QueryBuilder
	Validation  []ValidationRule
	Hidden      []string
	Visible     []string
	Fillable    []string
	Guarded     []string
	Appends     []string
	Casts       map[string]string
	Events      map[string][]func(interface{}) error
}

// ModelHooks represents model lifecycle hooks
type ModelHooks struct {
	BeforeCreate []func(interface{}) error
	AfterCreate  []func(interface{}) error
	BeforeUpdate []func(interface{}) error
	AfterUpdate  []func(interface{}) error
	BeforeDelete []func(interface{}) error
	AfterDelete  []func(interface{}) error
	BeforeSave   []func(interface{}) error
	AfterSave    []func(interface{}) error
}

// WhereCondition represents a WHERE clause condition
type WhereCondition struct {
	Field    string
	Operator string
	Value    interface{}
	Logical  string // AND, OR
	// New advanced features
	Raw      bool
	SubQuery QueryBuilder
	Nested   []WhereCondition
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
	// New advanced features
	Alias    string
	SubQuery QueryBuilder
}

// PaginationResult represents pagination result
type PaginationResult struct {
	Data        []interface{}
	Total       int64
	PerPage     int
	CurrentPage int
	LastPage    int
	From        int
	To          int
	HasMore     bool
	NextCursor  interface{}
	PrevCursor  interface{}
}

// Cache interface for query caching
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl int) error
	Delete(key string) error
	Clear() error
	Has(key string) bool
}

// QueryLog represents a query log entry
type QueryLog struct {
	SQL      string
	Args     []interface{}
	Duration time.Duration
	Time     time.Time
	Error    error
}

// QueryLogger interface for query logging
type QueryLogger interface {
	Log(log QueryLog)
	GetLogs() []QueryLog
	ClearLogs()
}
