// +build integration

package purchasepersister

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/diegoholiveira/bookstore-sample/internal/database"
	"github.com/diegoholiveira/bookstore-sample/internal/schema"
	"github.com/diegoholiveira/bookstore-sample/purchases"
	"github.com/diegoholiveira/bookstore-sample/purchases/purchasepersister/internal/seeder"
	"github.com/diegoholiveira/bookstore-sample/users"
)

var db *sql.DB

type FakeMailer struct{}

func (f FakeMailer) Send(ctx context.Context, u *users.User) {
}

func TestPurchaseIsPersistedSuccessfully_with_registered_user(t *testing.T) {
	schema.Up(db)
	defer schema.Down(db)

	seeder.Seed(db)

	s := NewPurchaserWithRegisteredUser(db)

	err := s.MakePurchase(context.Background(), purchases.Purchase{
		UserID: 1,
		Books: purchases.Books{
			{
				ID:       1,
				Price:    19.9,
				Quantity: 2,
			},
		},
	})

	if assert.Nil(t, err) {
		assertPurchaseIsSaved(t)
	}
}

func TestPurchaseIsPersistedSuccessfully_with_new_user(t *testing.T) {
	schema.Up(db)
	defer schema.Down(db)

	seeder.Seed(db)

	u := &users.User{
		Name:  "Diego Oliveira",
		Email: "contato@diegoholiveira.com",
	}

	registrationServiceFactory := NewUserRegisterServiceFactory(FakeMailer{})

	s := NewPurchaserWithNewUser(db, registrationServiceFactory)

	err := s.MakePurchase(context.Background(), purchases.Purchase{
		User: u,
		Books: purchases.Books{
			{
				ID:       1,
				Price:    19.9,
				Quantity: 2,
			},
		},
	})

	if assert.Nil(t, err) {
		assertUserIsSaved(t)
		assertPurchaseIsSaved(t)
	}
}

func TestPurchaseWithoutCopiesAvailable(t *testing.T) {
	schema.Up(db)
	defer schema.Down(db)

	seeder.Seed(db)

	u := &users.User{
		Name:  "Diego Oliveira",
		Email: "contato@diegoholiveira.com",
	}

	registrationServiceFactory := NewUserRegisterServiceFactory(FakeMailer{})

	s := NewPurchaserWithNewUser(db, registrationServiceFactory)

	err := s.MakePurchase(context.Background(), purchases.Purchase{
		User: u,
		Books: purchases.Books{
			{
				ID:       2,
				Price:    19.9,
				Quantity: 3,
			},
		},
	})

	assert.EqualError(t, err, "The book 2 does not have 3 copies")
}

func assertPurchaseIsSaved(t *testing.T) {
	var (
		rows     = 0
		expected = 1

		amount          = 0.0
		amount_expected = 39.8

		available          = 0
		available_expected = 98

		err error
	)

	err = db.QueryRow("SELECT count(id) as total FROM purchases").Scan(&rows)
	assert.Nil(t, err)
	assert.Equal(t, expected, rows)

	err = db.QueryRow("SELECT count(book_id) as total FROM purchased_books").Scan(&rows)
	assert.Nil(t, err)
	assert.Equal(t, expected, rows)

	err = db.QueryRow("SELECT amount FROM purchases WHERE id = ? LIMIT 1", 1).Scan(&amount)
	assert.Nil(t, err)
	assert.Equal(t, amount_expected, amount)

	err = db.QueryRow("SELECT available FROM books WHERE id = ? LIMIT 1", 1).Scan(&available)
	assert.Nil(t, err)
	assert.Equal(t, available_expected, available)
}

func assertUserIsSaved(t *testing.T) {
	var (
		rows     = 0
		expected = 2
	)

	err := db.QueryRow("SELECT count(id) as total FROM users").Scan(&rows)
	assert.Nil(t, err)
	assert.Equal(t, expected, rows)
}

func init() {
	var err error
	db, err = database.NewMySQLConnection()
	if err != nil {
		panic(err)
	}
}
