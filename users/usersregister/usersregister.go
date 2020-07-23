package usersregister

import (
	"database/sql"

	"github.com/go-chi/chi"
)

func SetupModule(router chi.Router, db *sql.DB, mailer WelcomeMailer) {
	handler := InitHandler(db, mailer)

	router.Method("POST", "/users", handler)
}
