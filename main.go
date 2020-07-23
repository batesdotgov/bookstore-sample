package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"

	"github.com/diegoholiveira/bookstore-sample/books/booksretriever"
	"github.com/diegoholiveira/bookstore-sample/internal/database"
	"github.com/diegoholiveira/bookstore-sample/internal/schema"
	"github.com/diegoholiveira/bookstore-sample/purchases/purchasepersister"
	"github.com/diegoholiveira/bookstore-sample/users/purchaseshistory"
	"github.com/diegoholiveira/bookstore-sample/users/usersregister"
)

type (
	Middleware  func(http.Handler) http.Handler
	Middlewares []Middleware
)

func NewHTTPRouter(middlewares Middlewares) chi.Router {
	router := chi.NewRouter()

	log.Printf("Registering %d middlewares into the http server", len(middlewares))

	for _, middleware := range middlewares {
		router.Use(middleware)
	}

	return router
}

func StartHTTPServer(ctx context.Context, router chi.Router) {
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	<-ctx.Done()

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("server Shutdown Failed:%+s", err)
	}
}

func CreateDatabaseSchema(ctx context.Context, db *sql.DB) {
	appEnvironment := os.Getenv("APP_ENVIRONMENT")
	if appEnvironment != "local" {
		return
	}
	schema.Up(db)
	<-ctx.Done()
	schema.Down(db)
}

func main() {
	router := InitHTTPRouter()
	db, err := database.NewMySQLConnection()
	if err != nil {
		panic(err)
	}
	welcomeMailer := usersregister.NewWelcomeMailer()

	// Configure our application modules
	booksretriever.SetupModule(router, db)
	purchasepersister.SetupModule(router, db, welcomeMailer)
	usersregister.SetupModule(router, db, welcomeMailer)
	purchaseshistory.SetupModule(router, db)

	// start the application
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-done
		cancel()
	}()
	go CreateDatabaseSchema(ctx, db)
	StartHTTPServer(ctx, router)
}
