package dialect

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ESGI-M2/GO/orm/core/interfaces"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// PostgresDialect implements the interfaces.Dialect interface for PostgreSQL
type PostgresDialect struct {
	db *sql.DB
}

// NewPostgresDialect creates a new Postgres dialect instance
func NewPostgresDialect() *PostgresDialect {
	return &PostgresDialect{}
}

// loadPostgresEnvFile attempts to load environment variables from .env files
func loadPostgresEnvFile() {
	envFiles := []string{".env", "../.env", "../../.env", ".env.local"}
	for _, envFile := range envFiles {
		if err := godotenv.Load(envFile); err == nil {
			break
		}
	}
}

// NewPostgresConnectionConfigFromEnv creates a connection config from environment variables
func NewPostgresConnectionConfigFromEnv() interfaces.ConnectionConfig {
	loadPostgresEnvFile()
	port := 5432 // default
	if portStr := os.Getenv("POSTGRES_PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}
	return interfaces.ConnectionConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     port,
		Database: os.Getenv("POSTGRES_DB"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	}
}

// Connect establishes a connection to PostgreSQL database
func (p *PostgresDialect) Connect(config interfaces.ConnectionConfig) error {
	var err error
	loadPostgresEnvFile()
	user := config.Username
	if user == "" {
		user = os.Getenv("POSTGRES_USER")
	}
	pass := config.Password
	if pass == "" {
		pass = os.Getenv("POSTGRES_PASSWORD")
	}
	host := config.Host
	if host == "" {
		host = os.Getenv("POSTGRES_HOST")
	}
	database := config.Database
	if database == "" {
		database = os.Getenv("POSTGRES_DB")
	}
	port := config.Port
	if port == 0 {
		if portStr := os.Getenv("POSTGRES_PORT"); portStr != "" {
			if p, err := strconv.Atoi(portStr); err == nil {
				port = p
			}
		} else {
			port = 5432
		}
	}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, database)
	p.db, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	p.db.SetMaxOpenConns(25)
	p.db.SetMaxIdleConns(5)
	p.db.SetConnMaxLifetime(5 * time.Minute)
	if err := p.Ping(); err != nil {
		return fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}
	log.Println("Connected to PostgreSQL successfully")
	return nil
}

func (p *PostgresDialect) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

func (p *PostgresDialect) Ping() error {
	if p.db == nil {
		return fmt.Errorf("database connection not established")
	}
	return p.db.Ping()
}

func (p *PostgresDialect) Exec(query string, args ...interface{}) (sql.Result, error) {
	if p.db == nil {
		return nil, fmt.Errorf("database connection not established")
	}
	return p.db.Exec(query, args...)
}

func (p *PostgresDialect) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if p.db == nil {
		return nil, fmt.Errorf("database connection not established")
	}
	return p.db.Query(query, args...)
}

func (p *PostgresDialect) QueryRow(query string, args ...interface{}) *sql.Row {
	if p.db == nil {
		return nil
	}
	return p.db.QueryRow(query, args...)
}

func (p *PostgresDialect) Begin() (interfaces.Transaction, error) {
	if p.db == nil {
		return nil, fmt.Errorf("database connection not established")
	}
	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}
	return &PostgresTransaction{tx: tx}, nil
}

func (p *PostgresDialect) BeginTx(ctx context.Context, opts *sql.TxOptions) (interfaces.Transaction, error) {
	if p.db == nil {
		return nil, fmt.Errorf("database connection not established")
	}
	tx, err := p.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &PostgresTransaction{tx: tx}, nil
}

func (p *PostgresDialect) CreateTable(tableName string, columns []interfaces.Column) error {
	var columnDefs []string
	for _, col := range columns {
		def := p.buildColumnDefinition(col)
		columnDefs = append(columnDefs, def)
	}
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n  %s\n)",
		tableName, strings.Join(columnDefs, ",\n  "))
	_, err := p.Exec(query)
	return err
}

func (p *PostgresDialect) DropTable(tableName string) error {
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	_, err := p.Exec(query)
	return err
}

func (p *PostgresDialect) TableExists(tableName string) (bool, error) {
	query := "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1"
	var count int
	err := p.QueryRow(query, tableName).Scan(&count)
	return count > 0, err
}

func (p *PostgresDialect) GetSQLType(goType reflect.Type) string {
	switch goType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return "INTEGER"
	case reflect.Int64:
		return "BIGINT"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return "INTEGER"
	case reflect.Uint64:
		return "BIGINT"
	case reflect.Float32:
		return "REAL"
	case reflect.Float64:
		return "DOUBLE PRECISION"
	case reflect.String:
		return "VARCHAR(255)"
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.Struct:
		if goType == reflect.TypeOf(time.Time{}) {
			return "TIMESTAMP"
		}
		return "TEXT"
	case reflect.Slice:
		if goType.Elem().Kind() == reflect.Uint8 {
			return "BYTEA"
		}
		return "TEXT"
	default:
		return "TEXT"
	}
}

func (p *PostgresDialect) GetPlaceholder(index int) string {
	return fmt.Sprintf("$%d", index+1)
}

func (p *PostgresDialect) FullTextSearch(field, query string) string {
	return fmt.Sprintf("to_tsvector('english', %s) @@ plainto_tsquery('english', '%s')", field, query)
}

func (p *PostgresDialect) GetRandomFunction() string {
	return "RANDOM()"
}

func (p *PostgresDialect) GetDateFunction() string {
	return "NOW()"
}

func (p *PostgresDialect) GetJSONExtract() string {
	return "jsonb_extract_path_text"
}

func (p *PostgresDialect) buildColumnDefinition(col interfaces.Column) string {
	var parts []string
	var typeDef string

	// Handle auto-increment for PostgreSQL
	if col.AutoIncrement {
		// Use SERIAL or BIGSERIAL depending on the Go type
		if col.Type == "BIGINT" {
			typeDef = "BIGSERIAL"
		} else {
			typeDef = "SERIAL"
		}
	} else {
		typeDef = col.Type
		if col.Length > 0 && (strings.Contains(typeDef, "VARCHAR") || strings.Contains(typeDef, "CHAR")) {
			typeDef = fmt.Sprintf("%s(%d)", typeDef, col.Length)
		}
	}

	parts = append(parts, fmt.Sprintf("%s %s", col.Name, typeDef))
	if !col.Nullable {
		parts = append(parts, "NOT NULL")
	}
	if col.Default != nil {
		defaultVal := fmt.Sprintf("%v", col.Default)
		if col.Type == "VARCHAR" || col.Type == "TEXT" {
			defaultVal = fmt.Sprintf("'%s'", defaultVal)
		}
		parts = append(parts, fmt.Sprintf("DEFAULT %s", defaultVal))
	}
	if col.PrimaryKey {
		parts = append(parts, "PRIMARY KEY")
	}
	if col.Unique {
		parts = append(parts, "UNIQUE")
	}
	if col.ForeignKey != nil {
		fkDef := fmt.Sprintf("REFERENCES %s(%s)",
			col.ForeignKey.ReferencedTable, col.ForeignKey.ReferencedColumn)
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

// PostgresTransaction implements core.Transaction for PostgreSQL
type PostgresTransaction struct {
	tx *sql.Tx
}

func (pt *PostgresTransaction) Commit() error {
	return pt.tx.Commit()
}

func (pt *PostgresTransaction) Rollback() error {
	return pt.tx.Rollback()
}

func (pt *PostgresTransaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	return pt.tx.Exec(query, args...)
}

func (pt *PostgresTransaction) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return pt.tx.Query(query, args...)
}

func (pt *PostgresTransaction) QueryRow(query string, args ...interface{}) *sql.Row {
	return pt.tx.QueryRow(query, args...)
}
