//go:generate mockery -name=UserPersister -output=./internal/mocks
package usersregister

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/diegoholiveira/bookstore-sample/pkg/http/render"
	"github.com/diegoholiveira/bookstore-sample/users"
)

type (
	UserPersister interface {
		Persist(context.Context, users.User) error
	}

	UserRegisterHandler struct {
		persister UserPersister
	}
)

func NewUserRegisterHandler(persister UserPersister) UserRegisterHandler {
	return UserRegisterHandler{
		persister: persister,
	}
}

func (h UserRegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var u users.User
	err = dec.Decode(&u)
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

	err = u.ValidateRegisterInput()
	if err != nil {
		render.JSON(w, http.StatusBadRequest, err)
		return
	}

	err = h.persister.Persist(r.Context(), u)
	if err != nil {
		render.JSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	w.WriteHeader(http.StatusCreated)
}
