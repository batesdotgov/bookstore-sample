package purchases

import "math"

type (
	Purchase struct {
		UserID uint64 `json:"user_id"`
		Books  Books  `json:"books"`
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
