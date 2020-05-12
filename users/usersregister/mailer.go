package usersregister

import (
	"context"

	"github.com/diegoholiveira/bookstore-sample/users"
)

type mailer struct {
}

func NewWelcomeMailer() WelcomeMailer {
	return mailer{}
}

func (m mailer) Send(ctx context.Context, u *users.User) {
	// do nothing here
}
