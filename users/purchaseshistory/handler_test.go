package purchaseshistory

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/diegoholiveira/bookstore-sample/users"
	"github.com/diegoholiveira/bookstore-sample/users/purchaseshistory/internal/mocks"
)

func TestPurchasesHistoryHandler_with_user_not_registered(t *testing.T) {
	r := httptest.NewRequest("GET", "/users/1/purchases", nil)
	w := httptest.NewRecorder()

	var (
		userFinder      = new(mocks.UserFinder)
		purchasesFinder = new(mocks.PurchasesFinder)
	)

	var expectedID uint64 = 1

	userFinder.
		On("FindUserByID", mock.Anything, expectedID).
		Return(users.User{}, users.ErrUserNotFound{})

	h := createHandler(userFinder, purchasesFinder)
	h.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	userFinder.AssertExpectations(t)

	expected := `{
		"error": "User not found"
	}`

	if assert.Equal(t, http.StatusNotFound, resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		assert.JSONEq(t, expected, string(body))
	}
}

func TestPurchasesHistoryHandler_with_invalid_user_id(t *testing.T) {
	r := httptest.NewRequest("GET", "/users/abc/purchases", nil)
	w := httptest.NewRecorder()

	var (
		userFinder      = new(mocks.UserFinder)
		purchasesFinder = new(mocks.PurchasesFinder)
	)

	h := createHandler(userFinder, purchasesFinder)
	h.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	expected := `{
		"error": "Invalid user id"
	}`

	if assert.Equal(t, http.StatusBadRequest, resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		assert.JSONEq(t, expected, string(body))
	}
}

func TestPurchasesHistoryHandler_successfully(t *testing.T) {
	user := users.User{
		ID: 1,
	}

	purchases := users.Purchases{
		{
			ID:     1,
			Amount: 19.9,
			Books: users.Books{
				{
					Title:    "The Hitchhiker's Guide to the Galaxy",
					Author:   "Douglas Adams",
					Price:    19.9,
					Quantity: 1,
				},
			},
		},
		{
			ID:     2,
			Amount: 59.7,
			Books: users.Books{
				{
					Title:    "The Restaurant at the End of the Universe",
					Author:   "Douglas Adams",
					Price:    19.9,
					Quantity: 1,
				},
				{
					Title:    "Life, the Universe and Everything",
					Author:   "Douglas Adams",
					Price:    19.9,
					Quantity: 1,
				},
				{
					Title:    "So Long, and Thanks For All the Fish",
					Author:   "Douglas Adams",
					Price:    19.9,
					Quantity: 1,
				},
			},
		},
	}

	r := httptest.NewRequest("GET", "/users/1/purchases", nil)
	w := httptest.NewRecorder()

	var (
		userFinder      = new(mocks.UserFinder)
		purchasesFinder = new(mocks.PurchasesFinder)
	)

	var expectedUserID uint64 = 1

	userFinder.
		On("FindUserByID", mock.Anything, expectedUserID).
		Return(user, nil)

	purchasesFinder.
		On("FindPurchasesByUser", mock.Anything, user).
		Return(purchases, nil)

	h := createHandler(userFinder, purchasesFinder)
	h.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	expected := `[
		{
			"id": 1,
			"amount": 19.9,
			"books": [
				{
					"title": "The Hitchhiker's Guide to the Galaxy",
					"author": "Douglas Adams",
					"price": 19.9,
					"quantity": 1
				}
			]
		},
		{
			"id": 2,
			"amount": 59.7,
			"books": [
				{
					"title": "The Restaurant at the End of the Universe",
					"author": "Douglas Adams",
					"price": 19.9,
					"quantity": 1
				},
				{
					"title": "Life, the Universe and Everything",
					"author": "Douglas Adams",
					"price": 19.9,
					"quantity": 1
				},
				{
					"title": "So Long, and Thanks For All the Fish",
					"author": "Douglas Adams",
					"price": 19.9,
					"quantity": 1
				}
			]
		}
	]`

	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		assert.JSONEq(t, expected, string(body))
	}
}

// createHandler give us a handler wrapped with the chi.Router
// to allow chi.URLParam work properly
func createHandler(userFinder UserFinder, purchasesFinder PurchasesFinder) chi.Router {
	r := chi.NewRouter()
	r.Method("GET", "/users/{id}/purchases", NewPurchasesHistoryHandler(userFinder, purchasesFinder))
	return r
}
