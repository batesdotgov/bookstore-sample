package purchasepersister

import (
	"context"
	"database/sql"

	"github.com/diegoholiveira/bookstore-sample/purchases"
)

type persister struct {
	db *sql.DB
}

func NewPurchasePersister(db *sql.DB) PurchasePersister {
	return &persister{
		db: db,
	}
}

func (p persister) Persist(ctx context.Context, purchase purchases.Purchase) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	id, err := p.persistPurchase(ctx, tx, purchase)
	if err != nil {
		_ = tx.Rollback()

		return err
	}

	err = p.persistPurchasedBooks(ctx, tx, id, purchase.Books)
	if err != nil {
		_ = tx.Rollback()

		return err
	}

	return tx.Commit()
}

func (p persister) persistPurchase(ctx context.Context, tx *sql.Tx, purchase purchases.Purchase) (int64, error) {
	r, err := tx.ExecContext(
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

func (p persister) persistPurchasedBooks(ctx context.Context, tx *sql.Tx, id int64, books purchases.Books) error {
	for _, b := range books {
		_, err := tx.ExecContext(
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
