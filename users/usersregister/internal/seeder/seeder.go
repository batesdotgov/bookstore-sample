package seeder

import "database/sql"

func Seed(db *sql.DB) {
	insert(
		db,
		"INSERT INTO users (name, email, password) VALUES (?, ?, ?)",
		"Diego Henrique Oliveira",
		"contact@diegoholiveira.com",
		"12345678",
	)
}

func insert(db *sql.DB, sql string, args ...interface{}) {
	_, err := db.Exec(sql, args...)
	if err != nil {
		panic(err)
	}
}
