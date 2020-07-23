//+build wireinject

package purchaseshistory

import (
	"database/sql"

	"github.com/google/wire"
)

func InitHandler(db *sql.DB) PurchasesHistoryHandler {
	panic(wire.Build(
		NewPurchasesRepository,
		wire.Bind(new(PurchasesFinder), new(PurchasesRepository)),
		wire.Bind(new(UserFinder), new(PurchasesRepository)),
		NewPurchasesHistoryHandler,
	))
}
