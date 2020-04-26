package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/fx"

	"github.com/diegoholiveira/bookstore-sample/books/booksretriever"
	"github.com/diegoholiveira/bookstore-sample/internal/database"
	"github.com/diegoholiveira/bookstore-sample/internal/schema"
	"github.com/diegoholiveira/bookstore-sample/purchases/purchasepersister"
	"github.com/diegoholiveira/bookstore-sample/users/usersregister"
)

type (
	Middleware func(http.Handler) http.Handler

	RouterParams struct {
		fx.In

		Middlewares []Middleware `group:"middlewares"`
	}
)

func NewHTTPRouter(params RouterParams) chi.Router {
	router := chi.NewRouter()

	log.Printf("Registering %d middlewares into the http server", len(params.Middlewares))

	for _, middleware := range params.Middlewares {
		router.Use(middleware)
	}

	return router
}

func StartHTTPServer(lc fx.Lifecycle, router chi.Router) {
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Print("Starting the server...")

			go srv.ListenAndServe() // nolint:errcheck

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Print("Stopping the server...")

			return srv.Shutdown(ctx)
		},
	})
}

func CreateDatabaseSchema(lc fx.Lifecycle, db *sql.DB) {
	appEnvironment := os.Getenv("APP_ENVIRONMENT")
	if appEnvironment != "local" {
		return
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			schema.Up(db)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			schema.Down(db)
			return nil
		},
	})
}

func main() {
	ServerDependencies := fx.Provide(
		NewHTTPRouter,
		database.NewMySQLConnection,
	)

	fx.New(
		fx.Options(
			Middlewares,
			ServerDependencies,
			usersregister.Module,
			booksretriever.Module,
			purchasepersister.Module,
		),
		fx.Invoke(CreateDatabaseSchema, StartHTTPServer),
	).Run()
}
