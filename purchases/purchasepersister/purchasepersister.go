package purchasepersister

import (
	"github.com/go-chi/chi"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Invoke(registerEndpoints),
	factories,
)

var factories = fx.Provide(
	NewPurchaseHandler,
	NewUserRegisterServiceFactory,
	NewPurchaserService,
	NewPurchaserWithNewUser,
	NewPurchaserWithRegisteredUser,
)

func registerEndpoints(router chi.Router, handler PurchaseHandler) {
	router.Method("POST", "/purchases", handler)
}
