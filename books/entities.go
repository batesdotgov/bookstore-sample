package books

import "errors"

type (
	Book struct {
		ID          uint64  `json:"id"`
		Title       string  `json:"title"`
		Author      string  `json:"author"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Available   uint32  `json:"available"`
	}

	Books []Book
)

var BookNotFoundErr error = errors.New("Book not found")
