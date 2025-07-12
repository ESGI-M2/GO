package unit

import (
	"strings"
	"testing"

	"github.com/ESGI-M2/GO/orm/core/interfaces"
	"github.com/ESGI-M2/GO/orm/factory"
)

// TestDialectFactory tests the dialect factory functionality
func TestDialectFactory(t *testing.T) {
	df := factory.NewDialectFactory()

	if df == nil {
		t.Fatal("NewDialectFactory should not return nil")
	}

	t.Run("Create MySQL Dialect", func(t *testing.T) {
		dialect, err := df.Create(factory.MySQL)
		if err != nil {
			t.Errorf("Create MySQL dialect failed: %v", err)
		}
		if dialect == nil {
			t.Error("MySQL dialect should not be nil")
		}
	})

	t.Run("Create PostgreSQL Dialect", func(t *testing.T) {
		dialect, err := df.Create(factory.PostgreSQL)
		if err != nil {
			t.Errorf("Create PostgreSQL dialect failed: %v", err)
		}
		if dialect == nil {
			t.Error("PostgreSQL dialect should not be nil")
		}
	})

	t.Run("Create Postgres Dialect (alias)", func(t *testing.T) {
		dialect, err := df.Create(factory.Postgres)
		if err != nil {
			t.Errorf("Create Postgres dialect failed: %v", err)
		}
		if dialect == nil {
			t.Error("Postgres dialect should not be nil")
		}
	})

	t.Run("Create Mock Dialect", func(t *testing.T) {
		dialect, err := df.Create(factory.Mock)
		if err != nil {
			t.Errorf("Create Mock dialect failed: %v", err)
		}
		if dialect == nil {
			t.Error("Mock dialect should not be nil")
		}
	})

	t.Run("Create Unsupported Dialect", func(t *testing.T) {
		dialect, err := df.Create("unsupported")
		if err == nil {
			t.Error("Expected error for unsupported dialect")
		}
		if dialect != nil {
			t.Error("Unsupported dialect should return nil")
		}
		if !strings.Contains(err.Error(), "unsupported dialect type") {
			t.Errorf("Expected unsupported dialect error message, got: %v", err)
		}
	})
}

func TestDialectFactory_CreateFromString(t *testing.T) {
	df := factory.NewDialectFactory()

	t.Run("Create from string - mysql", func(t *testing.T) {
		dialect, err := df.CreateFromString("mysql")
		if err != nil {
			t.Errorf("CreateFromString mysql failed: %v", err)
		}
		if dialect == nil {
			t.Error("MySQL dialect should not be nil")
		}
	})

	t.Run("Create from string - postgresql", func(t *testing.T) {
		dialect, err := df.CreateFromString("postgresql")
		if err != nil {
			t.Errorf("CreateFromString postgresql failed: %v", err)
		}
		if dialect == nil {
			t.Error("PostgreSQL dialect should not be nil")
		}
	})

	t.Run("Create from string - postgres", func(t *testing.T) {
		dialect, err := df.CreateFromString("postgres")
		if err != nil {
			t.Errorf("CreateFromString postgres failed: %v", err)
		}
		if dialect == nil {
			t.Error("Postgres dialect should not be nil")
		}
	})

	t.Run("Create from string - mock", func(t *testing.T) {
		dialect, err := df.CreateFromString("mock")
		if err != nil {
			t.Errorf("CreateFromString mock failed: %v", err)
		}
		if dialect == nil {
			t.Error("Mock dialect should not be nil")
		}
	})

	t.Run("Create from string - case insensitive", func(t *testing.T) {
		dialect, err := df.CreateFromString("MYSQL")
		if err != nil {
			t.Errorf("CreateFromString MYSQL failed: %v", err)
		}
		if dialect == nil {
			t.Error("MySQL dialect should not be nil")
		}
	})

	t.Run("Create from string - unsupported", func(t *testing.T) {
		dialect, err := df.CreateFromString("oracle")
		if err == nil {
			t.Error("Expected error for unsupported dialect")
		}
		if dialect != nil {
			t.Error("Unsupported dialect should return nil")
		}
	})
}

func TestDialectFactory_GetAvailableDialects(t *testing.T) {
	df := factory.NewDialectFactory()
	dialects := df.GetAvailableDialects()

	if len(dialects) == 0 {
		t.Error("GetAvailableDialects should return at least one dialect")
	}

	// Check that all expected dialects are present
	expectedDialects := []factory.DialectType{
		factory.MySQL,
		factory.Postgres,
		factory.Mock,
	}

	for _, expected := range expectedDialects {
		found := false
		for _, available := range dialects {
			if available == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected dialect %s not found in available dialects", expected)
		}
	}
}

func TestDialectFactory_IsSupported(t *testing.T) {
	df := factory.NewDialectFactory()

	t.Run("Supported dialects", func(t *testing.T) {
		supportedDialects := []factory.DialectType{
			factory.MySQL,
			factory.PostgreSQL,
			factory.Postgres,
			factory.SQLite,
			factory.Mock,
		}

		for _, dialect := range supportedDialects {
			if !df.IsSupported(dialect) {
				t.Errorf("Dialect %s should be supported", dialect)
			}
		}
	})

	t.Run("Case insensitive support", func(t *testing.T) {
		if !df.IsSupported("MYSQL") {
			t.Error("MYSQL should be supported (case insensitive)")
		}
		if !df.IsSupported("mysql") {
			t.Error("mysql should be supported")
		}
		if !df.IsSupported("MySQL") {
			t.Error("MySQL should be supported")
		}
	})

	t.Run("Unsupported dialects", func(t *testing.T) {
		unsupportedDialects := []factory.DialectType{
			"oracle",
			"sqlite3",
			"mongodb",
			"redis",
		}

		for _, dialect := range unsupportedDialects {
			if df.IsSupported(dialect) {
				t.Errorf("Dialect %s should not be supported", dialect)
			}
		}
	})
}

func TestDialectFactory_GlobalFunctions(t *testing.T) {
	t.Run("CreateDialect", func(t *testing.T) {
		dialect, err := factory.CreateDialect(factory.Mock)
		if err != nil {
			t.Errorf("CreateDialect failed: %v", err)
		}
		if dialect == nil {
			t.Error("CreateDialect should not return nil")
		}
	})

	t.Run("CreateDialectFromString", func(t *testing.T) {
		dialect, err := factory.CreateDialectFromString("mock")
		if err != nil {
			t.Errorf("CreateDialectFromString failed: %v", err)
		}
		if dialect == nil {
			t.Error("CreateDialectFromString should not return nil")
		}
	})
}

func TestDatabaseCreator(t *testing.T) {
	creator := factory.NewDatabaseCreator()

	if creator == nil {
		t.Fatal("NewDatabaseCreator should not return nil")
	}

	t.Run("CreateDatabaseIfNotExists Mock", func(t *testing.T) {
		config := interfaces.ConnectionConfig{
			Host:     "localhost",
			Port:     3306,
			Username: "test",
			Password: "test",
			Database: "test_db",
		}

		// Mock should not fail
		err := creator.CreateDatabaseIfNotExists(config, factory.Mock)
		if err != nil {
			t.Errorf("CreateDatabaseIfNotExists for mock should not fail: %v", err)
		}
	})

	t.Run("CreateDatabaseIfNotExists Unsupported", func(t *testing.T) {
		config := interfaces.ConnectionConfig{
			Host:     "localhost",
			Port:     3306,
			Username: "test",
			Password: "test",
			Database: "test_db",
		}

		err := creator.CreateDatabaseIfNotExists(config, "unsupported")
		if err == nil {
			t.Error("Expected error for unsupported dialect")
		}
		if !strings.Contains(err.Error(), "database creation not supported") {
			t.Errorf("Expected unsupported dialect error, got: %v", err)
		}
	})

	t.Run("EnsureDatabaseExists Mock", func(t *testing.T) {
		config := interfaces.ConnectionConfig{
			Host:     "localhost",
			Port:     3306,
			Username: "test",
			Password: "test",
			Database: "test_db",
		}

		// Mock should not fail
		err := creator.EnsureDatabaseExists(config, factory.Mock)
		if err != nil {
			t.Errorf("EnsureDatabaseExists for mock should not fail: %v", err)
		}
	})
}

func TestDatabaseCreator_GlobalFunctions(t *testing.T) {
	config := interfaces.ConnectionConfig{
		Host:     "localhost",
		Port:     3306,
		Username: "test",
		Password: "test",
		Database: "test_db",
	}

	t.Run("CreateDatabaseIfNotExists", func(t *testing.T) {
		err := factory.CreateDatabaseIfNotExists(config, factory.Mock)
		if err != nil {
			t.Errorf("CreateDatabaseIfNotExists failed: %v", err)
		}
	})

	t.Run("EnsureDatabaseExists", func(t *testing.T) {
		err := factory.EnsureDatabaseExists(config, factory.Mock)
		if err != nil {
			t.Errorf("EnsureDatabaseExists failed: %v", err)
		}
	})
}

// Test dialect type constants
func TestDialectTypeConstants(t *testing.T) {
	t.Run("Dialect type values", func(t *testing.T) {
		if factory.MySQL != "mysql" {
			t.Errorf("Expected MySQL to be 'mysql', got %s", factory.MySQL)
		}
		if factory.Postgres != "postgres" {
			t.Errorf("Expected Postgres to be 'postgres', got %s", factory.Postgres)
		}
		if factory.Mock != "mock" {
			t.Errorf("Expected Mock to be 'mock', got %s", factory.Mock)
		}
	})
}
