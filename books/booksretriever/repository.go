package booksretriever

import (
	"context"
	"database/sql"

	"github.com/diegoholiveira/bookstore-sample/books"
)

type (
	repository struct {
		db *sql.DB
	}
)

func NewBooksRetrieverRepository(db *sql.DB) BooksFinder {
	return repository{
		db: db,
	}
}

func (r repository) FindBookByID(ctx context.Context, id uint64) (books.Book, error) {
	book := books.Book{}

	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, title, author, description, price FROM books WHERE id = ? LIMIT 1",
		id,
	).Scan(&book.ID, &book.Title, &book.Author, &book.Description, &book.Price)

	switch {
	case err == sql.ErrNoRows:
		return book, books.BookNotFoundErr
	default:
		return book, err
	}
}

func (r repository) FindRecents(ctx context.Context) (books.Books, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT id, title, author, description, price FROM books ORDER BY id DESC",
	)

	if err != nil {
		return books.Books{}, err
	}

	_books := make(books.Books, 0)
	for rows.Next() {
		book := books.Book{}

		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Description,
			&book.Price,
		)
		if err != nil {
			return books.Books{}, err
		}

		_books = append(_books, book)
	}

	return _books, nil
}
