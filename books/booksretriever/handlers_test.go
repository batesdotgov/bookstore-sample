package booksretriever

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/diegoholiveira/bookstore-sample/books"
	"github.com/diegoholiveira/bookstore-sample/books/booksretriever/internal/mocks"
)

func TestListingBooks(t *testing.T) {
	books := books.Books{
		{
			ID:          1,
			Title:       "The Hitchhiker's Guide to the Galaxy",
			Author:      "Douglas Adams",
			Description: "A great book, for sure. Read it",
			Price:       19.9,
		},
		{
			ID:          2,
			Title:       "The Restaurant at the End of the Universe",
			Author:      "Douglas Adams",
			Description: "A great book, for sure. Read it",
			Price:       19.9,
		},
		{
			ID:          3,
			Title:       "Life, the Universe and Everything",
			Author:      "Douglas Adams",
			Description: "A great book, for sure. Read it",
			Price:       19.9,
		},
		{
			ID:          4,
			Title:       "So Long, and Thanks For All the Fish",
			Author:      "Douglas Adams",
			Description: "A great book, for sure. Read it",
			Price:       19.9,
		},
	}

	r := httptest.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()

	finder := new(mocks.BooksFinder)
	finder.
		On("FindRecents", mock.Anything).
		Return(books, nil)

	h := NewBooksAvailableHandler(finder)
	h.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	finder.AssertExpectations(t)

	expected := `[
		{
			"id": 1,
			"title": "The Hitchhiker's Guide to the Galaxy",
			"author": "Douglas Adams",
			"description": "A great book, for sure. Read it",
			"price": 19.9
		},
		{
			"id": 2,
			"title": "The Restaurant at the End of the Universe",
			"author": "Douglas Adams",
			"description": "A great book, for sure. Read it",
			"price": 19.9
		},
		{
			"id": 3,
			"title": "Life, the Universe and Everything",
			"author": "Douglas Adams",
			"description": "A great book, for sure. Read it",
			"price": 19.9
		},
		{
			"id": 4,
			"title": "So Long, and Thanks For All the Fish",
			"author": "Douglas Adams",
			"description": "A great book, for sure. Read it",
			"price": 19.9
		}
	]`

	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		assert.JSONEq(t, expected, string(body))
	}
}

func TestGettingBookDetails_bad_request(t *testing.T) {
	r := httptest.NewRequest("GET", "/books/abc", nil)
	w := httptest.NewRecorder()

	finder := new(mocks.BooksFinder)

	h := createHandler(finder)
	h.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	finder.AssertExpectations(t)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestGettingBookDetails_not_found(t *testing.T) {
	r := httptest.NewRequest("GET", "/books/5", nil)
	w := httptest.NewRecorder()

	var bookID uint64 = 5

	finder := new(mocks.BooksFinder)
	finder.
		On("FindBookByID", mock.Anything, bookID).
		Return(books.Book{}, books.BookNotFoundErr)

	h := createHandler(finder)
	h.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	finder.AssertExpectations(t)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestGettingBookDetails_successfully(t *testing.T) {
	book := books.Book{
		ID:          1,
		Title:       "The Hitchhiker's Guide to the Galaxy",
		Author:      "Douglas Adams",
		Description: "A great book, for sure. Read it",
		Price:       19.9,
	}

	r := httptest.NewRequest("GET", "/books/1", nil)
	w := httptest.NewRecorder()

	var bookID uint64 = 1

	finder := new(mocks.BooksFinder)
	finder.
		On("FindBookByID", mock.Anything, bookID).
		Return(book, nil)

	h := createHandler(finder)
	h.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	finder.AssertExpectations(t)

	expected := `{
		"id": 1,
		"title": "The Hitchhiker's Guide to the Galaxy",
		"author": "Douglas Adams",
		"description": "A great book, for sure. Read it",
		"price": 19.9
	}`

	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		assert.JSONEq(t, expected, string(body))
	}
}

// createHandler give us a handler wrapped with the chi.Router
// to allow chi.URLParam work properly
func createHandler(f BooksFinder) chi.Router {
	r := chi.NewRouter()

	registerEndpoints(r, Handlers{
		BookDetails: NewBookDetailsHandler(f),
	})

	return r
}
