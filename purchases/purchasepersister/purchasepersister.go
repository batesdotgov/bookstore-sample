package purchasepersister

import (
	"database/sql"

	"github.com/go-chi/chi"

	"github.com/diegoholiveira/bookstore-sample/users/usersregister"
)

func SetupModule(router chi.Router, db *sql.DB, mailer usersregister.WelcomeMailer) {
	handler := InitHandler(db, mailer)

	router.Method("POST", "/purchases", handler)
}
