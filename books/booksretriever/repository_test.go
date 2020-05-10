// +build integration

package booksretriever

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/diegoholiveira/bookstore-sample/books"
	"github.com/diegoholiveira/bookstore-sample/books/booksretriever/internal/seeder"
	"github.com/diegoholiveira/bookstore-sample/internal/database"
	"github.com/diegoholiveira/bookstore-sample/internal/schema"
)

var db *sql.DB

func TestBookFindByID_not_found(t *testing.T) {
	schema.Up(db)
	defer schema.Down(db)

	r := NewBooksRetrieverRepository(db)

	_, err := r.FindBookByID(context.Background(), 1)

	assert.True(t, err == books.BookNotFoundErr)
}

func TestBookFindByID_success(t *testing.T) {
	schema.Up(db)
	defer schema.Down(db)

	seeder.Seed(db)

	r := NewBooksRetrieverRepository(db)

	book, err := r.FindBookByID(context.Background(), 1)

	expected := books.Book{
		ID:          1,
		Title:       "The Hitchhiker's Guide to the Galaxy",
		Author:      "Douglas Adams",
		Description: "A great book, please read it",
		Price:       19.9,
		Available:   100,
	}

	assert.Nil(t, err)
	assert.Equal(t, expected, book)
}

func TestBookFindRecents(t *testing.T) {
	schema.Up(db)
	defer schema.Down(db)

	seeder.Seed(db)

	r := NewBooksRetrieverRepository(db)

	expected := books.Books{
		{
			ID:          4,
			Title:       "So Long, and Thanks For All the Fish",
			Author:      "Douglas Adams",
			Description: "A great book, please read it",
			Price:       19.9,
			Available:   100,
		},
		{
			ID:          3,
			Title:       "Life, the Universe and Everything",
			Author:      "Douglas Adams",
			Description: "A great book, please read it",
			Price:       19.9,
			Available:   100,
		},
		{
			ID:          2,
			Title:       "The Restaurant at the End of the Universe",
			Author:      "Douglas Adams",
			Description: "A great book, please read it",
			Price:       19.9,
			Available:   100,
		},
		{
			ID:          1,
			Title:       "The Hitchhiker's Guide to the Galaxy",
			Author:      "Douglas Adams",
			Description: "A great book, please read it",
			Price:       19.9,
			Available:   100,
		},
	}

	books, err := r.FindRecents(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, expected, books)
}

func init() {
	var err error
	db, err = database.NewMySQLConnection()
	if err != nil {
		panic(err)
	}
}
