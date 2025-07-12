package builder

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ESGI-M2/GO/orm/core/interfaces"
	"github.com/ESGI-M2/GO/orm/factory"
	"github.com/joho/godotenv"
)

// ConfigBuilder provides a fluent interface for building ORM configuration
type ConfigBuilder struct {
	config      interfaces.ConnectionConfig
	dialectType factory.DialectType
	autoCreate  bool
	envLoaded   bool
}

// NewConfigBuilder creates a new configuration builder
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: interfaces.ConnectionConfig{
			Host: "localhost",
			Port: 3306,
		},
		dialectType: factory.MySQL,
		autoCreate:  false,
		envLoaded:   false,
	}
}

// WithDialect sets the database dialect
func (cb *ConfigBuilder) WithDialect(dialectType factory.DialectType) *ConfigBuilder {
	cb.dialectType = dialectType

	// Set default port based on dialect
	switch dialectType {
	case factory.MySQL:
		if cb.config.Port == 0 || cb.config.Port == 5432 {
			cb.config.Port = 3306
		}
	case factory.PostgreSQL, factory.Postgres:
		if cb.config.Port == 0 || cb.config.Port == 3306 {
			cb.config.Port = 5432
		}
	}

	return cb
}

// WithHost sets the database host
func (cb *ConfigBuilder) WithHost(host string) *ConfigBuilder {
	cb.config.Host = host
	return cb
}

// WithPort sets the database port
func (cb *ConfigBuilder) WithPort(port int) *ConfigBuilder {
	cb.config.Port = port
	return cb
}

// WithDatabase sets the database name
func (cb *ConfigBuilder) WithDatabase(database string) *ConfigBuilder {
	cb.config.Database = database
	return cb
}

// WithUsername sets the database username
func (cb *ConfigBuilder) WithUsername(username string) *ConfigBuilder {
	cb.config.Username = username
	return cb
}

// WithPassword sets the database password
func (cb *ConfigBuilder) WithPassword(password string) *ConfigBuilder {
	cb.config.Password = password
	return cb
}

// WithCredentials sets both username and password
func (cb *ConfigBuilder) WithCredentials(username, password string) *ConfigBuilder {
	cb.config.Username = username
	cb.config.Password = password
	return cb
}

// WithConnectionPool sets connection pool settings
func (cb *ConfigBuilder) WithConnectionPool(maxOpen, maxIdle int) *ConfigBuilder {
	cb.config.MaxOpenConns = maxOpen
	cb.config.MaxIdleConns = maxIdle
	return cb
}

// WithConnectionLifetime sets connection lifetime in seconds
func (cb *ConfigBuilder) WithConnectionLifetime(lifetimeSeconds int) *ConfigBuilder {
	cb.config.ConnMaxLifetime = lifetimeSeconds
	return cb
}

// WithAutoCreateDatabase enables automatic database creation
func (cb *ConfigBuilder) WithAutoCreateDatabase() *ConfigBuilder {
	cb.autoCreate = true
	return cb
}

// WithEnvFile loads configuration from environment variables
func (cb *ConfigBuilder) WithEnvFile(envFiles ...string) *ConfigBuilder {
	if len(envFiles) == 0 {
		envFiles = []string{".env", "../.env", "../../.env", ".env.local"}
	}

	for _, envFile := range envFiles {
		if err := godotenv.Load(envFile); err == nil {
			cb.envLoaded = true
			break
		}
	}

	return cb
}

// FromEnv loads configuration from environment variables based on dialect
func (cb *ConfigBuilder) FromEnv() *ConfigBuilder {
	if !cb.envLoaded {
		cb.WithEnvFile()
	}

	switch cb.dialectType {
	case factory.MySQL:
		cb.fromMySQLEnv()
	case factory.PostgreSQL, factory.Postgres:
		cb.fromPostgresEnv()
	}

	return cb
}

// fromMySQLEnv loads MySQL configuration from environment
func (cb *ConfigBuilder) fromMySQLEnv() {
	if host := os.Getenv("MYSQL_HOST"); host != "" {
		cb.config.Host = host
	}
	if portStr := os.Getenv("MYSQL_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			cb.config.Port = port
		}
	}
	if database := os.Getenv("MYSQL_DATABASE"); database != "" {
		cb.config.Database = database
	}
	if username := os.Getenv("MYSQL_USER"); username != "" {
		cb.config.Username = username
	}
	if password := os.Getenv("MYSQL_PASSWORD"); password != "" {
		cb.config.Password = password
	}
}

// fromPostgresEnv loads PostgreSQL configuration from environment
func (cb *ConfigBuilder) fromPostgresEnv() {
	if host := os.Getenv("POSTGRES_HOST"); host != "" {
		cb.config.Host = host
	}
	if portStr := os.Getenv("POSTGRES_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			cb.config.Port = port
		}
	}
	if database := os.Getenv("POSTGRES_DB"); database != "" {
		cb.config.Database = database
	}
	if username := os.Getenv("POSTGRES_USER"); username != "" {
		cb.config.Username = username
	}
	if password := os.Getenv("POSTGRES_PASSWORD"); password != "" {
		cb.config.Password = password
	}
}

// Build builds the configuration
func (cb *ConfigBuilder) Build() (interfaces.ConnectionConfig, factory.DialectType, bool, error) {
	if cb.config.Database == "" {
		return cb.config, cb.dialectType, cb.autoCreate, fmt.Errorf("database name is required")
	}

	if cb.config.Username == "" {
		return cb.config, cb.dialectType, cb.autoCreate, fmt.Errorf("username is required")
	}

	return cb.config, cb.dialectType, cb.autoCreate, nil
}

// GetConfig returns the current configuration
func (cb *ConfigBuilder) GetConfig() interfaces.ConnectionConfig {
	return cb.config
}

// GetDialectType returns the current dialect type
func (cb *ConfigBuilder) GetDialectType() factory.DialectType {
	return cb.dialectType
}

// ShouldAutoCreateDatabase returns whether to auto-create the database
func (cb *ConfigBuilder) ShouldAutoCreateDatabase() bool {
	return cb.autoCreate
}

// Convenience functions for common configurations

// MySQL creates a MySQL configuration builder
func MySQL() *ConfigBuilder {
	return NewConfigBuilder().WithDialect(factory.MySQL)
}

// PostgreSQL creates a PostgreSQL configuration builder
func PostgreSQL() *ConfigBuilder {
	return NewConfigBuilder().WithDialect(factory.PostgreSQL)
}

// Mock creates a mock configuration builder for testing
func Mock() *ConfigBuilder {
	return NewConfigBuilder().WithDialect(factory.Mock)
}

// FromMySQLEnv creates a MySQL configuration from environment variables
func FromMySQLEnv() *ConfigBuilder {
	return MySQL().FromEnv()
}

// FromPostgresEnv creates a PostgreSQL configuration from environment variables
func FromPostgresEnv() *ConfigBuilder {
	return PostgreSQL().FromEnv()
}
