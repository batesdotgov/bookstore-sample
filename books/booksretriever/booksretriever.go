package booksretriever

import (
	"database/sql"

	"github.com/go-chi/chi"
)

func SetupModule(router chi.Router, db *sql.DB) {
	var (
		booksAvailableHandler = InitBooksAvailableHandler(db)
		booksDetailsHandler   = InitBookDetailsHandler(db)
	)

	router.Method("GET", "/books", booksAvailableHandler)
	router.Method("GET", "/books/{id}", booksDetailsHandler)
}
