package purchases

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBooksCalculateAmountOfPurchase(t *testing.T) {
	books := Books{
		{
			Quantity: 1,
			Price:    19.9,
		},
		{
			Quantity: 2,
			Price:    19.9,
		},
	}

	expected := 59.7

	assert.Equal(t, expected, books.CalculatePurchaseAmount())
}
