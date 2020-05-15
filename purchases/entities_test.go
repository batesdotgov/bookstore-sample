package purchases

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/diegoholiveira/bookstore-sample/users"
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

func TestPurchaseHasRegisteredUser(t *testing.T) {
	scenarios := []struct {
		purchase Purchase
		expected bool
	}{
		{
			purchase: Purchase{UserID: 1},
			expected: true,
		},
		{
			purchase: Purchase{},
			expected: false,
		},
	}

	for _, scenario := range scenarios {
		assert.Equal(t, scenario.expected, scenario.purchase.HasRegisteredUser())
	}
}

func TestPurchaseHasNewUser(t *testing.T) {
	scenarios := []struct {
		purchase Purchase
		expected bool
	}{
		{
			purchase: Purchase{
				UserID: 1,
				User: &users.User{
					ID: 2,
				},
			},
			expected: false,
		},
		{
			purchase: Purchase{},
			expected: false,
		},
		{
			purchase: Purchase{
				User: &users.User{
					ID: 1,
				},
			},
			expected: true,
		},
	}

	for _, scenario := range scenarios {
		assert.Equal(t, scenario.expected, scenario.purchase.HasNewUser())
	}
}
