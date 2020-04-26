package books

import "errors"

type (
	Book struct {
		ID          uint64  `json:"id"`
		Title       string  `json:"title"`
		Author      string  `json:"author"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}

	Books []Book
)

var BookNotFoundErr error = errors.New("Book not found")
