package usersregister

import (
	"database/sql"

	"github.com/go-chi/chi"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Invoke(registerEndpoints),
	factories,
)

var factories = fx.Provide(
	NewUserRegisterHandler,
	NewUserRegisterService,
	NewUserConstraints,
	NewDatabasePersister,
	NewWelcomeMailer,
	// translate *sql.DB into our database interfaces
	newDatabaseExecutor,
	newDatabaseQuerier,
)

func newDatabaseExecutor(db *sql.DB) DatabaseExecutor {
	return db
}

func newDatabaseQuerier(db *sql.DB) DatabaseQuerier {
	return db
}

func registerEndpoints(router chi.Router, handler UserRegisterHandler) {
	router.Method("POST", "/users", handler)
}
