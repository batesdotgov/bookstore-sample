//go:generate mockery -name=UserRegister -output=./internal/mocks
package usersregister

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/diegoholiveira/bookstore-sample/pkg/http/render"
	"github.com/diegoholiveira/bookstore-sample/users"
)

type (
	UserRegister interface {
		Register(context.Context, *users.User) error
	}

	UserRegisterHandler struct {
		persister UserRegister
	}
)

func NewUserRegisterHandler(persister UserRegister) UserRegisterHandler {
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

	var emailInUseErr users.ErrEmailAlreadyInUse

	err = h.persister.Register(r.Context(), &u)
	if errors.As(err, &emailInUseErr) {
		render.JSON(w, http.StatusBadRequest, map[string]string{
			"error": emailInUseErr.Error(),
		})
		return
	}

	if err == nil {
		w.WriteHeader(http.StatusCreated)
		return
	}

	render.JSON(w, http.StatusInternalServerError, map[string]string{
		"error": err.Error(),
	})
}
