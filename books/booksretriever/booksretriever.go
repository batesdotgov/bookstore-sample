package booksretriever

import (
	"database/sql"

	"github.com/go-chi/chi"
)

func SetupModule(router chi.Router, db *sql.DB) {
	repo := NewBooksRetrieverRepository(db)

	booksAvailableHandler := NewBooksAvailableHandler(repo)
	booksDetailsHandler := NewBookDetailsHandler(repo)

	router.Method("GET", "/books", booksAvailableHandler)
	router.Method("GET", "/books/{id}", booksDetailsHandler)
}
