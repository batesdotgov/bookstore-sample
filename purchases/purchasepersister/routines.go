//go:generate mockery -name=Purchaser -output=./internal/mocks

package purchasepersister

import (
	"context"
	"database/sql"

	"github.com/diegoholiveira/bookstore-sample/purchases"
	"github.com/diegoholiveira/bookstore-sample/users"
)

type (
	PurchaseRoutineWithNewUser struct {
		db *sql.DB

		registrationServiceFactory UserRegisterServiceFactory
	}

	PurchaseRoutineWithRegisteredUser struct {
		db *sql.DB
	}
)

func NewPurchaserWithNewUser(db *sql.DB, registrationServiceFactory UserRegisterServiceFactory) PurchaserWithNewUser {
	return PurchaseRoutineWithNewUser{
		db: db,

		registrationServiceFactory: registrationServiceFactory,
	}
}

func NewPurchaserWithRegisteredUser(db *sql.DB) PurchaserWithRegisteredUser {
	return PurchaseRoutineWithRegisteredUser{
		db: db,
	}
}

func (s PurchaseRoutineWithRegisteredUser) MakePurchase(ctx context.Context, p purchases.Purchase) error {
	tx, err := prepare(ctx, s.db, p)
	if err != nil {
		return err
	}
	return persist(ctx, tx, p)
}

func (s PurchaseRoutineWithNewUser) MakePurchase(ctx context.Context, p purchases.Purchase) error {
	tx, err := prepare(ctx, s.db, p)
	if err != nil {
		return err
	}
	// register the user
	err = s.register(ctx, tx, p.User)
	if err == nil {
		p.UserID = p.User.ID

		return persist(ctx, tx, p)
	}

	_ = tx.Rollback()

	return err
}

func (s PurchaseRoutineWithNewUser) register(ctx context.Context, tx *sql.Tx, u *users.User) error {
	service := s.registrationServiceFactory.Create(tx)
	return service.Register(ctx, u)
}

func prepare(ctx context.Context, db *sql.DB, p purchases.Purchase) (*sql.Tx, error) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return nil, err
	}

	constraints := constraints{
		db: tx,
	}
	err = constraints.CheckAvailability(ctx, p.Books)
	if err != nil {
		_ = tx.Rollback()

		return nil, err
	}

	return tx, nil
}

func persist(ctx context.Context, tx *sql.Tx, p purchases.Purchase) error {
	persister := persister{
		db: tx,
	}
	err := persister.Persist(ctx, p)
	if err == nil {
		return tx.Commit()
	}

	_ = tx.Rollback()

	return err
}
