package purchasepersister

import (
	"database/sql"

	"github.com/diegoholiveira/bookstore-sample/users/usersregister"
)

type (
	UserRegisterServiceFactory struct {
		mailer usersregister.WelcomeMailer
	}
)

func NewUserRegisterServiceFactory(mailer usersregister.WelcomeMailer) UserRegisterServiceFactory {
	return UserRegisterServiceFactory{
		mailer: mailer,
	}
}

func (factory UserRegisterServiceFactory) Create(tx *sql.Tx) usersregister.UserRegister {
	constraints := usersregister.NewUserConstraints(tx)
	persister := usersregister.NewDatabasePersister(tx)

	return usersregister.NewUserRegisterService(
		constraints,
		persister,
		factory.mailer,
	)
}
