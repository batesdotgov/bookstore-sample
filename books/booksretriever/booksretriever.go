package booksretriever

import (
	"github.com/go-chi/chi"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Invoke(registerEndpoints),
	factories,
)

var factories = fx.Provide(
	NewBooksAvailableHandler,
	NewBookDetailsHandler,
	NewBooksRetrieverRepository,
)

type Handlers struct {
	fx.In

	BooksAvailable BooksAvailableHandler
	BookDetails    BookDetailsHandler
}

func registerEndpoints(router chi.Router, handlers Handlers) {
	router.Method("GET", "/books", handlers.BooksAvailable)
	router.Method("GET", "/books/{id}", handlers.BookDetails)
}
