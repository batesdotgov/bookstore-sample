package schema

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func trim(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func runSQLFromFile(db *sql.DB, file string) {
	schema, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	statements := strings.Split(string(schema), ";")
	for _, statement := range statements {
		sql := trim(statement)
		if sql == "" {
			continue
		}
		_, err = db.Exec(sql)
		if err == nil {
			continue
		}
		panic(err)
	}
}

// Up runs the up.sql into our MySQL server
func Up(db *sql.DB) {
	dir := os.Getenv("APP_DIR")

	if dir == "" {
		dir = "/app"
	}

	file := fmt.Sprintf("%s/resources/schema/up.sql", dir)

	runSQLFromFile(db, file)
}

// Down runs the up.sql into our MySQL server
func Down(db *sql.DB) {
	dir := os.Getenv("APP_DIR")

	if dir == "" {
		dir = "/app"
	}

	file := fmt.Sprintf("%s/resources/schema/down.sql", dir)

	var err error
	// disable foreign keys
	_, err = db.Exec("SET FOREIGN_KEY_CHECKS=0")
	if err != nil {
		panic(err)
	}
	runSQLFromFile(db, file)
	// enable the foreign keys back
	_, err = db.Exec("SET FOREIGN_KEY_CHECKS=1")
	if err != nil {
		panic(err)
	}
}
