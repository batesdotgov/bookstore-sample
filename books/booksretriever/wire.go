//+build wireinject

package booksretriever

import (
	"database/sql"

	"github.com/google/wire"
)

func InitBooksAvailableHandler(db *sql.DB) BooksAvailableHandler {
	panic(wire.Build(
		NewBooksRetrieverRepository,
		NewBooksAvailableHandler,
	))
}

func InitBookDetailsHandler(db *sql.DB) BookDetailsHandler {
	panic(wire.Build(
		NewBooksRetrieverRepository,
		NewBookDetailsHandler,
	))
}
