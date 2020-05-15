package purchasepersister

import (
	"context"

	"github.com/diegoholiveira/bookstore-sample/purchases"
)

type (
	PurchaserWithNewUser        Purchaser
	PurchaserWithRegisteredUser Purchaser

	PurchaserService struct {
		newUser        PurchaserWithNewUser
		registeredUser PurchaserWithRegisteredUser
	}
)

func NewPurchaserService(registeredUser PurchaserWithRegisteredUser, newUser PurchaserWithNewUser) Purchaser {
	return PurchaserService{
		newUser:        newUser,
		registeredUser: registeredUser,
	}
}

func (s PurchaserService) MakePurchase(ctx context.Context, p purchases.Purchase) error {
	if p.HasNewUser() {
		return s.newUser.MakePurchase(ctx, p)
	}

	if p.HasRegisteredUser() {
		return s.registeredUser.MakePurchase(ctx, p)
	}

	return ErrPurchaseInvalid{
		Message: "A purchase must have a registered user or a new user",
	}
}
