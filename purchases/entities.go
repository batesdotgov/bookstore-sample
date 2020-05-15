package purchases

import (
	"math"

	"github.com/diegoholiveira/bookstore-sample/users"
)

type (
	Purchase struct {
		User   *users.User `json:"user,omitempty"`
		UserID uint64      `json:"user_id"`
		Books  Books       `json:"books"`
	}

	Book struct {
		ID       uint64 `json:"id"`
		Quantity uint16 `json:"quantity"`
		Price    float64
	}

	Books []Book
)

func (books Books) CalculatePurchaseAmount() float64 {
	var sum float64
	for _, book := range books {
		sum = sum + (float64(book.Quantity) * book.Price)
	}
	return math.Round(sum*100) / 100
}

func (p Purchase) HasRegisteredUser() bool {
	return p.UserID > 0
}

func (p Purchase) HasNewUser() bool {
	var empty *users.User
	return p.User != empty && p.UserID == 0
}
