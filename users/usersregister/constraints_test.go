// +build integration

package usersregister

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/diegoholiveira/bookstore-sample/internal/database"
	"github.com/diegoholiveira/bookstore-sample/internal/schema"
	"github.com/diegoholiveira/bookstore-sample/users/usersregister/internal/seeder"
)

var db *sql.DB

func TestConstraints_email_is_used(t *testing.T) {
	schema.Up(db)
	defer schema.Down(db)
	seeder.Seed(db)

	scenarios := []struct {
		Email    string
		Expected bool
	}{
		{
			Email:    "contato@diegoholiveira.com",
			Expected: false,
		},
		{
			Email:    "contact@diegoholiveira.com",
			Expected: true,
		},
	}

	constraints := NewUserConstraints(db)

	for _, scenario := range scenarios {
		assert.Equal(
			t,
			scenario.Expected,
			constraints.EmailIsUsed(context.Background(), scenario.Email),
		)
	}
}

func init() {
	var err error
	db, err = database.NewMySQLConnection()
	if err != nil {
		panic(err)
	}
}
