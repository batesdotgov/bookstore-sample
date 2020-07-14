package purchasepersister

import (
	"database/sql"

	"github.com/go-chi/chi"

	"github.com/diegoholiveira/bookstore-sample/users/usersregister"
)

func SetupModule(router chi.Router, db *sql.DB, mailer usersregister.WelcomeMailer) {
	registrationServiceFactory := NewUserRegisterServiceFactory(mailer)
	newUserRoutine := NewPurchaserWithNewUser(db, registrationServiceFactory)
	registeredUserRoutine := NewPurchaserWithRegisteredUser(db)
	purchaser := NewPurchaserService(registeredUserRoutine, newUserRoutine)
	handler := NewPurchaseHandler(purchaser)

	router.Method("POST", "/purchases", handler)
}
