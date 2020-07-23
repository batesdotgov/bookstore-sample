//+build wireinject

package usersregister

import (
	"database/sql"

	"github.com/google/wire"
)

func InitHandler(db *sql.DB, mailer WelcomeMailer) UserRegisterHandler {
	panic(wire.Build(
		wire.Bind(new(DatabaseQuerier), new(*sql.DB)),
		wire.Bind(new(DatabaseExecutor), new(*sql.DB)),
		NewUserConstraints,
		NewDatabasePersister,
		NewUserRegisterService,
		NewUserRegisterHandler,
	))
}
