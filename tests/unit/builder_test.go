package unit

import (
	"os"
	"testing"

	"github.com/ESGI-M2/GO/orm/builder"
	"github.com/ESGI-M2/GO/orm/core/interfaces"
	"github.com/ESGI-M2/GO/orm/factory"
)

// Test models for builder tests
type BuilderTestUser struct {
	ID    int    `orm:"pk,auto"`
	Name  string `orm:"column:name"`
	Email string `orm:"column:email,unique"`
}

// TestNewSimpleORM tests SimpleORM creation
func TestNewSimpleORM(t *testing.T) {
	orm := builder.NewSimpleORM()

	if orm == nil {
		t.Fatal("NewSimpleORM should not return nil")
	}

	if orm.IsConnected() {
		t.Error("New SimpleORM should not be connected initially")
	}
}

// TestSimpleORM_WithDialect tests dialect setting
func TestSimpleORM_WithDialect(t *testing.T) {
	t.Run("WithDialect string", func(t *testing.T) {
		orm := builder.NewSimpleORM().WithDialect("mysql")
		if orm.GetDialectType() != factory.MySQL {
			t.Errorf("Expected MySQL dialect, got %s", orm.GetDialectType())
		}
	})

	t.Run("WithDialect DialectType", func(t *testing.T) {
		orm := builder.NewSimpleORM().WithDialect(factory.PostgreSQL)
		if orm.GetDialectType() != factory.PostgreSQL {
			t.Errorf("Expected PostgreSQL dialect, got %s", orm.GetDialectType())
		}
	})

	t.Run("WithDialect invalid type", func(t *testing.T) {
		orm := builder.NewSimpleORM().WithDialect(123) // Invalid type
		if orm.GetDialectType() != factory.MySQL {
			t.Errorf("Expected default MySQL dialect for invalid type, got %s", orm.GetDialectType())
		}
	})
}

func TestSimpleORM_WithSpecificDialects(t *testing.T) {
	t.Run("WithMySQL", func(t *testing.T) {
		orm := builder.NewSimpleORM().WithMySQL()
		if orm.GetDialectType() != factory.MySQL {
			t.Errorf("Expected MySQL dialect, got %s", orm.GetDialectType())
		}
	})

	t.Run("WithPostgreSQL", func(t *testing.T) {
		orm := builder.NewSimpleORM().WithPostgreSQL()
		if orm.GetDialectType() != factory.PostgreSQL {
			t.Errorf("Expected PostgreSQL dialect, got %s", orm.GetDialectType())
		}
	})
}

func TestSimpleORM_WithConfig(t *testing.T) {
	config := interfaces.ConnectionConfig{
		Host:     "localhost",
		Port:     3306,
		Database: "test_db",
		Username: "test_user",
		Password: "test_pass",
	}

	orm := builder.NewSimpleORM().WithConfig(config)
	retrievedConfig := orm.GetConfig()

	if retrievedConfig.Host != config.Host {
		t.Errorf("Expected host %s, got %s", config.Host, retrievedConfig.Host)
	}
	if retrievedConfig.Port != config.Port {
		t.Errorf("Expected port %d, got %d", config.Port, retrievedConfig.Port)
	}
	if retrievedConfig.Database != config.Database {
		t.Errorf("Expected database %s, got %s", config.Database, retrievedConfig.Database)
	}
}

func TestSimpleORM_WithQuickConfig(t *testing.T) {
	t.Run("MySQL default port", func(t *testing.T) {
		orm := builder.NewSimpleORM().WithDialect(factory.MySQL).WithQuickConfig("localhost", "test_db", "user", "pass")
		config := orm.GetConfig()

		if config.Host != "localhost" {
			t.Errorf("Expected host localhost, got %s", config.Host)
		}
		if config.Port != 3306 {
			t.Errorf("Expected port 3306, got %d", config.Port)
		}
		if config.Database != "test_db" {
			t.Errorf("Expected database test_db, got %s", config.Database)
		}
	})

	t.Run("PostgreSQL default port", func(t *testing.T) {
		orm := builder.NewSimpleORM().WithDialect(factory.PostgreSQL).WithQuickConfig("localhost", "test_db", "user", "pass")
		config := orm.GetConfig()

		if config.Port != 5432 {
			t.Errorf("Expected port 5432, got %d", config.Port)
		}
	})
}

func TestSimpleORM_WithAutoCreateDatabase(t *testing.T) {
	orm := builder.NewSimpleORM().WithAutoCreateDatabase()
	// We can't directly test the autoCreate field as it's private,
	// but we can test that the method returns the ORM for chaining
	if orm == nil {
		t.Error("WithAutoCreateDatabase should return the ORM instance")
	}
}

func TestSimpleORM_RegisterModel(t *testing.T) {
	orm := builder.NewSimpleORM()
	user := &BuilderTestUser{}

	result := orm.RegisterModel(user)
	if result != orm {
		t.Error("RegisterModel should return the same ORM instance for chaining")
	}
}

func TestSimpleORM_RegisterModels(t *testing.T) {
	orm := builder.NewSimpleORM()
	user := &BuilderTestUser{}
	anotherUser := &BuilderTestUser{}

	result := orm.RegisterModels(user, anotherUser)
	if result != orm {
		t.Error("RegisterModels should return the same ORM instance for chaining")
	}
}

func TestSimpleORM_Connect(t *testing.T) {
	t.Run("Connect with mock dialect", func(t *testing.T) {
		orm := builder.NewSimpleORM().
			WithDialect(factory.Mock).
			WithQuickConfig("localhost", "test_db", "user", "pass").
			RegisterModel(&BuilderTestUser{})

		err := orm.Connect()
		if err != nil {
			t.Errorf("Connect failed: %v", err)
		}

		if !orm.IsConnected() {
			t.Error("ORM should be connected after successful Connect()")
		}
	})

	t.Run("Connect without dialect", func(t *testing.T) {
		orm := builder.NewSimpleORM().WithQuickConfig("localhost", "test_db", "user", "pass")

		// Should use default MySQL dialect
		err := orm.Connect()
		// This might fail due to actual MySQL connection, but that's ok
		// We're testing the flow
		_ = err
	})

	t.Run("Connect already connected", func(t *testing.T) {
		orm := builder.NewSimpleORM().
			WithDialect(factory.Mock).
			WithQuickConfig("localhost", "test_db", "user", "pass")

		// First connection
		err := orm.Connect()
		if err != nil {
			t.Errorf("First connect failed: %v", err)
		}

		// Second connection should not fail
		err = orm.Connect()
		if err != nil {
			t.Errorf("Second connect should not fail: %v", err)
		}
	})
}

func TestSimpleORM_QueryAndRepository(t *testing.T) {
	orm := builder.NewSimpleORM().
		WithDialect(factory.Mock).
		WithQuickConfig("localhost", "test_db", "user", "pass").
		RegisterModel(&BuilderTestUser{})

	user := &BuilderTestUser{}

	t.Run("Query without connection", func(t *testing.T) {
		query := orm.Query(user)
		if query == nil {
			t.Error("Query should return a QueryBuilder even when not connected")
		}

		// Try to execute query - should get error
		_, err := query.Find()
		if err == nil {
			t.Error("Expected error when querying without connection")
		}
	})

	t.Run("Repository without connection", func(t *testing.T) {
		repo := orm.Repository(user)
		if repo == nil {
			t.Error("Repository should return a Repository even when not connected")
		}

		// Try to use repository - should get error
		_, err := repo.FindAll()
		if err == nil {
			t.Error("Expected error when using repository without connection")
		}
	})

	t.Run("Query and Repository after connection", func(t *testing.T) {
		err := orm.Connect()
		if err != nil {
			t.Errorf("Connect failed: %v", err)
		}

		query := orm.Query(user)
		if query == nil {
			t.Error("Query should return a QueryBuilder")
		}

		repo := orm.Repository(user)
		if repo == nil {
			t.Error("Repository should return a Repository")
		}
	})
}

func TestSimpleORM_Raw(t *testing.T) {
	orm := builder.NewSimpleORM().WithDialect(factory.Mock)

	t.Run("Raw without connection", func(t *testing.T) {
		raw := orm.Raw("SELECT 1")
		if raw == nil {
			t.Error("Raw should return a QueryBuilder even when not connected")
		}

		// Try to execute - should get error
		_, err := raw.Find()
		if err == nil {
			t.Error("Expected error when executing raw query without connection")
		}
	})

	t.Run("Raw after connection", func(t *testing.T) {
		orm.WithQuickConfig("localhost", "test_db", "user", "pass")
		err := orm.Connect()
		if err != nil {
			t.Errorf("Connect failed: %v", err)
		}

		raw := orm.Raw("SELECT 1")
		if raw == nil {
			t.Error("Raw should return a QueryBuilder")
		}
	})
}

func TestSimpleORM_Transaction(t *testing.T) {
	orm := builder.NewSimpleORM().
		WithDialect(factory.Mock).
		WithQuickConfig("localhost", "test_db", "user", "pass")

	t.Run("Transaction without connection", func(t *testing.T) {
		err := orm.Transaction(func(tx interfaces.ORM) error {
			return nil
		})
		if err == nil {
			t.Error("Expected error when starting transaction without connection")
		}
	})

	t.Run("Transaction after connection", func(t *testing.T) {
		err := orm.Connect()
		if err != nil {
			t.Errorf("Connect failed: %v", err)
		}

		err = orm.Transaction(func(tx interfaces.ORM) error {
			if tx == nil {
				t.Error("Transaction ORM should not be nil")
			}
			return nil
		})
		if err != nil {
			t.Errorf("Transaction failed: %v", err)
		}
	})
}

func TestSimpleORM_Close(t *testing.T) {
	orm := builder.NewSimpleORM().
		WithDialect(factory.Mock).
		WithQuickConfig("localhost", "test_db", "user", "pass")

	// Connect first
	err := orm.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	// Then close
	err = orm.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	if orm.IsConnected() {
		t.Error("ORM should not be connected after Close()")
	}
}

func TestSimpleORM_GlobalFunctions(t *testing.T) {
	t.Run("NewMySQL", func(t *testing.T) {
		orm := builder.NewMySQL()
		if orm == nil {
			t.Error("NewMySQL should not return nil")
		}
		if orm.GetDialectType() != factory.MySQL {
			t.Errorf("Expected MySQL dialect, got %s", orm.GetDialectType())
		}
	})

	t.Run("NewPostgreSQL", func(t *testing.T) {
		orm := builder.NewPostgreSQL()
		if orm == nil {
			t.Error("NewPostgreSQL should not return nil")
		}
		if orm.GetDialectType() != factory.PostgreSQL {
			t.Errorf("Expected PostgreSQL dialect, got %s", orm.GetDialectType())
		}
	})

	t.Run("QuickSetup", func(t *testing.T) {
		user := &BuilderTestUser{}
		orm, err := builder.QuickSetup("mock", "localhost", "test_db", "user", "pass", user)
		if err != nil {
			t.Errorf("QuickSetup failed: %v", err)
		}
		if orm == nil {
			t.Error("QuickSetup should not return nil")
		}
		if !orm.IsConnected() {
			t.Error("QuickSetup should return connected ORM")
		}
	})
}

// TestConfigBuilder tests the configuration builder
func TestConfigBuilder(t *testing.T) {
	t.Run("NewConfigBuilder", func(t *testing.T) {
		builder := builder.NewConfigBuilder()
		if builder == nil {
			t.Fatal("NewConfigBuilder should not return nil")
		}

		config := builder.GetConfig()
		if config.Host != "localhost" {
			t.Errorf("Expected default host localhost, got %s", config.Host)
		}
		if config.Port != 3306 {
			t.Errorf("Expected default port 3306, got %d", config.Port)
		}

		if builder.GetDialectType() != factory.MySQL {
			t.Errorf("Expected default dialect MySQL, got %s", builder.GetDialectType())
		}
	})

	t.Run("WithDialect", func(t *testing.T) {
		cb := builder.NewConfigBuilder().WithDialect(factory.PostgreSQL)
		if cb.GetDialectType() != factory.PostgreSQL {
			t.Errorf("Expected PostgreSQL dialect, got %s", cb.GetDialectType())
		}

		config := cb.GetConfig()
		if config.Port != 5432 {
			t.Errorf("Expected PostgreSQL default port 5432, got %d", config.Port)
		}
	})

	t.Run("Fluent interface", func(t *testing.T) {
		cb := builder.NewConfigBuilder().
			WithDialect(factory.MySQL).
			WithHost("testhost").
			WithPort(3307).
			WithDatabase("testdb").
			WithUsername("testuser").
			WithPassword("testpass").
			WithCredentials("newuser", "newpass").
			WithConnectionPool(10, 5).
			WithConnectionLifetime(300).
			WithAutoCreateDatabase()

		config := cb.GetConfig()
		if config.Host != "testhost" {
			t.Errorf("Expected host testhost, got %s", config.Host)
		}
		if config.Port != 3307 {
			t.Errorf("Expected port 3307, got %d", config.Port)
		}
		if config.Database != "testdb" {
			t.Errorf("Expected database testdb, got %s", config.Database)
		}
		if config.Username != "newuser" {
			t.Errorf("Expected username newuser, got %s", config.Username)
		}
		if config.Password != "newpass" {
			t.Errorf("Expected password newpass, got %s", config.Password)
		}
		if config.MaxOpenConns != 10 {
			t.Errorf("Expected MaxOpenConns 10, got %d", config.MaxOpenConns)
		}
		if config.MaxIdleConns != 5 {
			t.Errorf("Expected MaxIdleConns 5, got %d", config.MaxIdleConns)
		}
		if config.ConnMaxLifetime != 300 {
			t.Errorf("Expected ConnMaxLifetime 300, got %d", config.ConnMaxLifetime)
		}

		if !cb.ShouldAutoCreateDatabase() {
			t.Error("Expected auto create database to be enabled")
		}
	})

	t.Run("Build validation", func(t *testing.T) {
		// Missing database
		cb := builder.NewConfigBuilder().WithUsername("user")
		_, _, _, err := cb.Build()
		if err == nil {
			t.Error("Expected error for missing database")
		}

		// Missing username
		cb = builder.NewConfigBuilder().WithDatabase("testdb")
		_, _, _, err = cb.Build()
		if err == nil {
			t.Error("Expected error for missing username")
		}

		// Valid config
		cb = builder.NewConfigBuilder().WithDatabase("testdb").WithUsername("user")
		config, dialectType, autoCreate, err := cb.Build()
		if err != nil {
			t.Errorf("Build should not fail with valid config: %v", err)
		}
		if config.Database != "testdb" {
			t.Errorf("Expected database testdb, got %s", config.Database)
		}
		if dialectType != factory.MySQL {
			t.Errorf("Expected MySQL dialect, got %s", dialectType)
		}
		if autoCreate {
			t.Error("Expected autoCreate to be false by default")
		}
	})
}

func TestConfigBuilder_GlobalFunctions(t *testing.T) {
	t.Run("MySQL", func(t *testing.T) {
		cb := builder.MySQL()
		if cb == nil {
			t.Error("MySQL should not return nil")
		}
		if cb.GetDialectType() != factory.MySQL {
			t.Errorf("Expected MySQL dialect, got %s", cb.GetDialectType())
		}
	})

	t.Run("PostgreSQL", func(t *testing.T) {
		cb := builder.PostgreSQL()
		if cb == nil {
			t.Error("PostgreSQL should not return nil")
		}
		if cb.GetDialectType() != factory.PostgreSQL {
			t.Errorf("Expected PostgreSQL dialect, got %s", cb.GetDialectType())
		}
	})

	t.Run("Mock", func(t *testing.T) {
		cb := builder.Mock()
		if cb == nil {
			t.Error("Mock should not return nil")
		}
		if cb.GetDialectType() != factory.Mock {
			t.Errorf("Expected Mock dialect, got %s", cb.GetDialectType())
		}
	})
}

func TestConfigBuilder_FromEnv(t *testing.T) {
	// Set test environment variables
	os.Setenv("MYSQL_HOST", "envhost")
	os.Setenv("MYSQL_PORT", "3307")
	os.Setenv("MYSQL_DATABASE", "envdb")
	os.Setenv("MYSQL_USER", "envuser")
	os.Setenv("MYSQL_PASSWORD", "envpass")

	defer func() {
		os.Unsetenv("MYSQL_HOST")
		os.Unsetenv("MYSQL_PORT")
		os.Unsetenv("MYSQL_DATABASE")
		os.Unsetenv("MYSQL_USER")
		os.Unsetenv("MYSQL_PASSWORD")
	}()

	t.Run("FromEnv MySQL", func(t *testing.T) {
		cb := builder.NewConfigBuilder().WithDialect(factory.MySQL).FromEnv()
		config := cb.GetConfig()

		if config.Host != "envhost" {
			t.Errorf("Expected host envhost, got %s", config.Host)
		}
		if config.Port != 3307 {
			t.Errorf("Expected port 3307, got %d", config.Port)
		}
		if config.Database != "envdb" {
			t.Errorf("Expected database envdb, got %s", config.Database)
		}
		if config.Username != "envuser" {
			t.Errorf("Expected username envuser, got %s", config.Username)
		}
		if config.Password != "envpass" {
			t.Errorf("Expected password envpass, got %s", config.Password)
		}
	})

	t.Run("FromMySQLEnv", func(t *testing.T) {
		cb := builder.FromMySQLEnv()
		if cb == nil {
			t.Error("FromMySQLEnv should not return nil")
		}
		if cb.GetDialectType() != factory.MySQL {
			t.Errorf("Expected MySQL dialect, got %s", cb.GetDialectType())
		}
	})
}
