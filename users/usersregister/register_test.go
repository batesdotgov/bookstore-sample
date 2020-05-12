package usersregister

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/diegoholiveira/bookstore-sample/users"
	"github.com/diegoholiveira/bookstore-sample/users/usersregister/internal/mocks"
)

func TestRegisterPersistTheUserSuccessfully(t *testing.T) {
	u := users.User{
		Name:     "Diego Henrique Oliveira",
		Email:    "contact@diegoholiveira.com",
		Password: "12345678",
	}

	constraints := new(mocks.UserConstraints)
	constraints.
		On("EmailIsUsed", mock.Anything, u.Email).
		Return(false)

	persister := new(mocks.DatabasePersister)
	persister.
		On("Persist", mock.Anything, &u).
		Return(nil)

	mailer := new(mocks.WelcomeMailer)
	mailer.On("Send", mock.Anything, &u)

	service := NewUserRegisterService(constraints, persister, mailer)
	err := service.Register(context.Background(), &u)

	if assert.Nil(t, err) {
		constraints.AssertExpectations(t)
		persister.AssertExpectations(t)
		mailer.AssertExpectations(t)
	}
}
