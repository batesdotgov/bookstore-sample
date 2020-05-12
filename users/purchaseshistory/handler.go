//go:generate mockery -name=PurchasesFinder -output=./internal/mocks
//go:generate mockery -name=UserFinder -output=./internal/mocks
package purchaseshistory

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/diegoholiveira/bookstore-sample/pkg/http/render"
	"github.com/diegoholiveira/bookstore-sample/users"
)

type (
	PurchasesHistoryHandler struct {
		userFinder      UserFinder
		purchasesFinder PurchasesFinder
	}

	PurchasesFinder interface {
		FindPurchasesByUser(context.Context, users.User) (users.Purchases, error)
	}

	UserFinder interface {
		FindUserByID(context.Context, uint64) (users.User, error)
	}
)

func NewPurchasesHistoryHandler(userFinder UserFinder, purchasesFinder PurchasesFinder) PurchasesHistoryHandler {
	return PurchasesHistoryHandler{
		userFinder:      userFinder,
		purchasesFinder: purchasesFinder,
	}
}

func (h PurchasesHistoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid user id",
		})
		return // stop here
	}

	var errNotFound users.ErrUserNotFound
	user, err := h.userFinder.FindUserByID(r.Context(), userID)
	if errors.As(err, &errNotFound) {
		render.JSON(w, http.StatusNotFound, map[string]string{
			"error": errNotFound.Error(),
		})

		return // stop here
	}

	purchases, _ := h.purchasesFinder.FindPurchasesByUser(r.Context(), user)

	render.JSON(w, http.StatusOK, purchases)
}
