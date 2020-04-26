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
)

var db *sql.DB

func TestPurchaseIsPersistedSuccessfully(t *testing.T) {
	schema.Up(db)
	defer schema.Down(db)

	seeder.Seed(db)

	r := NewPurchasePersister(db)
	err := r.Persist(context.Background(), purchases.Purchase{
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

func assertPurchaseIsSaved(t *testing.T) {
	var (
		rows     = 0
		expected = 1

		amount          = 0.0
		amount_expected = 39.8

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
}

func init() {
	var err error
	db, err = database.NewMySQLConnection()
	if err != nil {
		panic(err)
	}
}
