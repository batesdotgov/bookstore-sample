package purchasepersister

import (
	"context"
	"database/sql"

	"github.com/diegoholiveira/bookstore-sample/purchases"
)

type (
	DatabaseExecutor interface {
		ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	}

	persister struct {
		db DatabaseExecutor
	}
)

func (p persister) Persist(ctx context.Context, purchase purchases.Purchase) error {
	id, err := p.persistPurchase(ctx, purchase)
	if err != nil {
		return err
	}

	err = p.persistPurchasedBooks(ctx, id, purchase.Books)
	if err != nil {
		return err
	}

	err = p.decrementCopiesAvailable(ctx, purchase.Books)
	if err != nil {
		return err
	}

	return nil
}

func (p persister) persistPurchase(ctx context.Context, purchase purchases.Purchase) (int64, error) {
	r, err := p.db.ExecContext(
		ctx,
		"INSERT INTO purchases (user_id, amount) VALUES (?, ?)",
		purchase.UserID,
		purchase.Books.CalculatePurchaseAmount(),
	)
	if err != nil {
		return 0, err
	}

	return r.LastInsertId()
}

func (p persister) persistPurchasedBooks(ctx context.Context, id int64, books purchases.Books) error {
	for _, b := range books {
		_, err := p.db.ExecContext(
			ctx,
			"INSERT INTO purchased_books VALUES (?, ?, ?)",
			id,
			b.ID,
			b.Quantity,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p persister) decrementCopiesAvailable(ctx context.Context, books purchases.Books) error {
	for _, b := range books {
		_, err := p.db.ExecContext(
			ctx,
			"UPDATE books SET available = available - ? WHERE id = ?",
			b.Quantity,
			b.ID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
