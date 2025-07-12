package integration

import (
	"testing"

	"github.com/ESGI-M2/GO/orm/core/interfaces"
)

// TestAdvancedQueryFeatures tests advanced query building features
func TestAdvancedQueryFeatures(t *testing.T) {
	resetMockUsers()
	// Setup
	db := setupTestDB(t)
	defer db.Close()

	// Create test tables
	createTestTables(t, db)

	// Insert test data
	insertTestData(t, db)

	t.Run("OR Conditions", func(t *testing.T) {
		query := db.Query(&User{})
		results, err := query.WhereOr(
			interfaces.WhereCondition{Field: "name", Operator: "=", Value: "John"},
			interfaces.WhereCondition{Field: "age", Operator: ">", Value: 25},
		).Find()

		if err != nil {
			t.Fatalf("Failed to execute OR query: %v", err)
		}

		if len(results) == 0 {
			t.Error("Expected results from OR query")
		}

		t.Logf("OR query returned %d results", len(results))
	})

	t.Run("Raw WHERE Conditions", func(t *testing.T) {
		query := db.Query(&User{})
		results, err := query.WhereRaw("name LIKE ? AND age > ?", "%John%", 20).Find()

		if err != nil {
			t.Fatalf("Failed to execute raw WHERE query: %v", err)
		}

		if len(results) == 0 {
			t.Error("Expected results from raw WHERE query")
		}

		t.Logf("Raw WHERE query returned %d results", len(results))
	})

	t.Run("BETWEEN Conditions", func(t *testing.T) {
		query := db.Query(&User{})
		results, err := query.WhereBetween("age", 20, 30).Find()

		if err != nil {
			t.Fatalf("Failed to execute BETWEEN query: %v", err)
		}

		if len(results) == 0 {
			t.Error("Expected results from BETWEEN query")
		}

		t.Logf("BETWEEN query returned %d results", len(results))
	})

	t.Run("NULL Conditions", func(t *testing.T) {
		query := db.Query(&User{})
		results, err := query.WhereNotNull("email").Find()

		if err != nil {
			t.Fatalf("Failed to execute NOT NULL query: %v", err)
		}

		t.Logf("NOT NULL query returned %d results", len(results))
	})

	t.Run("LIKE Conditions", func(t *testing.T) {
		query := db.Query(&User{})
		results, err := query.WhereLike("name", "%John%").Find()

		if err != nil {
			t.Fatalf("Failed to execute LIKE query: %v", err)
		}

		if len(results) == 0 {
			t.Error("Expected results from LIKE query")
		}

		t.Logf("LIKE query returned %d results", len(results))
	})

	t.Run("DISTINCT Query", func(t *testing.T) {
		query := db.Query(&User{})
		results, err := query.Distinct().Select("age").Find()

		if err != nil {
			t.Fatalf("Failed to execute DISTINCT query: %v", err)
		}

		t.Logf("DISTINCT query returned %d results", len(results))
	})

	t.Run("FOR UPDATE Lock", func(t *testing.T) {
		query := db.Query(&User{})
		results, err := query.ForUpdate().Where("id", "=", 1).Find()

		if err != nil {
			t.Fatalf("Failed to execute FOR UPDATE query: %v", err)
		}

		if len(results) == 0 {
			t.Error("Expected results from FOR UPDATE query")
		}

		t.Logf("FOR UPDATE query returned %d results", len(results))
	})
}

// TestPaginationFeatures tests pagination features
func TestPaginationFeatures(t *testing.T) {
	resetMockUsers()
	// Setup
	db := setupTestDB(t)
	defer db.Close()

	// Create test tables
	createTestTables(t, db)

	// Insert test data
	insertTestData(t, db)

	t.Run("Offset Pagination", func(t *testing.T) {
		query := db.Query(&User{})
		results, err := query.OffsetPaginate(1, 5).Find()

		if err != nil {
			t.Fatalf("Failed to execute offset pagination: %v", err)
		}

		if len(results) > 5 {
			t.Error("Expected maximum 5 results from pagination")
		}

		t.Logf("Offset pagination returned %d results", len(results))
	})

	t.Run("Cursor Pagination", func(t *testing.T) {
		query := db.Query(&User{})
		results, err := query.CursorPaginate("id", nil, 3).Find()

		if err != nil {
			t.Fatalf("Failed to execute cursor pagination: %v", err)
		}

		if len(results) > 3 {
			t.Error("Expected maximum 3 results from cursor pagination")
		}

		t.Logf("Cursor pagination returned %d results", len(results))
	})
}

// TestCachingFeatures tests query caching features
func TestCachingFeatures(t *testing.T) {
	resetMockUsers()
	// Setup
	db := setupTestDB(t)
	defer db.Close()

	// Create test tables
	createTestTables(t, db)

	// Insert test data
	insertTestData(t, db)

	t.Run("Cache Enabled", func(t *testing.T) {
		query := db.Query(&User{})
		results1, err := query.Cache(300).Where("id", "=", 1).Find()

		if err != nil {
			t.Fatalf("Failed to execute cached query: %v", err)
		}

		// Execute same query again (should use cache)
		results2, err := query.Cache(300).Where("id", "=", 1).Find()

		if err != nil {
			t.Fatalf("Failed to execute cached query (second time): %v", err)
		}

		if len(results1) != len(results2) {
			t.Error("Cached results should be identical")
		}

		t.Logf("Cache test completed successfully")
	})

	t.Run("Cache Disabled", func(t *testing.T) {
		query := db.Query(&User{})
		results, err := query.WithoutCache().Where("id", "=", 1).Find()

		if err != nil {
			t.Fatalf("Failed to execute non-cached query: %v", err)
		}

		if len(results) == 0 {
			t.Error("Expected results from non-cached query")
		}

		t.Logf("Non-cached query returned %d results", len(results))
	})
}

// TestBatchOperations tests batch operations
func TestBatchOperations(t *testing.T) {
	resetMockUsers()
	// Setup
	db := setupTestDB(t)
	defer db.Close()

	// Create test tables
	createTestTables(t, db)

	t.Run("Batch Create", func(t *testing.T) {
		users := []interface{}{
			&User{Name: "Batch User 1", Age: 25, Email: "batch1@test.com"},
			&User{Name: "Batch User 2", Age: 30, Email: "batch2@test.com"},
			&User{Name: "Batch User 3", Age: 35, Email: "batch3@test.com"},
		}

		repo := db.Repository(&User{})
		err := repo.BatchCreate(users)

		if err != nil {
			t.Fatalf("Failed to batch create users: %v", err)
		}

		// Verify creation
		count, err := repo.Count()
		if err != nil {
			t.Fatalf("Failed to count users: %v", err)
		}

		if count < 3 {
			t.Error("Expected at least 3 users after batch create")
		}

		t.Logf("Batch create completed successfully, total users: %d", count)
	})

	t.Run("Batch Update", func(t *testing.T) {
		// First create some users
		users := []interface{}{
			&User{Name: "Update User 1", Age: 25, Email: "update1@test.com"},
			&User{Name: "Update User 2", Age: 30, Email: "update2@test.com"},
		}

		repo := db.Repository(&User{})
		err := repo.BatchCreate(users)
		if err != nil {
			t.Fatalf("Failed to create users for batch update: %v", err)
		}

		// Update users
		updatedUsers := []interface{}{
			&User{ID: 1, Name: "Updated User 1", Age: 26, Email: "updated1@test.com"},
			&User{ID: 2, Name: "Updated User 2", Age: 31, Email: "updated2@test.com"},
		}

		err = repo.BatchUpdate(updatedUsers)
		if err != nil {
			t.Fatalf("Failed to batch update users: %v", err)
		}

		t.Logf("Batch update completed successfully")
	})

	t.Run("Batch Delete", func(t *testing.T) {
		// First create some users
		users := []interface{}{
			&User{Name: "Delete User 1", Age: 25, Email: "delete1@test.com"},
			&User{Name: "Delete User 2", Age: 30, Email: "delete2@test.com"},
		}

		repo := db.Repository(&User{})
		err := repo.BatchCreate(users)
		if err != nil {
			t.Fatalf("Failed to create users for batch delete: %v", err)
		}

		// Delete users
		err = repo.BatchDelete(users)
		if err != nil {
			t.Fatalf("Failed to batch delete users: %v", err)
		}

		t.Logf("Batch delete completed successfully")
	})
}

// TestSoftDeleteFeatures tests soft delete features
func TestSoftDeleteFeatures(t *testing.T) {
	resetMockUsers()
	// Setup
	db := setupTestDB(t)
	defer db.Close()

	// Create test tables with soft deletes
	createTestTablesWithSoftDeletes(t, db)

	t.Run("Soft Delete", func(t *testing.T) {
		user := &UserWithSoftDelete{
			Name:  "Soft Delete User",
			Age:   25,
			Email: "softdelete@test.com",
		}

		repo := db.Repository(&UserWithSoftDelete{})
		err := repo.Save(user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		// Soft delete
		err = repo.SoftDelete(user)
		if err != nil {
			t.Fatalf("Failed to soft delete user: %v", err)
		}

		// Verify user is not in normal queries
		results, err := repo.FindAll()
		if err != nil {
			t.Fatalf("Failed to find all users: %v", err)
		}

		if len(results) > 0 {
			t.Error("Soft deleted user should not appear in normal queries")
		}

		// Verify user is in trashed queries
		trashed, err := repo.FindTrashed()
		if err != nil {
			t.Fatalf("Failed to find trashed users: %v", err)
		}

		if len(trashed) == 0 {
			t.Error("Soft deleted user should appear in trashed queries")
		}

		t.Logf("Soft delete test completed successfully")
	})

	t.Run("Restore", func(t *testing.T) {
		user := &UserWithSoftDelete{
			Name:  "Restore User",
			Age:   30,
			Email: "restore@test.com",
		}

		repo := db.Repository(&UserWithSoftDelete{})
		err := repo.Save(user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		// Soft delete
		err = repo.SoftDelete(user)
		if err != nil {
			t.Fatalf("Failed to soft delete user: %v", err)
		}

		// Restore
		err = repo.Restore(user)
		if err != nil {
			t.Fatalf("Failed to restore user: %v", err)
		}

		// Verify user is back in normal queries
		results, err := repo.FindAll()
		if err != nil {
			t.Fatalf("Failed to find all users: %v", err)
		}

		if len(results) == 0 {
			t.Error("Restored user should appear in normal queries")
		}

		t.Logf("Restore test completed successfully")
	})
}

// TestEagerLoadingFeatures tests eager loading features
func TestEagerLoadingFeatures(t *testing.T) {
	resetMockUsers()
	// Setup
	db := setupTestDB(t)
	defer db.Close()

	// Create test tables with relations
	createTestTablesWithRelations(t, db)

	// Insert test data with relations
	insertTestDataWithRelations(t, db)

	t.Run("Eager Loading", func(t *testing.T) {
		query := db.Query(&User{})
		results, err := query.With("posts", func(q interfaces.QueryBuilder) interfaces.QueryBuilder {
			return q
		}).Find()

		if err != nil {
			t.Fatalf("Failed to execute eager loading query: %v", err)
		}

		if len(results) == 0 {
			t.Error("Expected results from eager loading query")
		}

		// Verify relations are loaded
		for _, result := range results {
			if posts, exists := result["posts"]; exists {
				t.Logf("User has %d posts", len(posts.([]map[string]interface{})))
			}
		}

		t.Logf("Eager loading test completed successfully")
	})

	t.Run("With Count", func(t *testing.T) {
		query := db.Query(&User{})
		results, err := query.WithCount("posts").Find()

		if err != nil {
			t.Fatalf("Failed to execute with count query: %v", err)
		}

		if len(results) == 0 {
			t.Error("Expected results from with count query")
		}

		t.Logf("With count test completed successfully")
	})
}

// TestChunkingFeatures tests chunking features
func TestChunkingFeatures(t *testing.T) {
	resetMockUsers()
	// Setup
	db := setupTestDB(t)
	defer db.Close()

	// Create test tables
	createTestTables(t, db)

	// Insert test data
	insertTestData(t, db)

	t.Run("Chunk Processing", func(t *testing.T) {
		repo := db.Repository(&User{})
		chunkCount := 0
		totalProcessed := 0

		err := repo.Chunk(2, func(chunk []interface{}) error {
			chunkCount++
			totalProcessed += len(chunk)
			t.Logf("Processing chunk %d with %d items", chunkCount, len(chunk))
			return nil
		})

		if err != nil {
			t.Fatalf("Failed to process chunks: %v", err)
		}

		t.Logf("Chunk processing completed: %d chunks, %d total items", chunkCount, totalProcessed)
	})

	t.Run("Each Processing", func(t *testing.T) {
		repo := db.Repository(&User{})
		processedCount := 0

		err := repo.Each(func(item interface{}) error {
			processedCount++
			user := item.(map[string]interface{})
			t.Logf("Processing user: %s", user["name"])
			return nil
		})

		if err != nil {
			t.Fatalf("Failed to process each item: %v", err)
		}

		t.Logf("Each processing completed: %d items", processedCount)
	})
}

// TestIncrementDecrementFeatures tests increment/decrement features
func TestIncrementDecrementFeatures(t *testing.T) {
	resetMockUsers()
	// Setup
	db := setupTestDB(t)
	defer db.Close()

	// Create test tables
	createTestTables(t, db)

	// Insert test data
	insertTestData(t, db)

	t.Run("Increment", func(t *testing.T) {
		repo := db.Repository(&User{})

		// Get initial value
		initialValue, err := repo.Value("age")
		if err != nil {
			t.Fatalf("Failed to get initial age: %v", err)
		}

		// Increment
		err = repo.Increment("age", 5)
		if err != nil {
			t.Fatalf("Failed to increment age: %v", err)
		}

		// Get new value
		newValue, err := repo.Value("age")
		if err != nil {
			t.Fatalf("Failed to get new age: %v", err)
		}

		if newValue.(int) <= initialValue.(int) {
			t.Error("Age should have been incremented")
		}

		t.Logf("Increment test completed: %d -> %d", initialValue, newValue)
	})

	t.Run("Decrement", func(t *testing.T) {
		repo := db.Repository(&User{})

		// Get initial value
		initialValue, err := repo.Value("age")
		if err != nil {
			t.Fatalf("Failed to get initial age: %v", err)
		}

		// Decrement
		err = repo.Decrement("age", 3)
		if err != nil {
			t.Fatalf("Failed to decrement age: %v", err)
		}

		// Get new value
		newValue, err := repo.Value("age")
		if err != nil {
			t.Fatalf("Failed to get new age: %v", err)
		}

		if newValue.(int) >= initialValue.(int) {
			t.Error("Age should have been decremented")
		}

		t.Logf("Decrement test completed: %d -> %d", initialValue, newValue)
	})
}

// Helper functions

func createTestTables(t *testing.T, db interfaces.ORM) {
	// Create users table
	err := db.CreateTable(&User{})
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}
}

func createTestTablesWithSoftDeletes(t *testing.T, db interfaces.ORM) {
	// Create users table with soft deletes
	err := db.CreateTable(&UserWithSoftDelete{})
	if err != nil {
		t.Fatalf("Failed to create users table with soft deletes: %v", err)
	}
}

func createTestTablesWithRelations(t *testing.T, db interfaces.ORM) {
	// Create users table
	err := db.CreateTable(&User{})
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}

	// Create posts table
	err = db.CreateTable(&Post{})
	if err != nil {
		t.Fatalf("Failed to create posts table: %v", err)
	}
}

func insertTestData(t *testing.T, db interfaces.ORM) {
	users := []*User{
		{Name: "John Doe", Age: 25, Email: "john@test.com"},
		{Name: "Jane Smith", Age: 30, Email: "jane@test.com"},
		{Name: "Bob Johnson", Age: 35, Email: "bob@test.com"},
		{Name: "Alice Brown", Age: 28, Email: "alice@test.com"},
		{Name: "Charlie Wilson", Age: 32, Email: "charlie@test.com"},
	}

	repo := db.Repository(&User{})
	for _, user := range users {
		err := repo.Save(user)
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}
}

func insertTestDataWithRelations(t *testing.T, db interfaces.ORM) {
	// Create users
	users := []*User{
		{Name: "John Doe", Age: 25, Email: "john@test.com"},
		{Name: "Jane Smith", Age: 30, Email: "jane@test.com"},
	}

	userRepo := db.Repository(&User{})
	for _, user := range users {
		err := userRepo.Save(user)
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}

	// Create posts
	posts := []*Post{
		{Title: "First Post", Content: "Content 1", UserID: 1},
		{Title: "Second Post", Content: "Content 2", UserID: 1},
		{Title: "Third Post", Content: "Content 3", UserID: 2},
	}

	postRepo := db.Repository(&Post{})
	for _, post := range posts {
		err := postRepo.Save(post)
		if err != nil {
			t.Fatalf("Failed to create test post: %v", err)
		}
	}
}
