module github.com/ESGI-M2/GO

go 1.24.1

// Retract broken versions that don't have the FullTextSearch method
retract (
	v0.1.3-dev
	v0.0.3-dev
	v0.0.2-dev
	v0.0.1-dev
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-sql-driver/mysql v1.9.2 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
