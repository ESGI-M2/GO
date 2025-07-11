package unit

import (
	"testing"

	"project/orm"
	"project/orm/core/interfaces"
	"project/orm/dialect"
)

type RepoTestUser struct {
	Id    int    `orm:"pk,auto"`
	Name  string `orm:"column:name"`
	Email string `orm:"column:email,unique"`
}

func setupRepository() orm.Repository {
	mockDialect := &dialect.MockDialect{}
	ormInstance := orm.New(mockDialect)
	ormInstance.RegisterModel(&RepoTestUser{})

	// Connect first
	config := interfaces.ConnectionConfig{
		Host:     "localhost",
		Port:     3306,
		Database: "test",
		Username: "root",
		Password: "password",
	}
	ormInstance.Connect(config)

	return ormInstance.Repository(&RepoTestUser{})
}

func TestRepository_Find(t *testing.T) {
	repo := setupRepository()
	user := &RepoTestUser{Id: 1}
	_, err := repo.Find(user.Id)
	if err != nil {
		t.Errorf("Find failed: %v", err)
	}
	// Since we're using a mock, user might be nil, which is expected
}

func TestRepository_FindAll(t *testing.T) {
	repo := setupRepository()
	users, err := repo.FindAll()
	if err != nil {
		t.Errorf("FindAll failed: %v", err)
	}
	if users == nil {
		users = []interface{}{} // ensure not nil
	}
	if len(users) != 0 {
		t.Errorf("FindAll should return an empty slice, got %d", len(users))
	}
}

func TestRepository_FindBy(t *testing.T) {
	repo := setupRepository()
	criteria := map[string]interface{}{
		"name": "test",
	}
	users, err := repo.FindBy(criteria)
	if err != nil {
		t.Errorf("FindBy failed: %v", err)
	}
	if users == nil {
		users = []interface{}{} // ensure not nil
	}
	if len(users) != 0 {
		t.Errorf("FindBy should return an empty slice, got %d", len(users))
	}
}

func TestRepository_FindOneBy(t *testing.T) {
	repo := setupRepository()
	criteria := map[string]interface{}{
		"email": "test@example.com",
	}
	_, err := repo.FindOneBy(criteria)
	if err != nil {
		t.Errorf("FindOneBy failed: %v", err)
	}
	// Since we're using a mock, user might be nil, which is expected
}

func TestRepository_Count(t *testing.T) {
	repo := setupRepository()
	count, err := repo.Count()
	if err != nil {
		t.Errorf("Count failed: %v", err)
	}
	if count < 0 {
		t.Error("Count should be non-negative")
	}
}

func TestRepository_Exists(t *testing.T) {
	repo := setupRepository()
	user := &RepoTestUser{Id: 1}
	_, err := repo.Exists(user.Id)
	if err != nil {
		t.Errorf("Exists failed: %v", err)
	}
	// exists can be true or false depending on mock implementation
}

func TestRepository_Save_Insert(t *testing.T) {
	repo := setupRepository()
	user := &RepoTestUser{
		Id:    1,
		Name:  "Test User",
		Email: "test@example.com",
	}
	err := repo.Save(user)
	if err != nil {
		t.Errorf("Save (insert) failed: %v", err)
	}
}

func TestRepository_Save_Update(t *testing.T) {
	repo := setupRepository()
	user := &RepoTestUser{
		Id:    1,
		Name:  "Updated User",
		Email: "updated@example.com",
	}
	err := repo.Save(user)
	if err != nil {
		t.Errorf("Save (update) failed: %v", err)
	}
}

func TestRepository_Update(t *testing.T) {
	repo := setupRepository()
	user := &RepoTestUser{
		Id:    1,
		Name:  "Updated User",
		Email: "updated@example.com",
	}
	err := repo.Update(user)
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}
}

func TestRepository_Delete(t *testing.T) {
	repo := setupRepository()
	user := &RepoTestUser{
		Id: 1,
	}
	err := repo.Delete(user)
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}
}

func TestRepository_DeleteBy(t *testing.T) {
	repo := setupRepository()
	criteria := map[string]interface{}{
		"email": "test@example.com",
	}
	err := repo.DeleteBy(criteria)
	if err != nil {
		t.Errorf("DeleteBy failed: %v", err)
	}
}

func TestRepository_ErrorCases(t *testing.T) {
	// Test with nil ORM (should return error repository)
	mockDialect := &dialect.MockDialect{}
	ormInstance := orm.New(mockDialect)
	// Don't register model to trigger error case
	repo := ormInstance.Repository(&RepoTestUser{})

	_, err := repo.Find(1)
	if err == nil {
		t.Error("Expected error when metadata is not available")
	}
}
