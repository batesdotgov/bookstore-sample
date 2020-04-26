package main

import (
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/fx"
)

type MiddlewareOutput struct {
	fx.Out

	Middleware Middleware `group:"middlewares"`
}

const DefaultHTTPTimeout = 30 * time.Second

var Middlewares = fx.Provide(
	NewTimeoutMiddleware,
	NewHTTPRecovererMiddleware,
	NewRequestIDMiddleware,
)

func NewTimeoutMiddleware() MiddlewareOutput {
	return MiddlewareOutput{
		Middleware: middleware.Timeout(DefaultHTTPTimeout),
	}
}

func NewHTTPRecovererMiddleware() MiddlewareOutput {
	return MiddlewareOutput{
		Middleware: middleware.Recoverer,
	}
}

func NewRequestIDMiddleware() MiddlewareOutput {
	return MiddlewareOutput{
		Middleware: middleware.RequestID,
	}
}
