package purchaseshistory

import (
	"github.com/go-chi/chi"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Invoke(registerEndpoints),
	factories,
)

var factories = fx.Provide(
	NewPurchasesHistoryHandler,
	NewPurchasesRepository,
)

func registerEndpoints(router chi.Router, handler PurchasesHistoryHandler) {
	router.Method("GET", "/users/{id}/purchases", handler)
}
