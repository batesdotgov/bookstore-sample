package purchasepersister

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/diegoholiveira/bookstore-sample/purchases"
	"github.com/diegoholiveira/bookstore-sample/purchases/purchasepersister/internal/mocks"
)

func TestPurchase_success(t *testing.T) {
	input := `{
		"user_id": 1,
		"books": [
			{
				"id": 1,
				"quantity": 2,
				"price": 19.9
			}
		]
	}`

	purchase := purchases.Purchase{
		UserID: 1,
		Books: purchases.Books{
			{
				ID:       1,
				Quantity: 2,
				Price:    19.9,
			},
		},
	}

	r := httptest.NewRequest("POST", "/purchases", strings.NewReader(input))
	w := httptest.NewRecorder()

	persister := new(mocks.PurchasePersister)
	persister.
		On("Persist", mock.Anything, purchase).
		Return(nil)

	h := NewPurchaseHandler(persister)
	h.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if assert.Equal(t, http.StatusCreated, resp.StatusCode) {
		persister.AssertExpectations(t)
	}
}
