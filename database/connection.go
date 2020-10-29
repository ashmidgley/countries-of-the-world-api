package database

import "fmt"

const (
	host     = ""
	port     = 5432
	user     = ""
	password = ""
	dbname   = ""
)

// GetConnectionString returns the database connection string.
func GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}
