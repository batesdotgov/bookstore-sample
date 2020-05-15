package purchasepersister

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/diegoholiveira/bookstore-sample/purchases"
)

type (
	DatabaseQuerier interface {
		QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	}

	constraints struct {
		db DatabaseQuerier
	}
)

func (c constraints) CheckAvailability(ctx context.Context, books purchases.Books) error {
	for _, book := range books {
		var available uint16

		err := c.db.QueryRowContext(
			ctx,
			"SELECT available FROM books WHERE id = ? FOR UPDATE",
			book.ID,
		).Scan(&available)
		if err == sql.ErrNoRows {
			return ErrPurchaseInvalid{
				Message: fmt.Sprintf("The book %d is not available", book.ID),
			}
		}
		if err != nil {
			return err
		}
		if book.Quantity > available {
			return ErrPurchaseInvalid{
				Message: fmt.Sprintf("The book %d does not have %d copies", book.ID, book.Quantity),
			}
		}
	}
	return nil
}
