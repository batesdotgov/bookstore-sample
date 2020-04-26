//go:generate mockery -name=PurchasePersister -output=./internal/mocks

package purchasepersister

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/diegoholiveira/bookstore-sample/pkg/http/render"
	"github.com/diegoholiveira/bookstore-sample/purchases"
)

type (
	PurchasePersister interface {
		Persist(context.Context, purchases.Purchase) error
	}

	PurchaseHandler struct {
		persister PurchasePersister
	}
)

func NewPurchaseHandler(persister PurchasePersister) PurchaseHandler {
	return PurchaseHandler{
		persister: persister,
	}
}

func (h PurchaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var purchase purchases.Purchase
	err = dec.Decode(&purchase)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{
			"error": "Error while decoding the JSON payload",
		})
		return
	}

	if dec.More() {
		render.JSON(w, http.StatusBadRequest, map[string]string{
			"error": "Request body must only contain a single JSON object",
		})
		return
	}

	err = h.persister.Persist(r.Context(), purchase)
	if err != nil {
		render.JSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	w.WriteHeader(http.StatusCreated)
}
