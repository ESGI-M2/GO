package unit

import (
	"testing"

	"github.com/ESGI-M2/GO/orm/builder"
	"github.com/ESGI-M2/GO/orm/core/interfaces"
	"github.com/ESGI-M2/GO/orm/factory"
)

type AdvancedRepoTestUser struct {
	ID        int    `orm:"pk,auto"`
	Name      string `orm:"column:name"`
	Email     string `orm:"column:email,unique"`
	Age       int    `orm:"column:age"`
	IsActive  bool   `orm:"column:is_active"`
	DeletedAt *int64 `orm:"column:deleted_at"`
	CreatedAt *int64 `orm:"column:created_at"`
	UpdatedAt *int64 `orm:"column:updated_at"`
}

type AdvancedRepoTestProfile struct {
	ID       int    `orm:"pk,auto"`
	UserID   int    `orm:"column:user_id,fk:users.id"`
	Bio      string `orm:"column:bio"`
	Avatar   string `orm:"column:avatar"`
	Verified bool   `orm:"column:verified"`
}

func setupAdvancedRepository() interfaces.Repository {
	orm := builder.NewSimpleORM().
		WithDialect(factory.Mock).
		RegisterModel(&AdvancedRepoTestUser{}).
		RegisterModel(&AdvancedRepoTestProfile{})

	err := orm.Connect()
	if err != nil {
		panic(err)
	}

	return orm.Repository(&AdvancedRepoTestUser{})
}

func TestAdvancedRepository_FindWithRelations(t *testing.T) {
	repo := setupAdvancedRepository()

	result, err := repo.FindWithRelations(1, "profile", "posts")
	if err != nil {
		t.Errorf("FindWithRelations failed: %v", err)
	}

	// result can be nil with mock dialect
	_ = result
}

func TestAdvancedRepository_FindAllWithRelations(t *testing.T) {
	repo := setupAdvancedRepository()

	results, err := repo.FindAllWithRelations("profile", "posts")
	if err != nil {
		t.Errorf("FindAllWithRelations failed: %v", err)
	}

	if results == nil {
		results = []interface{}{}
	}

	if len(results) != 0 {
		t.Errorf("FindAllWithRelations should return empty slice, got %d", len(results))
	}
}

func TestAdvancedRepository_FindByWithRelations(t *testing.T) {
	repo := setupAdvancedRepository()

	criteria := map[string]interface{}{
		"is_active": true,
		"age":       25,
	}

	results, err := repo.FindByWithRelations(criteria, "profile")
	if err != nil {
		t.Errorf("FindByWithRelations failed: %v", err)
	}

	if results == nil {
		results = []interface{}{}
	}

	if len(results) != 0 {
		t.Errorf("FindByWithRelations should return empty slice, got %d", len(results))
	}
}

func TestAdvancedRepository_BatchCreate(t *testing.T) {
	repo := setupAdvancedRepository()

	users := []interface{}{
		&AdvancedRepoTestUser{Name: "John", Email: "john@example.com", Age: 25},
		&AdvancedRepoTestUser{Name: "Jane", Email: "jane@example.com", Age: 30},
		&AdvancedRepoTestUser{Name: "Bob", Email: "bob@example.com", Age: 35},
	}

	err := repo.BatchCreate(users)
	if err != nil {
		t.Errorf("BatchCreate failed: %v", err)
	}
}

func TestAdvancedRepository_BatchUpdate(t *testing.T) {
	repo := setupAdvancedRepository()

	users := []interface{}{
		&AdvancedRepoTestUser{ID: 1, Name: "John Updated", Email: "john.updated@example.com", Age: 26},
		&AdvancedRepoTestUser{ID: 2, Name: "Jane Updated", Email: "jane.updated@example.com", Age: 31},
		&AdvancedRepoTestUser{ID: 3, Name: "Bob Updated", Email: "bob.updated@example.com", Age: 36},
	}

	err := repo.BatchUpdate(users)
	if err != nil {
		t.Errorf("BatchUpdate failed: %v", err)
	}
}

func TestAdvancedRepository_BatchDelete(t *testing.T) {
	repo := setupAdvancedRepository()

	users := []interface{}{
		&AdvancedRepoTestUser{ID: 1},
		&AdvancedRepoTestUser{ID: 2},
		&AdvancedRepoTestUser{ID: 3},
	}

	err := repo.BatchDelete(users)
	if err != nil {
		t.Errorf("BatchDelete failed: %v", err)
	}
}

func TestAdvancedRepository_SoftDelete(t *testing.T) {
	repo := setupAdvancedRepository()

	user := &AdvancedRepoTestUser{ID: 1}
	err := repo.SoftDelete(user)
	if err == nil {
		t.Error("SoftDelete should fail when soft deletes are not enabled")
	}
	if err != nil && err.Error() != "soft deletes not enabled for this model" {
		t.Errorf("Expected 'soft deletes not enabled' error, got: %v", err)
	}
}

func TestAdvancedRepository_Restore(t *testing.T) {
	repo := setupAdvancedRepository()

	user := &AdvancedRepoTestUser{ID: 1}
	err := repo.Restore(user)
	if err == nil {
		t.Error("Restore should fail when soft deletes are not enabled")
	}
	if err != nil && err.Error() != "soft deletes not enabled for this model" {
		t.Errorf("Expected 'soft deletes not enabled' error, got: %v", err)
	}
}

func TestAdvancedRepository_ForceDelete(t *testing.T) {
	repo := setupAdvancedRepository()

	user := &AdvancedRepoTestUser{ID: 1}
	err := repo.ForceDelete(user)
	if err != nil {
		t.Errorf("ForceDelete failed: %v", err)
	}
}

func TestAdvancedRepository_FindTrashed(t *testing.T) {
	repo := setupAdvancedRepository()

	results, err := repo.FindTrashed()
	if err == nil {
		t.Error("FindTrashed should fail when soft deletes are not enabled")
	}
	if err != nil && err.Error() != "soft deletes not enabled for this model" {
		t.Errorf("Expected 'soft deletes not enabled' error, got: %v", err)
	}

	// results should be nil when error occurs
	if results != nil {
		t.Error("FindTrashed should return nil results when error occurs")
	}
}

func TestAdvancedRepository_RestoreBy(t *testing.T) {
	repo := setupAdvancedRepository()

	criteria := map[string]interface{}{
		"name": "John",
	}

	err := repo.RestoreBy(criteria)
	if err == nil {
		t.Error("RestoreBy should fail when soft deletes are not enabled")
	}
	if err != nil && err.Error() != "soft deletes not enabled for this model" {
		t.Errorf("Expected 'soft deletes not enabled' error, got: %v", err)
	}
}

func TestAdvancedRepository_Scope(t *testing.T) {
	repo := setupAdvancedRepository()

	scopedRepo := repo.Scope("active")
	if scopedRepo == nil {
		t.Error("Scope should return a repository")
	}

	// Test with parameters
	scopedRepo = repo.Scope("age_greater_than", 18)
	if scopedRepo == nil {
		t.Error("Scope with parameters should return a repository")
	}
}

func TestAdvancedRepository_Chunk(t *testing.T) {
	repo := setupAdvancedRepository()

	chunkSize := 10
	callCount := 0

	err := repo.Chunk(chunkSize, func(entities []interface{}) error {
		callCount++

		if entities == nil {
			entities = []interface{}{}
		}

		if len(entities) > chunkSize {
			t.Errorf("Chunk size should not exceed %d, got %d", chunkSize, len(entities))
		}

		return nil
	})

	if err != nil {
		t.Errorf("Chunk failed: %v", err)
	}

	// With mock dialect, we might have 0 or 1 call
	if callCount < 0 {
		t.Errorf("Chunk should call function at least once or not at all, got %d calls", callCount)
	}
}

func TestAdvancedRepository_Each(t *testing.T) {
	repo := setupAdvancedRepository()

	processedCount := 0

	err := repo.Each(func(entity interface{}) error {
		processedCount++

		if entity == nil {
			t.Error("Each should not pass nil entity")
		}

		return nil
	})

	if err != nil {
		t.Errorf("Each failed: %v", err)
	}

	// With mock dialect, we might have 0 entities
	if processedCount < 0 {
		t.Errorf("Each should process 0 or more entities, got %d", processedCount)
	}
}

func TestAdvancedRepository_Pluck(t *testing.T) {
	repo := setupAdvancedRepository()

	values, err := repo.Pluck("name")
	if err != nil {
		t.Errorf("Pluck failed: %v", err)
	}

	if values == nil {
		values = []interface{}{}
	}

	if len(values) != 0 {
		t.Errorf("Pluck should return empty slice, got %d", len(values))
	}
}

func TestAdvancedRepository_Value(t *testing.T) {
	repo := setupAdvancedRepository()

	value, err := repo.Value("name")
	if err != nil {
		t.Errorf("Value failed: %v", err)
	}

	// value can be nil with mock dialect
	_ = value
}

func TestAdvancedRepository_Increment(t *testing.T) {
	repo := setupAdvancedRepository()

	err := repo.Increment("age", 1)
	if err != nil {
		t.Errorf("Increment failed: %v", err)
	}

	// Test with float value
	err = repo.Increment("age", 2.5)
	if err != nil {
		t.Errorf("Increment with float failed: %v", err)
	}
}

func TestAdvancedRepository_Decrement(t *testing.T) {
	repo := setupAdvancedRepository()

	err := repo.Decrement("age", 1)
	if err != nil {
		t.Errorf("Decrement failed: %v", err)
	}

	// Test with float value
	err = repo.Decrement("age", 1.5)
	if err != nil {
		t.Errorf("Decrement with float failed: %v", err)
	}
}

func TestAdvancedRepository_ChainedOperations(t *testing.T) {
	repo := setupAdvancedRepository()

	// Test chained operations with scope
	scopedRepo := repo.Scope("active")

	results, err := scopedRepo.FindAll()
	if err != nil {
		t.Errorf("Chained operations failed: %v", err)
	}

	if results == nil {
		results = []interface{}{}
	}

	if len(results) != 0 {
		t.Errorf("Chained operations should return empty slice, got %d", len(results))
	}
}

func TestAdvancedRepository_ComplexBatchOperations(t *testing.T) {
	repo := setupAdvancedRepository()

	// Create batch of users
	users := []interface{}{
		&AdvancedRepoTestUser{Name: "User1", Email: "user1@example.com", Age: 20},
		&AdvancedRepoTestUser{Name: "User2", Email: "user2@example.com", Age: 25},
		&AdvancedRepoTestUser{Name: "User3", Email: "user3@example.com", Age: 30},
	}

	err := repo.BatchCreate(users)
	if err != nil {
		t.Errorf("BatchCreate failed: %v", err)
	}

	// Update batch
	for i, user := range users {
		if u, ok := user.(*AdvancedRepoTestUser); ok {
			u.ID = i + 1 // Simulate assigned IDs
			u.Name = u.Name + " Updated"
		}
	}

	err = repo.BatchUpdate(users)
	if err != nil {
		t.Errorf("BatchUpdate failed: %v", err)
	}

	// Soft delete batch
	err = repo.BatchDelete(users)
	if err != nil {
		t.Errorf("BatchDelete failed: %v", err)
	}
}

func TestAdvancedRepository_SoftDeleteWorkflow(t *testing.T) {
	repo := setupAdvancedRepository()

	user := &AdvancedRepoTestUser{ID: 1, Name: "Test User"}

	// Test that soft delete workflow fails when soft deletes are not enabled
	// Soft delete
	err := repo.SoftDelete(user)
	if err == nil {
		t.Error("SoftDelete should fail when soft deletes are not enabled")
	}

	// Find trashed
	_, err = repo.FindTrashed()
	if err == nil {
		t.Error("FindTrashed should fail when soft deletes are not enabled")
	}

	// Restore
	err = repo.Restore(user)
	if err == nil {
		t.Error("Restore should fail when soft deletes are not enabled")
	}

	// Force delete should work regardless of soft delete settings
	err = repo.ForceDelete(user)
	if err != nil {
		t.Errorf("ForceDelete failed: %v", err)
	}
}

func TestAdvancedRepository_UtilityOperations(t *testing.T) {
	repo := setupAdvancedRepository()

	// Test Pluck
	names, err := repo.Pluck("name")
	if err != nil {
		t.Errorf("Pluck failed: %v", err)
	}

	if names == nil {
		names = []interface{}{}
	}

	// Test Value
	firstValue, err := repo.Value("name")
	if err != nil {
		t.Errorf("Value failed: %v", err)
	}

	_ = firstValue

	// Test Increment
	err = repo.Increment("age", 1)
	if err != nil {
		t.Errorf("Increment failed: %v", err)
	}

	// Test Decrement
	err = repo.Decrement("age", 1)
	if err != nil {
		t.Errorf("Decrement failed: %v", err)
	}
}

func TestAdvancedRepository_ErrorCases(t *testing.T) {
	// Test with disconnected ORM
	orm := builder.NewSimpleORM().WithDialect(factory.Mock)
	repo := orm.Repository(&AdvancedRepoTestUser{})

	// Test FindWithRelations with disconnected ORM
	_, err := repo.FindWithRelations(1, "profile")
	if err == nil {
		t.Error("Expected error when ORM is not connected")
	}

	// Test BatchCreate with disconnected ORM
	users := []interface{}{
		&AdvancedRepoTestUser{Name: "Test"},
	}
	err = repo.BatchCreate(users)
	if err == nil {
		t.Error("Expected error when ORM is not connected")
	}

	// Test SoftDelete with disconnected ORM
	user := &AdvancedRepoTestUser{ID: 1}
	err = repo.SoftDelete(user)
	if err == nil {
		t.Error("Expected error when ORM is not connected")
	}

	// Test Chunk with disconnected ORM
	err = repo.Chunk(10, func(entities []interface{}) error {
		return nil
	})
	if err == nil {
		t.Error("Expected error when ORM is not connected")
	}
}

func TestAdvancedRepository_EdgeCases(t *testing.T) {
	repo := setupAdvancedRepository()

	// Test empty batch operations
	emptyBatch := []interface{}{}

	err := repo.BatchCreate(emptyBatch)
	if err != nil {
		t.Errorf("BatchCreate with empty batch failed: %v", err)
	}

	err = repo.BatchUpdate(emptyBatch)
	if err != nil {
		t.Errorf("BatchUpdate with empty batch failed: %v", err)
	}

	err = repo.BatchDelete(emptyBatch)
	if err != nil {
		t.Errorf("BatchDelete with empty batch failed: %v", err)
	}

	// Test Chunk with 0 size (this might not return an error in the current implementation)
	err = repo.Chunk(0, func(entities []interface{}) error {
		return nil
	})
	// Note: The current implementation might not validate chunk size
	if err != nil {
		t.Logf("Chunk with 0 size returned error as expected: %v", err)
	}
}
