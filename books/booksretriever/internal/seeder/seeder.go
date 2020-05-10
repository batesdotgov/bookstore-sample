package seeder

import "database/sql"

func Seed(db *sql.DB) {
	insert(
		db,
		"The Hitchhiker's Guide to the Galaxy",
		"Douglas Adams",
		"A great book, please read it",
		19.9,
		100,
	)
	insert(
		db,
		"The Restaurant at the End of the Universe",
		"Douglas Adams",
		"A great book, please read it",
		19.9,
		100,
	)
	insert(
		db,
		"Life, the Universe and Everything",
		"Douglas Adams",
		"A great book, please read it",
		19.9,
		100,
	)
	insert(
		db,
		"So Long, and Thanks For All the Fish",
		"Douglas Adams",
		"A great book, please read it",
		19.9,
		100,
	)
}

func insert(db *sql.DB, args ...interface{}) {
	_, err := db.Exec(
		"INSERT INTO books (title, author, description, price, available) VALUES (?, ?, ?, ?, ?)",
		args...,
	)
	if err != nil {
		panic(err)
	}
}
