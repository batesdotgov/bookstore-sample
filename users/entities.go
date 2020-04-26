package users

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type (
	User struct {
		ID       uint64
		Name     string `json:"name,omitempty"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	Book struct {
		Title    string  `json:"title"`
		Author   string  `json:"author"`
		Price    float64 `json:"price"`
		Quantity uint16  `json:"quantity"`
	}

	Books []Book

	Purchase struct {
		ID     uint64  `json:"id"`
		Amount float64 `json:"amount"`
		Books  Books   `json:"books"`
	}

	Purchases []Purchase
)

var UserNotFoundErr error = errors.New("User not found")

func (u User) ValidateRegisterInput() error {
	return validation.ValidateStruct(&u,
		// Name can't be empty and must contain between 8 to 100 characters
		validation.Field(&u.Name, validation.Required, validation.Length(3, 100)),
		// Email can't be empty and must contain a valid email address
		validation.Field(&u.Email, validation.Required, is.Email),
		// Password can't be empty and must contain between 8 to 100 characters
		validation.Field(&u.Password, validation.Required, validation.Length(8, 100)),
	)
}
