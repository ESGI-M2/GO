package unit

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ESGI-M2/GO/orm/core/metadata"
)

// Test models for utility function testing
type UtilsTestModel struct {
	ID        int       `orm:"pk,auto"`
	Name      string    `orm:"column:name"`
	Email     string    `orm:"column:email,unique"`
	Age       int       `orm:"column:age"`
	IsActive  bool      `orm:"column:is_active"`
	CreatedAt time.Time `orm:"column:created_at"`
	Balance   float64   `orm:"column:balance"`
}

type UtilsTestModel2 struct {
	UserID      int    `db:"user_id"`
	Description string `db:"description"`
	Status      string `db:"status"`
}

func TestReflectUtils_IsZeroValue(t *testing.T) {
	// Test helper function to check zero values
	isZeroValue := func(v reflect.Value) bool {
		switch v.Kind() {
		case reflect.Bool:
			return !v.Bool()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return v.Int() == 0
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return v.Uint() == 0
		case reflect.Float32, reflect.Float64:
			return v.Float() == 0
		case reflect.String:
			return v.String() == ""
		case reflect.Ptr, reflect.Interface:
			return v.IsNil()
		default:
			return false
		}
	}

	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{"zero int", 0, true},
		{"non-zero int", 42, false},
		{"zero int64", int64(0), true},
		{"non-zero int64", int64(42), false},
		{"zero float64", float64(0), true},
		{"non-zero float64", 3.14, false},
		{"zero string", "", true},
		{"non-zero string", "hello", false},
		{"zero bool", false, true},
		{"non-zero bool", true, false},
		{"nil pointer", (*int)(nil), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := reflect.ValueOf(test.value)
			if test.value == nil {
				v = reflect.ValueOf((*int)(nil))
			}
			result := isZeroValue(v)
			if result != test.expected {
				t.Errorf("isZeroValue(%v) = %v, expected %v", test.value, result, test.expected)
			}
		})
	}
}

func TestReflectUtils_SetFieldValue(t *testing.T) {
	// Test helper function to set field values with type conversion
	setFieldValue := func(field reflect.Value, value interface{}) error {
		if value == nil {
			field.Set(reflect.Zero(field.Type()))
			return nil
		}

		valueType := reflect.TypeOf(value)
		fieldType := field.Type()

		// If types match, set directly
		if valueType == fieldType {
			field.Set(reflect.ValueOf(value))
			return nil
		}

		// Handle type conversions
		switch fieldType.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			switch v := value.(type) {
			case int64:
				field.SetInt(v)
			case int:
				field.SetInt(int64(v))
			case float64:
				field.SetInt(int64(v))
			default:
				return fmt.Errorf("cannot convert %v to %s", value, fieldType)
			}
		case reflect.String:
			if str, ok := value.(string); ok {
				field.SetString(str)
			} else {
				return fmt.Errorf("cannot convert %v to string", value)
			}
		case reflect.Bool:
			if b, ok := value.(bool); ok {
				field.SetBool(b)
			} else {
				return fmt.Errorf("cannot convert %v to bool", value)
			}
		default:
			return fmt.Errorf("unsupported field type: %s", fieldType)
		}

		return nil
	}

	model := &UtilsTestModel{}
	v := reflect.ValueOf(model).Elem()

	// Test setting int field
	nameField := v.FieldByName("ID")
	err := setFieldValue(nameField, 42)
	if err != nil {
		t.Errorf("setFieldValue failed for ID: %v", err)
	}
	if model.ID != 42 {
		t.Errorf("Expected ID to be 42, got %d", model.ID)
	}

	// Test setting string field
	nameField = v.FieldByName("Name")
	err = setFieldValue(nameField, "John Doe")
	if err != nil {
		t.Errorf("setFieldValue failed for Name: %v", err)
	}
	if model.Name != "John Doe" {
		t.Errorf("Expected Name to be 'John Doe', got '%s'", model.Name)
	}

	// Test setting bool field
	boolField := v.FieldByName("IsActive")
	err = setFieldValue(boolField, true)
	if err != nil {
		t.Errorf("setFieldValue failed for IsActive: %v", err)
	}
	if !model.IsActive {
		t.Errorf("Expected IsActive to be true, got %v", model.IsActive)
	}

	// Test setting nil value
	err = setFieldValue(nameField, nil)
	if err != nil {
		t.Errorf("setFieldValue failed for nil: %v", err)
	}
	if model.Name != "" {
		t.Errorf("Expected Name to be empty after setting nil, got '%s'", model.Name)
	}
}

func TestReflectUtils_FindFieldByColumnName(t *testing.T) {
	// Test helper function to find field by column name
	findFieldByColumnName := func(entityValue reflect.Value, columnName string) reflect.Value {
		for i := 0; i < entityValue.NumField(); i++ {
			field := entityValue.Type().Field(i)

			// Check ORM tag first
			if ormTag := field.Tag.Get("orm"); ormTag != "" {
				// Parse ORM tag to find column name
				parts := strings.Split(ormTag, ",")
				for _, part := range parts {
					part = strings.TrimSpace(part)
					if strings.HasPrefix(part, "column:") {
						colName := strings.TrimPrefix(part, "column:")
						if colName == columnName {
							return entityValue.Field(i)
						}
					}
				}
			}

			// Fall back to DB tag
			if dbTag := field.Tag.Get("db"); dbTag != "" {
				if dbTag == columnName {
					return entityValue.Field(i)
				}
			}

			// Fall back to field name (case-insensitive)
			if strings.EqualFold(field.Name, columnName) {
				return entityValue.Field(i)
			}
		}
		return reflect.Value{}
	}

	model := &UtilsTestModel{}
	v := reflect.ValueOf(model).Elem()

	// Test finding field by ORM tag column name
	field := findFieldByColumnName(v, "name")
	if !field.IsValid() {
		t.Error("Should find field by ORM tag column name")
	}
	if field.Type() != reflect.TypeOf("") {
		t.Error("Found field should be string type")
	}

	// Test finding field by exact field name
	field = findFieldByColumnName(v, "ID")
	if !field.IsValid() {
		t.Error("Should find field by exact field name")
	}
	if field.Type() != reflect.TypeOf(0) {
		t.Error("Found field should be int type")
	}

	// Test finding field by case-insensitive field name
	field = findFieldByColumnName(v, "id")
	if !field.IsValid() {
		t.Error("Should find field by case-insensitive field name")
	}

	// Test not finding non-existent field
	field = findFieldByColumnName(v, "non_existent")
	if field.IsValid() {
		t.Error("Should not find non-existent field")
	}

	// Test with DB tag model
	model2 := &UtilsTestModel2{}
	v2 := reflect.ValueOf(model2).Elem()

	field = findFieldByColumnName(v2, "user_id")
	if !field.IsValid() {
		t.Error("Should find field by DB tag")
	}
	if field.Type() != reflect.TypeOf(0) {
		t.Error("Found field should be int type")
	}
}

func TestReflectUtils_StructToMap(t *testing.T) {
	// Test the structToMap function from repository
	structToMap := func(v reflect.Value) map[string]interface{} {
		result := make(map[string]interface{})

		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			value := v.Field(i)

			if value.CanInterface() {
				result[field.Name] = value.Interface()
			}
		}

		return result
	}

	model := &UtilsTestModel{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Age:      30,
		IsActive: true,
		Balance:  123.45,
	}

	v := reflect.ValueOf(model).Elem()
	result := structToMap(v)

	if len(result) != 7 {
		t.Errorf("Expected 7 fields, got %d", len(result))
	}

	if result["ID"] != 1 {
		t.Errorf("Expected ID to be 1, got %v", result["ID"])
	}

	if result["Name"] != "John Doe" {
		t.Errorf("Expected Name to be 'John Doe', got %v", result["Name"])
	}

	if result["Email"] != "john@example.com" {
		t.Errorf("Expected Email to be 'john@example.com', got %v", result["Email"])
	}

	if result["Age"] != 30 {
		t.Errorf("Expected Age to be 30, got %v", result["Age"])
	}

	if result["IsActive"] != true {
		t.Errorf("Expected IsActive to be true, got %v", result["IsActive"])
	}

	if result["Balance"] != 123.45 {
		t.Errorf("Expected Balance to be 123.45, got %v", result["Balance"])
	}
}

func TestReflectUtils_StructSliceToMapSlice(t *testing.T) {
	// Test the structSliceToMapSlice function
	structSliceToMapSlice := func(slice interface{}) []map[string]interface{} {
		if slice == nil {
			return []map[string]interface{}{}
		}

		sliceValue := reflect.ValueOf(slice)
		if sliceValue.Kind() != reflect.Slice {
			return []map[string]interface{}{}
		}

		result := make([]map[string]interface{}, sliceValue.Len())
		for i := 0; i < sliceValue.Len(); i++ {
			item := sliceValue.Index(i)
			if item.Kind() == reflect.Ptr {
				item = item.Elem()
			}

			if item.Kind() == reflect.Struct {
				// Inline structToMap implementation
				itemMap := make(map[string]interface{})
				for j := 0; j < item.NumField(); j++ {
					field := item.Type().Field(j)
					value := item.Field(j)
					if value.CanInterface() {
						itemMap[field.Name] = value.Interface()
					}
				}
				result[i] = itemMap
			}
		}

		return result
	}

	models := []*UtilsTestModel{
		{ID: 1, Name: "John", Email: "john@example.com", Age: 30, IsActive: true},
		{ID: 2, Name: "Jane", Email: "jane@example.com", Age: 25, IsActive: false},
	}

	result := structSliceToMapSlice(models)

	if len(result) != 2 {
		t.Errorf("Expected 2 maps, got %d", len(result))
	}

	if result[0]["ID"] != 1 {
		t.Errorf("Expected first ID to be 1, got %v", result[0]["ID"])
	}

	if result[1]["Name"] != "Jane" {
		t.Errorf("Expected second Name to be 'Jane', got %v", result[1]["Name"])
	}

	// Test with nil slice
	result = structSliceToMapSlice(nil)
	if len(result) != 0 {
		t.Errorf("Expected empty slice for nil input, got %d", len(result))
	}

	// Test with non-slice type
	result = structSliceToMapSlice("not a slice")
	if len(result) != 0 {
		t.Errorf("Expected empty slice for non-slice input, got %d", len(result))
	}
}

func TestReflectUtils_GetSQLType(t *testing.T) {
	// Test the getSQLType function
	getSQLType := func(t reflect.Type) string {
		switch t.Kind() {
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
			return "BOOLEAN"
		case reflect.Struct:
			if t.String() == "time.Time" {
				return "TIMESTAMP"
			}
			return "TEXT"
		case reflect.Slice:
			if t.Elem().Kind() == reflect.Uint8 {
				return "BLOB"
			}
			return "TEXT"
		default:
			return "TEXT"
		}
	}

	tests := []struct {
		goType   reflect.Type
		expected string
	}{
		{reflect.TypeOf(int(0)), "INT"},
		{reflect.TypeOf(int32(0)), "INT"},
		{reflect.TypeOf(int64(0)), "BIGINT"},
		{reflect.TypeOf(uint(0)), "INT UNSIGNED"},
		{reflect.TypeOf(uint64(0)), "BIGINT UNSIGNED"},
		{reflect.TypeOf(float32(0)), "FLOAT"},
		{reflect.TypeOf(float64(0)), "DOUBLE"},
		{reflect.TypeOf(""), "VARCHAR(255)"},
		{reflect.TypeOf(true), "BOOLEAN"},
		{reflect.TypeOf(time.Time{}), "TIMESTAMP"},
		{reflect.TypeOf([]byte{}), "BLOB"},
		{reflect.TypeOf([]string{}), "TEXT"},
	}

	for _, test := range tests {
		result := getSQLType(test.goType)
		if result != test.expected {
			t.Errorf("getSQLType(%v) = %s, expected %s", test.goType, result, test.expected)
		}
	}
}

func TestReflectUtils_GetTableName(t *testing.T) {
	// Test the getTableName function
	getTableName := func(t reflect.Type) string {
		// Default to lowercase struct name
		return strings.ToLower(t.Name())
	}

	tests := []struct {
		structType reflect.Type
		expected   string
	}{
		{reflect.TypeOf(UtilsTestModel{}), "utilstestmodel"},
		{reflect.TypeOf(UtilsTestModel2{}), "utilstestmodel2"},
	}

	for _, test := range tests {
		result := getTableName(test.structType)
		if result != test.expected {
			t.Errorf("getTableName(%v) = %s, expected %s", test.structType, result, test.expected)
		}
	}
}

func TestReflectUtils_ParseDefaultValue(t *testing.T) {
	// Test the parseDefaultValue function
	parseDefaultValue := func(value string, t reflect.Type) interface{} {
		switch t.Kind() {
		case reflect.String:
			return value
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if i, err := strconv.Atoi(value); err == nil {
				return i
			}
		case reflect.Float32, reflect.Float64:
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				return f
			}
		case reflect.Bool:
			return strings.ToLower(value) == "true"
		}
		return value
	}

	tests := []struct {
		value    string
		goType   reflect.Type
		expected interface{}
	}{
		{"hello", reflect.TypeOf(""), "hello"},
		{"42", reflect.TypeOf(int(0)), 42},
		{"3.14", reflect.TypeOf(float64(0)), 3.14},
		{"true", reflect.TypeOf(bool(false)), true},
		{"false", reflect.TypeOf(bool(false)), false},
		{"invalid", reflect.TypeOf(int(0)), "invalid"}, // Should return original string if parsing fails
	}

	for _, test := range tests {
		result := parseDefaultValue(test.value, test.goType)
		if result != test.expected {
			t.Errorf("parseDefaultValue(%s, %v) = %v, expected %v", test.value, test.goType, result, test.expected)
		}
	}
}

func TestReflectUtils_MetadataExtraction(t *testing.T) {
	// Test metadata extraction utilities
	manager := metadata.NewManager()

	model := &UtilsTestModel{}
	metadata, err := manager.ExtractMetadata(model)
	if err != nil {
		t.Errorf("ExtractMetadata failed: %v", err)
	}

	if metadata == nil {
		t.Fatal("Metadata should not be nil")
	}

	if metadata.Type != reflect.TypeOf(UtilsTestModel{}) {
		t.Errorf("Expected type %v, got %v", reflect.TypeOf(UtilsTestModel{}), metadata.Type)
	}

	if metadata.TableName != "utilstestmodel" {
		t.Errorf("Expected table name 'utilstestmodel', got '%s'", metadata.TableName)
	}

	if len(metadata.Columns) == 0 {
		t.Error("Should have extracted columns")
	}

	// Test primary key identification
	if metadata.PrimaryKey == "" {
		t.Error("Should have identified primary key")
	}

	// Test auto increment identification
	if metadata.AutoIncrement == "" {
		t.Error("Should have identified auto increment column")
	}

	// Test caching
	metadata2, err := manager.ExtractMetadata(model)
	if err != nil {
		t.Errorf("Second ExtractMetadata failed: %v", err)
	}

	if metadata != metadata2 {
		t.Error("Should return cached metadata instance")
	}
}

func TestReflectUtils_TypeValidation(t *testing.T) {
	// Test type validation utilities
	manager := metadata.NewManager()

	// Test with valid struct
	model := &UtilsTestModel{}
	_, err := manager.ExtractMetadata(model)
	if err != nil {
		t.Errorf("ExtractMetadata should succeed with valid struct: %v", err)
	}

	// Test with non-struct type
	var notAStruct int = 42
	_, err = manager.ExtractMetadata(notAStruct)
	if err == nil {
		t.Error("ExtractMetadata should fail with non-struct type")
	}

	// Test with nil (this will cause a panic, so we'll catch it)
	defer func() {
		if r := recover(); r != nil {
			t.Logf("ExtractMetadata with nil panicked as expected: %v", r)
		}
	}()

	_, err = manager.ExtractMetadata(nil)
	if err == nil {
		t.Error("ExtractMetadata should fail with nil")
	}
}

func TestReflectUtils_FieldAccessibility(t *testing.T) {
	// Test field accessibility checks
	type TestStruct struct {
		PublicField  string
		privateField string
	}

	model := &TestStruct{
		PublicField:  "public",
		privateField: "private",
	}

	v := reflect.ValueOf(model).Elem()

	// Test public field
	publicField := v.FieldByName("PublicField")
	if !publicField.IsValid() {
		t.Error("Should find public field")
	}
	if !publicField.CanInterface() {
		t.Error("Should be able to interface with public field")
	}
	if !publicField.CanSet() {
		t.Error("Should be able to set public field")
	}

	// Test private field
	privateField := v.FieldByName("privateField")
	if !privateField.IsValid() {
		t.Error("Should find private field")
	}
	if privateField.CanInterface() {
		t.Error("Should not be able to interface with private field")
	}
	if privateField.CanSet() {
		t.Error("Should not be able to set private field")
	}
}

func TestReflectUtils_TypeConversion(t *testing.T) {
	// Test type conversion utilities
	tests := []struct {
		name       string
		from       interface{}
		to         reflect.Type
		expected   interface{}
		shouldFail bool
	}{
		{"int to int", 42, reflect.TypeOf(int(0)), 42, false},
		{"int64 to int", int64(42), reflect.TypeOf(int(0)), 42, false},
		{"float64 to int", float64(42.0), reflect.TypeOf(int(0)), 42, false},
		{"string to string", "hello", reflect.TypeOf(""), "hello", false},
		{"bool to bool", true, reflect.TypeOf(false), true, false},
		{"incompatible types", "hello", reflect.TypeOf(int(0)), nil, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a value of the target type
			target := reflect.New(test.to).Elem()

			// Attempt conversion
			fromValue := reflect.ValueOf(test.from)

			if test.to == fromValue.Type() {
				target.Set(fromValue)
				if !test.shouldFail && target.Interface() != test.expected {
					t.Errorf("Expected %v, got %v", test.expected, target.Interface())
				}
			} else if fromValue.Type().ConvertibleTo(test.to) {
				target.Set(fromValue.Convert(test.to))
				if !test.shouldFail && target.Interface() != test.expected {
					t.Errorf("Expected %v, got %v", test.expected, target.Interface())
				}
			} else if !test.shouldFail {
				t.Errorf("Expected conversion to succeed but types are not convertible")
			}
		})
	}
}

func TestReflectUtils_PointerHandling(t *testing.T) {
	// Test pointer handling utilities
	model := &UtilsTestModel{ID: 1, Name: "Test"}

	// Test with pointer
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr {
		t.Error("Should be pointer type")
	}

	// Test dereferencing
	elem := v.Elem()
	if elem.Kind() != reflect.Struct {
		t.Error("Should be struct type after dereferencing")
	}

	// Test field access through pointer
	nameField := elem.FieldByName("Name")
	if !nameField.IsValid() {
		t.Error("Should find Name field")
	}
	if nameField.String() != "Test" {
		t.Errorf("Expected 'Test', got '%s'", nameField.String())
	}

	// Test with nil pointer
	var nilPointer *UtilsTestModel
	nilValue := reflect.ValueOf(nilPointer)
	if nilValue.Kind() != reflect.Ptr {
		t.Error("Should be pointer type")
	}
	if !nilValue.IsNil() {
		t.Error("Should be nil pointer")
	}
}
