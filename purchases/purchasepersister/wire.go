//+build wireinject

package purchasepersister

import (
	"database/sql"

	"github.com/google/wire"

	"github.com/diegoholiveira/bookstore-sample/users/usersregister"
)

func InitHandler(db *sql.DB, mailer usersregister.WelcomeMailer) PurchaseHandler {
	panic(wire.Build(
		NewUserRegisterServiceFactory,
		NewPurchaserWithNewUser,
		NewPurchaserWithRegisteredUser,
		NewPurchaserService,
		NewPurchaseHandler,
	))
}
