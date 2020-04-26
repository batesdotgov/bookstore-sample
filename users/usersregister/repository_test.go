// +build integration

package usersregister

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/diegoholiveira/bookstore-sample/internal/database"
	"github.com/diegoholiveira/bookstore-sample/internal/schema"
	"github.com/diegoholiveira/bookstore-sample/users"
)

var db *sql.DB

func TestUserIsPersistedSuccessfully(t *testing.T) {
	schema.Up(db)
	defer schema.Down(db)

	r := NewUserPersister(db)
	err := r.Persist(context.Background(), users.User{
		Name:     "Diego Henrique Oliveira",
		Email:    "contact@diegoholiveira.com",
		Password: "12345678",
	})

	if assert.Nil(t, err) {
		assertUserIsSaved(t)
	}
}

func assertUserIsSaved(t *testing.T) {
	var (
		users    = 0
		expected = 1
	)

	err := db.QueryRow("SELECT count(id) as total FROM users").Scan(&users)
	assert.Nil(t, err)
	assert.Equal(t, expected, users)
}

func init() {
	var err error
	db, err = database.NewMySQLConnection()
	if err != nil {
		panic(err)
	}
}
