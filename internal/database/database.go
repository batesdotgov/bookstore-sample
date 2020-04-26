package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQLConnection() (*sql.DB, error) {
	host := os.Getenv("MYSQL_HOST")
	database := os.Getenv("MYSQL_DATABASE")
	username := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")

	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s)/%s",
		username,
		password,
		host,
		database,
	))
	return db, err
}
