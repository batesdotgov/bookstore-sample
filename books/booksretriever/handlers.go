//go:generate mockery -name=BooksFinder -output=./internal/mocks
package booksretriever

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/diegoholiveira/bookstore-sample/books"
	"github.com/diegoholiveira/bookstore-sample/pkg/http/render"
)

type (
	BooksFinder interface {
		FindRecents(context.Context) (books.Books, error)
		FindBookByID(context.Context, uint64) (books.Book, error)
	}

	BooksAvailableHandler struct {
		finder BooksFinder
	}

	BookDetailsHandler struct {
		finder BooksFinder
	}
)

func NewBooksAvailableHandler(f BooksFinder) BooksAvailableHandler {
	return BooksAvailableHandler{
		finder: f,
	}
}

func NewBookDetailsHandler(f BooksFinder) BookDetailsHandler {
	return BookDetailsHandler{
		finder: f,
	}
}

func (h BooksAvailableHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	books, _ := h.finder.FindRecents(r.Context())
	render.JSON(w, http.StatusOK, books)
}

func (h BookDetailsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bookID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid book id",
		})
		return // stop here
	}

	book, err := h.finder.FindBookByID(r.Context(), bookID)
	if err == books.BookNotFoundErr {
		render.JSON(w, http.StatusNotFound, map[string]string{
			"error": "Book not found",
		})
		return // stop here
	}

	render.JSON(w, http.StatusOK, book)
}
