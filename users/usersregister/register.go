//go:generate mockery -name=UserConstraints -output=./internal/mocks
//go:generate mockery -name=DatabasePersister -output=./internal/mocks
//go:generate mockery -name=WelcomeMailer -output=./internal/mocks
package usersregister

import (
	"context"

	"github.com/diegoholiveira/bookstore-sample/users"
)

type (
	UserConstraints interface {
		EmailIsUsed(context.Context, string) bool
	}

	DatabasePersister interface {
		Persist(context.Context, *users.User) error
	}

	WelcomeMailer interface {
		Send(context.Context, *users.User)
	}

	UserRegisterService struct {
		constraints UserConstraints
		db          DatabasePersister
		mailer      WelcomeMailer
	}
)

func NewUserRegisterService(
	constraints UserConstraints,
	db DatabasePersister,
	mailer WelcomeMailer,
) UserRegister {
	return UserRegisterService{
		constraints: constraints,
		db:          db,
		mailer:      mailer,
	}
}

func (s UserRegisterService) Register(ctx context.Context, u *users.User) error {
	if s.constraints.EmailIsUsed(ctx, u.Email) {
		return users.ErrEmailAlreadyInUse{
			Email: u.Email,
		}
	}

	err := s.db.Persist(ctx, u)
	if err == nil {
		s.mailer.Send(ctx, u)
		return nil
	}

	return err
}
