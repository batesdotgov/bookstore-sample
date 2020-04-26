// +build integration

package purchaseshistory

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/diegoholiveira/bookstore-sample/internal/database"
	"github.com/diegoholiveira/bookstore-sample/internal/schema"
	"github.com/diegoholiveira/bookstore-sample/users"
	"github.com/diegoholiveira/bookstore-sample/users/purchaseshistory/internal/seeder"
)

var db *sql.DB

func TestPurchasesFinder(t *testing.T) {
	schema.Up(db)
	defer schema.Down(db)
	seeder.Seed(db)

	expected := users.Purchases{
		users.Purchase{
			ID:     1,
			Amount: 19.9,
			Books: users.Books{
				users.Book{
					Title:    "The Hitchhiker's Guide to the Galaxy",
					Author:   "Douglas Adams",
					Price:    19.9,
					Quantity: 1,
				},
			},
		},
	}

	user := users.User{
		ID: 1,
	}

	repository := NewPurchasesRepository(db)

	purchases, err := repository.FindPurchasesByUser(context.Background(), user)

	assert.Nil(t, err)
	assert.Equal(t, expected, purchases)
}

func init() {
	var err error
	db, err = database.NewMySQLConnection()
	if err != nil {
		panic(err)
	}
}
