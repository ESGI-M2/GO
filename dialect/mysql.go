package dialect

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitMySQL() {
	var err error
	godotenv.Load("../.env")

	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	db := os.Getenv("MYSQL_DATABASE")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, host, db)

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	if err := DB.Ping(); err != nil {
		log.Fatalf("DB unreachable: %v", err)
	}
	log.Println("Connected to MySQL")
}
