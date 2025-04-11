package orm

import (
	"testing"
	"project/dialect"
	"os"
	"strings"
	"regexp"
	"fmt"
)

func TestMain(m *testing.M) {
	dialect.InitMySQL()

	code := m.Run()

	if dialect.DB != nil {
		dialect.DB.Close()
	}

	os.Exit(code)
}

func TestQuery(t *testing.T) {
	selectQuery := Select("users", []string{"name", "age"}, "age > 30", 10)
    fmt.Println("SELECT Query:", selectQuery)
}

func TestSelect(t *testing.T) {
	tests := []struct {
		table  string
		fields []string
		where  string
		limit  int
		expected string
	}{
		{
			table:  "users",
			fields: []string{"name", "age"},
			where:  "age > 30",
			limit:  10,
			expected: "SELECT name, age FROM users WHERE age > 30 LIMIT 10",
		},
		{
			table:  "users",
			fields: []string{},
			where:  "",
			limit:  0,
			expected: "SELECT * FROM users",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("SELECT %s", tt.table), func(t *testing.T) {
			query := Select(tt.table, tt.fields, tt.where, tt.limit)
			if query != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, query)
			}
		})
	}
	
}


func TestInsert(t *testing.T) {
	type User struct {
		ID   int    `db:"id" autoincrement:"true"`
		Name string `db:"name"`
		Age  int    `db:"age"`
	}	

	tests := []struct {
		model    User
		expected string
	}{
		{
			model:    User{Name: "Jean", Age: 30},
			expected: "INSERT INTO users (name, age) VALUES (?, ?)",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("INSERT User %s", tt.model.Name), func(t *testing.T) {
			query, _ := Insert("users", tt.model)
			if query != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, query)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type User struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
		Age  int    `db:"age"`
	}

	tests := []struct {
		model    User
		where    string
		expected string
	}{
		{
			model:    User{Name: "Jean", Age: 31},
			where:    "id = 1",
			expected: "UPDATE users SET name = ?, age = ? WHERE id = 1",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("UPDATE User %d", tt.model.ID), func(t *testing.T) {
			query, _ := Update("users", tt.model, tt.where)
			if query != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, query)
			}
		})
	}
}

func TestCreateTableSQL(t *testing.T) {
	type User struct {
		ID   int    `db:"id" primary:"true" autoincrement:"true"`
		Name string `db:"name"`
		Age  int    `db:"age"`
	}
	
	type Post struct {
		ID     int    `db:"id" primary:"true" autoincrement:"true"`
		Title  string `db:"title"`
		UserID int    `db:"user_id" foreign:"users(id)"`
	}
	sql := CreateTableSQL("users", User{})
	expected := `CREATE TABLE users (
		id INT PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(255),
		age INT
		);`

	if oneLine(sql) != oneLine(expected) {
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s", expected, sql)
	}

	sqlPost := CreateTableSQL("posts", Post{})
	expectedPost := `CREATE TABLE posts (
		id INT PRIMARY KEY AUTO_INCREMENT,
		title VARCHAR(255),
		user_id INT,
		FOREIGN KEY (user_id) REFERENCES users(id)
		);`

	if oneLine(sqlPost) != oneLine(expectedPost) {
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s", expectedPost, sqlPost)
	}
}

func oneLine(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}
