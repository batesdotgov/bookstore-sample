package purchaseshistory

import (
	"database/sql"

	"github.com/go-chi/chi"
)

func SetupModule(router chi.Router, db *sql.DB) {
	handler := InitHandler(db)

	router.Method("GET", "/users/{id}/purchases", handler)
}
