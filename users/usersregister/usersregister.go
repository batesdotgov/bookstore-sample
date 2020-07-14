package usersregister

import (
	"database/sql"

	"github.com/go-chi/chi"
)

func SetupModule(router chi.Router, db *sql.DB, mailer WelcomeMailer) {
	constraints := NewUserConstraints(db)
	persister := NewDatabasePersister(db)
	service := NewUserRegisterService(constraints, persister, mailer)
	handler := NewUserRegisterHandler(service)

	router.Method("POST", "/users", handler)
}
