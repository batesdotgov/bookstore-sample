package purchaseshistory

import (
	"database/sql"

	"github.com/go-chi/chi"
)

func SetupModule(router chi.Router, db *sql.DB) {
	repo := NewPurchasesRepository(db)
	handler := NewPurchasesHistoryHandler(repo, repo)

	router.Method("GET", "/users/{id}/purchases", handler)
}
