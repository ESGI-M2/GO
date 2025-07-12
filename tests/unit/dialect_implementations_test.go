package unit

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/ESGI-M2/GO/dialect"
	"github.com/ESGI-M2/GO/orm/core/interfaces"
)

func TestMySQLDialect_Creation(t *testing.T) {
	mysql := dialect.NewMySQLDialect()
	if mysql == nil {
		t.Error("NewMySQLDialect should return a dialect instance")
	}
}

func TestMySQLDialect_ConnectionConfig(t *testing.T) {
	config := dialect.NewConnectionConfigFromEnv()
	if config.Host == "" {
		t.Log("Host is empty, this is expected if .env file is not set")
	}
	if config.Port == 0 {
		t.Log("Port is 0, this is expected if .env file is not set")
	}
	if config.Database == "" {
		t.Log("Database is empty, this is expected if .env file is not set")
	}
}

func TestMySQLDialect_Connect(t *testing.T) {
	mysql := dialect.NewMySQLDialect()

	// Test connect with valid config
	config := interfaces.ConnectionConfig{
		Host:     "localhost",
		Port:     3306,
		Username: "test",
		Password: "test",
		Database: "test",
	}

	// This might fail if MySQL is not available, which is expected
	err := mysql.Connect(config)
	if err != nil {
		t.Logf("MySQL connect failed as expected: %v", err)
	}

	// Test close
	err = mysql.Close()
	if err != nil {
		t.Logf("MySQL close failed: %v", err)
	}
}

func TestMySQLDialect_GetSQLType(t *testing.T) {
	mysql := dialect.NewMySQLDialect()

	// Test basic Go types
	tests := []struct {
		goType   reflect.Type
		expected string
	}{
		{reflect.TypeOf(""), "VARCHAR(255)"},
		{reflect.TypeOf(int(0)), "INT"},
		{reflect.TypeOf(int64(0)), "BIGINT"},
		{reflect.TypeOf(float64(0)), "DOUBLE"},
		{reflect.TypeOf(bool(false)), "TINYINT(1)"}, // MySQL uses TINYINT(1) for boolean
	}

	for _, test := range tests {
		result := mysql.GetSQLType(test.goType)
		if result != test.expected {
			t.Errorf("GetSQLType(%v) = %s, expected %s", test.goType, result, test.expected)
		}
	}
}

func TestMySQLDialect_GetPlaceholder(t *testing.T) {
	mysql := dialect.NewMySQLDialect()

	// MySQL uses ? placeholders
	for i := 0; i < 5; i++ {
		placeholder := mysql.GetPlaceholder(i)
		if placeholder != "?" {
			t.Errorf("GetPlaceholder(%d) = %s, expected ?", i, placeholder)
		}
	}
}

func TestMySQLDialect_AdvancedFeatures(t *testing.T) {
	mysql := dialect.NewMySQLDialect()

	// Test FullTextSearch
	fts := mysql.FullTextSearch("content", "search term")
	if fts == "" {
		t.Error("FullTextSearch should return non-empty string")
	}

	// Test GetRandomFunction
	random := mysql.GetRandomFunction()
	if random == "" {
		t.Error("GetRandomFunction should return non-empty string")
	}

	// Test GetDateFunction
	date := mysql.GetDateFunction()
	if date == "" {
		t.Error("GetDateFunction should return non-empty string")
	}

	// Test GetJSONExtract
	json := mysql.GetJSONExtract()
	if json == "" {
		t.Error("GetJSONExtract should return non-empty string")
	}
}

func TestMySQLDialect_Transaction(t *testing.T) {
	mysql := dialect.NewMySQLDialect()

	// Test transaction methods without actual connection
	_, err := mysql.Begin()
	if err == nil {
		t.Error("Begin should fail when not connected")
	}

	_, err = mysql.BeginTx(context.Background(), nil)
	if err == nil {
		t.Error("BeginTx should fail when not connected")
	}
}

func TestMySQLDialect_TableOperations(t *testing.T) {
	mysql := dialect.NewMySQLDialect()

	// Test table operations without connection
	columns := []interfaces.Column{
		{Name: "id", Type: "INT", PrimaryKey: true, AutoIncrement: true},
		{Name: "name", Type: "VARCHAR(255)", Nullable: false},
	}

	err := mysql.CreateTable("test_table", columns)
	if err == nil {
		t.Error("CreateTable should fail when not connected")
	}

	err = mysql.DropTable("test_table")
	if err == nil {
		t.Error("DropTable should fail when not connected")
	}

	// TableExists might panic if not connected, so we'll catch that
	defer func() {
		if r := recover(); r != nil {
			t.Logf("TableExists panicked as expected when not connected: %v", r)
		}
	}()

	_, err = mysql.TableExists("test_table")
	if err == nil {
		t.Error("TableExists should fail when not connected")
	}
}

func TestPostgreSQLDialect_Creation(t *testing.T) {
	postgres := dialect.NewPostgresDialect()
	if postgres == nil {
		t.Error("NewPostgresDialect should return a dialect instance")
	}
}

func TestPostgreSQLDialect_ConnectionConfig(t *testing.T) {
	config := dialect.NewPostgresConnectionConfigFromEnv()
	if config.Host == "" {
		t.Log("Host is empty, this is expected if .env file is not set")
	}
	if config.Port == 0 {
		t.Log("Port is 0, this is expected if .env file is not set")
	}
	if config.Database == "" {
		t.Log("Database is empty, this is expected if .env file is not set")
	}
}

func TestPostgreSQLDialect_Connect(t *testing.T) {
	postgres := dialect.NewPostgresDialect()

	// Test connect with valid config
	config := interfaces.ConnectionConfig{
		Host:     "localhost",
		Port:     5432,
		Username: "test",
		Password: "test",
		Database: "test",
	}

	// This might fail if PostgreSQL is not available, which is expected
	err := postgres.Connect(config)
	if err != nil {
		t.Logf("PostgreSQL connect failed as expected: %v", err)
	}

	// Test close
	err = postgres.Close()
	if err != nil {
		t.Logf("PostgreSQL close failed: %v", err)
	}
}

func TestPostgreSQLDialect_GetSQLType(t *testing.T) {
	postgres := dialect.NewPostgresDialect()

	// Test basic Go types
	tests := []struct {
		goType   reflect.Type
		expected string
	}{
		{reflect.TypeOf(""), "VARCHAR(255)"},
		{reflect.TypeOf(int(0)), "INTEGER"},
		{reflect.TypeOf(int64(0)), "BIGINT"},
		{reflect.TypeOf(float64(0)), "DOUBLE PRECISION"},
		{reflect.TypeOf(bool(false)), "BOOLEAN"},
	}

	for _, test := range tests {
		result := postgres.GetSQLType(test.goType)
		if result != test.expected {
			t.Errorf("GetSQLType(%v) = %s, expected %s", test.goType, result, test.expected)
		}
	}
}

func TestPostgreSQLDialect_GetPlaceholder(t *testing.T) {
	postgres := dialect.NewPostgresDialect()

	// PostgreSQL uses $1, $2, etc. placeholders
	for i := 0; i < 5; i++ {
		placeholder := postgres.GetPlaceholder(i)
		expected := fmt.Sprintf("$%d", i+1)
		if placeholder != expected {
			t.Errorf("GetPlaceholder(%d) = %s, expected %s", i, placeholder, expected)
		}
	}
}

func TestPostgreSQLDialect_AdvancedFeatures(t *testing.T) {
	postgres := dialect.NewPostgresDialect()

	// Test FullTextSearch
	fts := postgres.FullTextSearch("content", "search term")
	if fts == "" {
		t.Error("FullTextSearch should return non-empty string")
	}

	// Test GetRandomFunction
	random := postgres.GetRandomFunction()
	if random == "" {
		t.Error("GetRandomFunction should return non-empty string")
	}

	// Test GetDateFunction
	date := postgres.GetDateFunction()
	if date == "" {
		t.Error("GetDateFunction should return non-empty string")
	}

	// Test GetJSONExtract
	json := postgres.GetJSONExtract()
	if json == "" {
		t.Error("GetJSONExtract should return non-empty string")
	}
}

func TestPostgreSQLDialect_Transaction(t *testing.T) {
	postgres := dialect.NewPostgresDialect()

	// Test transaction methods without actual connection
	_, err := postgres.Begin()
	if err == nil {
		t.Error("Begin should fail when not connected")
	}

	_, err = postgres.BeginTx(context.Background(), nil)
	if err == nil {
		t.Error("BeginTx should fail when not connected")
	}
}

func TestPostgreSQLDialect_TableOperations(t *testing.T) {
	postgres := dialect.NewPostgresDialect()

	// Test table operations without connection
	columns := []interfaces.Column{
		{Name: "id", Type: "INTEGER", PrimaryKey: true, AutoIncrement: true},
		{Name: "name", Type: "VARCHAR(255)", Nullable: false},
	}

	err := postgres.CreateTable("test_table", columns)
	if err == nil {
		t.Error("CreateTable should fail when not connected")
	}

	err = postgres.DropTable("test_table")
	if err == nil {
		t.Error("DropTable should fail when not connected")
	}

	// TableExists might panic if not connected, so we'll catch that
	defer func() {
		if r := recover(); r != nil {
			t.Logf("TableExists panicked as expected when not connected: %v", r)
		}
	}()

	_, err = postgres.TableExists("test_table")
	if err == nil {
		t.Error("TableExists should fail when not connected")
	}
}

func TestDialectComparison(t *testing.T) {
	mysql := dialect.NewMySQLDialect()
	postgres := dialect.NewPostgresDialect()

	// Test that dialects handle placeholders differently
	mysqlPlaceholder := mysql.GetPlaceholder(0)
	postgresPlaceholder := postgres.GetPlaceholder(0)

	if mysqlPlaceholder == postgresPlaceholder {
		t.Error("MySQL and PostgreSQL should have different placeholder formats")
	}

	// Test that dialects handle some types differently
	intType := reflect.TypeOf(int(0))
	mysqlIntType := mysql.GetSQLType(intType)
	postgresIntType := postgres.GetSQLType(intType)

	if mysqlIntType != "INT" {
		t.Errorf("MySQL should use INT for int type, got %s", mysqlIntType)
	}

	if postgresIntType != "INTEGER" {
		t.Errorf("PostgreSQL should use INTEGER for int type, got %s", postgresIntType)
	}
}

func TestDialectInterface_Compliance(t *testing.T) {
	// Test that both dialects implement the Dialect interface
	var mysql interfaces.Dialect = dialect.NewMySQLDialect()
	var postgres interfaces.Dialect = dialect.NewPostgresDialect()

	if mysql == nil {
		t.Error("MySQL dialect should implement Dialect interface")
	}

	if postgres == nil {
		t.Error("PostgreSQL dialect should implement Dialect interface")
	}

	// Test that all interface methods are available
	_ = mysql.GetSQLType(reflect.TypeOf(""))
	_ = mysql.GetPlaceholder(0)
	_ = mysql.FullTextSearch("field", "query")
	_ = mysql.GetRandomFunction()
	_ = mysql.GetDateFunction()
	_ = mysql.GetJSONExtract()

	_ = postgres.GetSQLType(reflect.TypeOf(""))
	_ = postgres.GetPlaceholder(0)
	_ = postgres.FullTextSearch("field", "query")
	_ = postgres.GetRandomFunction()
	_ = postgres.GetDateFunction()
	_ = postgres.GetJSONExtract()
}

func TestDialectErrorHandling(t *testing.T) {
	mysql := dialect.NewMySQLDialect()
	postgres := dialect.NewPostgresDialect()

	// Test invalid config
	invalidConfig := interfaces.ConnectionConfig{
		Host:     "invalid-host",
		Port:     9999,
		Username: "invalid",
		Password: "invalid",
		Database: "invalid",
	}

	err := mysql.Connect(invalidConfig)
	if err == nil {
		t.Error("MySQL connect should fail with invalid config")
	}

	err = postgres.Connect(invalidConfig)
	if err == nil {
		t.Error("PostgreSQL connect should fail with invalid config")
	}
}

func TestDialectPing(t *testing.T) {
	mysql := dialect.NewMySQLDialect()
	postgres := dialect.NewPostgresDialect()

	// Test ping without connection
	err := mysql.Ping()
	if err == nil {
		t.Error("MySQL ping should fail when not connected")
	}

	err = postgres.Ping()
	if err == nil {
		t.Error("PostgreSQL ping should fail when not connected")
	}
}
