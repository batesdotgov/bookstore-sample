package purchasepersister

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/diegoholiveira/bookstore-sample/purchases"
	"github.com/diegoholiveira/bookstore-sample/purchases/purchasepersister/internal/mocks"
	"github.com/diegoholiveira/bookstore-sample/users"
)

func TestPurchasePersisterWithRegisteredUser(t *testing.T) {
	purchase := purchases.Purchase{
		UserID: 1,
		Books: purchases.Books{
			{
				ID:       1,
				Quantity: 2,
				Price:    19.9,
			},
		},
	}

	registeredUser := new(mocks.Purchaser)
	registeredUser.
		On("MakePurchase", mock.Anything, purchase).
		Return(nil)

	newUser := new(mocks.Purchaser)

	purchaser := NewPurchaserService(registeredUser, newUser)

	err := purchaser.MakePurchase(context.Background(), purchase)

	if assert.Nil(t, err) {
		registeredUser.AssertExpectations(t)
		newUser.AssertExpectations(t)
	}
}

func TestPurchasePersisterWithNewUser(t *testing.T) {
	purchase := purchases.Purchase{
		User: &users.User{
			Name:  "Diego",
			Email: "contact@diegoholiveira.com",
		},
		Books: purchases.Books{
			{
				ID:       1,
				Quantity: 2,
				Price:    19.9,
			},
		},
	}

	registeredUser := new(mocks.Purchaser)

	newUser := new(mocks.Purchaser)
	newUser.
		On("MakePurchase", mock.Anything, purchase).
		Return(nil)

	purchaser := NewPurchaserService(registeredUser, newUser)

	err := purchaser.MakePurchase(context.Background(), purchase)

	if assert.Nil(t, err) {
		registeredUser.AssertExpectations(t)
		newUser.AssertExpectations(t)
	}
}
