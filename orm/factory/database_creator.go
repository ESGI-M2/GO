package factory

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/ESGI-M2/GO/orm/core/interfaces"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// DatabaseCreator handles automatic database creation
type DatabaseCreator struct{}

// NewDatabaseCreator creates a new database creator instance
func NewDatabaseCreator() *DatabaseCreator {
	return &DatabaseCreator{}
}

// CreateDatabaseIfNotExists creates a database if it doesn't exist
func (dc *DatabaseCreator) CreateDatabaseIfNotExists(config interfaces.ConnectionConfig, dialectType DialectType) error {
	switch strings.ToLower(string(dialectType)) {
	case "mysql":
		return dc.createMySQLDatabase(config)
	case "postgresql", "postgres":
		return dc.createPostgresDatabase(config)
	case "mock":
		// No database creation needed for mock
		return nil
	default:
		return fmt.Errorf("database creation not supported for dialect: %s", dialectType)
	}
}

// createMySQLDatabase creates a MySQL database if it doesn't exist
func (dc *DatabaseCreator) createMySQLDatabase(config interfaces.ConnectionConfig) error {
	// Connect to MySQL server without specifying a database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?parseTime=true&loc=Local",
		config.Username, config.Password, config.Host, config.Port)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL server: %w", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping MySQL server: %w", err)
	}

	// Check if database exists
	var dbExists bool
	checkQuery := "SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?"
	err = db.QueryRow(checkQuery, config.Database).Scan(&dbExists)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check database existence: %w", err)
	}

	// Create database if it doesn't exist
	if err == sql.ErrNoRows {
		createQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", config.Database)
		_, err = db.Exec(createQuery)
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		log.Printf("Created MySQL database: %s", config.Database)
	} else {
		log.Printf("MySQL database already exists: %s", config.Database)
	}

	return nil
}

// createPostgresDatabase creates a PostgreSQL database if it doesn't exist
func (dc *DatabaseCreator) createPostgresDatabase(config interfaces.ConnectionConfig) error {
	// Connect to PostgreSQL server using the 'postgres' database
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable",
		config.Host, config.Port, config.Username, config.Password)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL server: %w", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping PostgreSQL server: %w", err)
	}

	// Check if database exists
	var dbExists bool
	checkQuery := "SELECT 1 FROM pg_database WHERE datname = $1"
	err = db.QueryRow(checkQuery, config.Database).Scan(&dbExists)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check database existence: %w", err)
	}

	// Create database if it doesn't exist
	if err == sql.ErrNoRows {
		createQuery := fmt.Sprintf("CREATE DATABASE \"%s\"", config.Database)
		_, err = db.Exec(createQuery)
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		log.Printf("Created PostgreSQL database: %s", config.Database)
	} else {
		log.Printf("PostgreSQL database already exists: %s", config.Database)
	}

	return nil
}

// EnsureDatabaseExists is a convenience function to ensure database exists
func (dc *DatabaseCreator) EnsureDatabaseExists(config interfaces.ConnectionConfig, dialectType DialectType) error {
	return dc.CreateDatabaseIfNotExists(config, dialectType)
}

// Global database creator instance
var DefaultDatabaseCreator = NewDatabaseCreator()

// CreateDatabaseIfNotExists creates a database using the default creator
func CreateDatabaseIfNotExists(config interfaces.ConnectionConfig, dialectType DialectType) error {
	return DefaultDatabaseCreator.CreateDatabaseIfNotExists(config, dialectType)
}

// EnsureDatabaseExists ensures a database exists using the default creator
func EnsureDatabaseExists(config interfaces.ConnectionConfig, dialectType DialectType) error {
	return DefaultDatabaseCreator.EnsureDatabaseExists(config, dialectType)
}
