// +build integration

package usersregister

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/diegoholiveira/bookstore-sample/internal/schema"
	"github.com/diegoholiveira/bookstore-sample/users"
)

func TestUserIsPersistedSuccessfully(t *testing.T) {
	schema.Up(db)
	defer schema.Down(db)

	u := users.User{
		Name:     "Diego Henrique Oliveira",
		Email:    "contact@diegoholiveira.com",
		Password: "12345678",
	}

	r := NewDatabasePersister(db)

	if assert.Nil(t, r.Persist(context.Background(), &u)) {
		var expectedID uint64 = 1

		assertUserIsSaved(t)

		assert.Equal(t, expectedID, u.ID)
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
