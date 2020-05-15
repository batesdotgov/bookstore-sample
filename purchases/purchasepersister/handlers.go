//go:generate mockery -name=Purchaser -output=./internal/mocks

package purchasepersister

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/diegoholiveira/bookstore-sample/pkg/http/render"
	"github.com/diegoholiveira/bookstore-sample/purchases"
)

type (
	Purchaser interface {
		MakePurchase(context.Context, purchases.Purchase) error
	}

	PurchaseHandler struct {
		purchaser Purchaser
	}
)

func NewPurchaseHandler(p Purchaser) PurchaseHandler {
	return PurchaseHandler{
		purchaser: p,
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

	err = h.purchaser.MakePurchase(r.Context(), purchase)
	if err == nil {
		w.WriteHeader(http.StatusCreated)
		return
	}

	var errPurchase ErrPurchaseInvalid
	if errors.As(err, &errPurchase) {
		render.JSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	render.JSON(w, http.StatusInternalServerError, map[string]string{
		"error": err.Error(),
	})
}
