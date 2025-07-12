package orm

import (
	"github.com/ESGI-M2/GO/dialect"
	"github.com/ESGI-M2/GO/orm/core/connection"
	"github.com/ESGI-M2/GO/orm/core/interfaces"
)

// ORM provides the main interface for the ORM
type ORM = interfaces.ORM

// New creates a new ORM instance
func New(dialect interfaces.Dialect) ORM {
	return connection.NewORM(dialect)
}

// NewWithMySQL creates a new ORM instance with MySQL dialect
func NewWithMySQL() ORM {
	return New(dialect.NewMySQLDialect())
}

// NewWithPostgres creates a new ORM instance with PostgreSQL dialect
func NewWithPostgres() ORM {
	return New(dialect.NewPostgresDialect())
}

// NewConnectionConfigFromEnv creates a connection config from environment variables
// It automatically loads .env files and reads MySQL environment variables
func NewConnectionConfigFromEnv() ConnectionConfig {
	return dialect.NewConnectionConfigFromEnv()
}

// NewPostgresConnectionConfigFromEnv creates a connection config from environment variables
// It automatically loads .env files and reads PostgreSQL environment variables
func NewPostgresConnectionConfigFromEnv() ConnectionConfig {
	return dialect.NewPostgresConnectionConfigFromEnv()
}

// ConnectionConfig represents database connection configuration
type ConnectionConfig = interfaces.ConnectionConfig

// QueryBuilder represents a query builder
type QueryBuilder = interfaces.QueryBuilder

// Repository represents a repository
type Repository = interfaces.Repository

// Dialect represents a database dialect
type Dialect = interfaces.Dialect

// Transaction represents a database transaction
type Transaction = interfaces.Transaction

// Column represents a database column
type Column = interfaces.Column

// ForeignKey represents a foreign key constraint
type ForeignKey = interfaces.ForeignKey

// Index represents a database index
type Index = interfaces.Index

// Relation represents a relationship between models
type Relation = interfaces.Relation

// RelationType represents the type of relationship
type RelationType = interfaces.RelationType

// ModelMetadata represents metadata for a model
type ModelMetadata = interfaces.ModelMetadata

// WhereCondition represents a WHERE clause condition
type WhereCondition = interfaces.WhereCondition

// OrderBy represents an ORDER BY clause
type OrderBy = interfaces.OrderBy

// Join represents a JOIN clause
type Join = interfaces.Join

// Constants for relation types
const (
	OneToOne   = interfaces.OneToOne
	OneToMany  = interfaces.OneToMany
	ManyToOne  = interfaces.ManyToOne
	ManyToMany = interfaces.ManyToMany
)
