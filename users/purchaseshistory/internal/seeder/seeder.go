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
	insert(
		db,
		"INSERT INTO books (title, author, description, price, available) VALUES (?, ?, ?, ?, ?)",
		"The Hitchhiker's Guide to the Galaxy",
		"Douglas Adams",
		"A great book, please read it",
		19.9,
		100,
	)
	insert(
		db,
		"INSERT INTO books (title, author, description, price, available) VALUES (?, ?, ?, ?, ?)",
		"The Restaurant at the End of the Universe",
		"Douglas Adams",
		"A great book, please read it",
		19.9,
		100,
	)
	insert(
		db,
		"INSERT INTO purchases (user_id, amount) VALUES (?, ?)",
		1,
		19.9,
	)
	insert(
		db,
		"INSERT INTO purchased_books VALUES (?, ?, ?)",
		1, // purchase_id
		1, // book_id
		1, // quantity
	)
}

func insert(db *sql.DB, sql string, args ...interface{}) {
	_, err := db.Exec(sql, args...)
	if err != nil {
		panic(err)
	}
}
