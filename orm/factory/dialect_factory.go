package factory

import (
	"fmt"
	"strings"

	"github.com/ESGI-M2/GO/dialect"
	"github.com/ESGI-M2/GO/orm/core/interfaces"
	ormDialect "github.com/ESGI-M2/GO/orm/dialect"
)

// DialectType represents supported database dialects
type DialectType string

const (
	MySQL    DialectType = "mysql"
	Postgres DialectType = "postgres"
	Mock     DialectType = "mock"
)

// DialectFactory provides easy dialect creation
type DialectFactory struct{}

// NewDialectFactory creates a new dialect factory instance
func NewDialectFactory() *DialectFactory {
	return &DialectFactory{}
}

// Create creates a dialect instance based on the dialect type
func (df *DialectFactory) Create(dialectType DialectType) (interfaces.Dialect, error) {
	switch strings.ToLower(string(dialectType)) {
	case "mysql":
		return dialect.NewMySQLDialect(), nil
	case "postgresql", "postgres":
		return dialect.NewPostgresDialect(), nil
	case "mock":
		return ormDialect.NewMockDialect(), nil
	default:
		return nil, fmt.Errorf("unsupported dialect type: %s", dialectType)
	}
}

// CreateFromString creates a dialect instance from a string
func (df *DialectFactory) CreateFromString(dialectStr string) (interfaces.Dialect, error) {
	return df.Create(DialectType(dialectStr))
}

// GetAvailableDialects returns a list of available dialect types
func (df *DialectFactory) GetAvailableDialects() []DialectType {
	return []DialectType{
		MySQL,
		Postgres,
		Mock,
	}
}

// IsSupported checks if a dialect type is supported
func (df *DialectFactory) IsSupported(dialectType DialectType) bool {
	for _, supported := range df.GetAvailableDialects() {
		if strings.EqualFold(string(supported), string(dialectType)) {
			return true
		}
	}
	return false
}

// Global factory instance for convenience
var DefaultFactory = NewDialectFactory()

// CreateDialect creates a dialect using the default factory
func CreateDialect(dialectType DialectType) (interfaces.Dialect, error) {
	return DefaultFactory.Create(dialectType)
}

// CreateDialectFromString creates a dialect from string using the default factory
func CreateDialectFromString(dialectStr string) (interfaces.Dialect, error) {
	return DefaultFactory.CreateFromString(dialectStr)
}
