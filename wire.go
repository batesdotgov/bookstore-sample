//+build wireinject

package main

import (
	"github.com/go-chi/chi"
	"github.com/google/wire"
)

func InitHTTPRouter() chi.Router {
	panic(wire.Build(
		NewHTTPRouter,
		NewMiddlewares,
	))
}
