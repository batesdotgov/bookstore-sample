package usersregister

import (
	"github.com/go-chi/chi"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Invoke(registerEndpoints),
	factories,
)

var factories = fx.Provide(
	NewUserRegisterHandler,
	NewUserPersister,
)

func registerEndpoints(router chi.Router, handler UserRegisterHandler) {
	router.Method("POST", "/users", handler)
}
