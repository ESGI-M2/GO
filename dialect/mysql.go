package dialect

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"project/orm/core/interfaces"
	"reflect"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// MySQLDialect implements the interfaces.Dialect interface for MySQL
type MySQLDialect struct {
	db *sql.DB
}

// NewMySQLDialect creates a new MySQL dialect instance
func NewMySQLDialect() *MySQLDialect {
	return &MySQLDialect{}
}

// Connect establishes a connection to MySQL database
func (m *MySQLDialect) Connect(config interfaces.ConnectionConfig) error {
	var err error

	// Load environment variables if .env file exists
	godotenv.Load("../.env")

	// Use config values or fallback to environment variables
	user := config.Username
	if user == "" {
		user = os.Getenv("MYSQL_USER")
	}

	pass := config.Password
	if pass == "" {
		pass = os.Getenv("MYSQL_PASSWORD")
	}

	host := config.Host
	if host == "" {
		host = os.Getenv("MYSQL_HOST")
	}

	database := config.Database
	if database == "" {
		database = os.Getenv("MYSQL_DATABASE")
	}

	port := config.Port
	if port == 0 {
		port = 3306
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		user, pass, host, port, database)

	m.db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	// Configure connection pool with default values
	m.db.SetMaxOpenConns(25)
	m.db.SetMaxIdleConns(5)
	m.db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := m.Ping(); err != nil {
		return fmt.Errorf("failed to ping MySQL: %w", err)
	}

	log.Println("Connected to MySQL successfully")
	return nil
}

// Close closes the database connection
func (m *MySQLDialect) Close() error {
	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

// Ping tests the database connection
func (m *MySQLDialect) Ping() error {
	if m.db == nil {
		return fmt.Errorf("database connection not established")
	}
	return m.db.Ping()
}

// Exec executes a query without returning rows
func (m *MySQLDialect) Exec(query string, args ...interface{}) (sql.Result, error) {
	if m.db == nil {
		return nil, fmt.Errorf("database connection not established")
	}
	return m.db.Exec(query, args...)
}

// Query executes a query that returns rows
func (m *MySQLDialect) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if m.db == nil {
		return nil, fmt.Errorf("database connection not established")
	}
	return m.db.Query(query, args...)
}

// QueryRow executes a query that returns a single row
func (m *MySQLDialect) QueryRow(query string, args ...interface{}) *sql.Row {
	if m.db == nil {
		return nil
	}
	return m.db.QueryRow(query, args...)
}

// Begin starts a new transaction
func (m *MySQLDialect) Begin() (interfaces.Transaction, error) {
	if m.db == nil {
		return nil, fmt.Errorf("database connection not established")
	}
	tx, err := m.db.Begin()
	if err != nil {
		return nil, err
	}
	return &MySQLTransaction{tx: tx}, nil
}

// BeginTx starts a new transaction with options
func (m *MySQLDialect) BeginTx(ctx context.Context, opts *sql.TxOptions) (interfaces.Transaction, error) {
	if m.db == nil {
		return nil, fmt.Errorf("database connection not established")
	}
	tx, err := m.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &MySQLTransaction{tx: tx}, nil
}

// CreateTable creates a table with the given columns
func (m *MySQLDialect) CreateTable(tableName string, columns []interfaces.Column) error {
	var columnDefs []string

	for _, col := range columns {
		def := m.buildColumnDefinition(col)
		columnDefs = append(columnDefs, def)
	}

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n  %s\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci",
		tableName, strings.Join(columnDefs, ",\n  "))

	_, err := m.Exec(query)
	return err
}

// DropTable drops a table
func (m *MySQLDialect) DropTable(tableName string) error {
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	_, err := m.Exec(query)
	return err
}

// TableExists checks if a table exists
func (m *MySQLDialect) TableExists(tableName string) (bool, error) {
	query := "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?"
	var count int
	err := m.QueryRow(query, tableName).Scan(&count)
	return count > 0, err
}

// GetSQLType maps Go types to MySQL SQL types
func (m *MySQLDialect) GetSQLType(goType reflect.Type) string {
	switch goType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return "INT"
	case reflect.Int64:
		return "BIGINT"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return "INT UNSIGNED"
	case reflect.Uint64:
		return "BIGINT UNSIGNED"
	case reflect.Float32:
		return "FLOAT"
	case reflect.Float64:
		return "DOUBLE"
	case reflect.String:
		return "VARCHAR(255)"
	case reflect.Bool:
		return "TINYINT(1)"
	case reflect.Struct:
		if goType == reflect.TypeOf(time.Time{}) {
			return "DATETIME"
		}
		return "TEXT"
	case reflect.Slice:
		if goType.Elem().Kind() == reflect.Uint8 {
			return "BLOB"
		}
		return "TEXT"
	default:
		return "TEXT"
	}
}

// GetPlaceholder returns the placeholder for parameterized queries
func (m *MySQLDialect) GetPlaceholder(index int) string {
	return "?"
}

// buildColumnDefinition builds a MySQL column definition
func (m *MySQLDialect) buildColumnDefinition(col interfaces.Column) string {
	var parts []string

	// Column name and type
	typeDef := col.Type
	if col.Length > 0 && (strings.Contains(typeDef, "VARCHAR") || strings.Contains(typeDef, "CHAR")) {
		typeDef = fmt.Sprintf("%s(%d)", typeDef, col.Length)
	}
	parts = append(parts, fmt.Sprintf("%s %s", col.Name, typeDef))

	// Nullable constraint
	if !col.Nullable {
		parts = append(parts, "NOT NULL")
	}

	// Auto increment
	if col.AutoIncrement {
		parts = append(parts, "AUTO_INCREMENT")
	}

	// Default value
	if col.Default != nil {
		defaultVal := fmt.Sprintf("%v", col.Default)
		if col.Type == "VARCHAR" || col.Type == "TEXT" {
			defaultVal = fmt.Sprintf("'%s'", defaultVal)
		}
		parts = append(parts, fmt.Sprintf("DEFAULT %s", defaultVal))
	}

	// Primary key
	if col.PrimaryKey {
		parts = append(parts, "PRIMARY KEY")
	}

	// Unique constraint
	if col.Unique {
		parts = append(parts, "UNIQUE")
	}

	// Foreign key
	if col.ForeignKey != nil {
		fkDef := fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(%s)",
			col.Name, col.ForeignKey.ReferencedTable, col.ForeignKey.ReferencedColumn)

		if col.ForeignKey.OnDelete != "" {
			fkDef += fmt.Sprintf(" ON DELETE %s", col.ForeignKey.OnDelete)
		}
		if col.ForeignKey.OnUpdate != "" {
			fkDef += fmt.Sprintf(" ON UPDATE %s", col.ForeignKey.OnUpdate)
		}

		parts = append(parts, fkDef)
	}

	return strings.Join(parts, " ")
}

// MySQLTransaction implements core.Transaction for MySQL
type MySQLTransaction struct {
	tx *sql.Tx
}

// Commit commits the transaction
func (mt *MySQLTransaction) Commit() error {
	return mt.tx.Commit()
}

// Rollback rolls back the transaction
func (mt *MySQLTransaction) Rollback() error {
	return mt.tx.Rollback()
}

// Exec executes a query within the transaction
func (mt *MySQLTransaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	return mt.tx.Exec(query, args...)
}

// Query executes a query that returns rows within the transaction
func (mt *MySQLTransaction) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return mt.tx.Query(query, args...)
}

// QueryRow executes a query that returns a single row within the transaction
func (mt *MySQLTransaction) QueryRow(query string, args ...interface{}) *sql.Row {
	return mt.tx.QueryRow(query, args...)
}

// Legacy function for backward compatibility
var DB *sql.DB

// InitMySQL initializes the legacy global DB variable
func InitMySQL() {
	dialect := NewMySQLDialect()
	config := interfaces.ConnectionConfig{
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     3306,
		Database: os.Getenv("MYSQL_DATABASE"),
		Username: os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
	}

	if err := dialect.Connect(config); err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	// Set the legacy DB variable for backward compatibility
	DB = dialect.db
}
