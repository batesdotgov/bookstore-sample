package usersregister

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/diegoholiveira/bookstore-sample/users"
	"github.com/diegoholiveira/bookstore-sample/users/usersregister/internal/mocks"
)

func TestUserRegisterHandler_perfect_payload(t *testing.T) {
	input := `{
		"name": "Diego Henrique Oliveira",
		"email": "contact@diegoholiveira.com",
		"password": "12345678"
	}`

	r := httptest.NewRequest("POST", "/users", strings.NewReader(input))
	w := httptest.NewRecorder()

	u := users.User{
		Name:     "Diego Henrique Oliveira",
		Email:    "contact@diegoholiveira.com",
		Password: "12345678",
	}

	persister := new(mocks.UserRegister)
	persister.
		On("Register", mock.Anything, &u).
		Return(nil)

	h := NewUserRegisterHandler(persister)
	h.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if assert.Equal(t, http.StatusCreated, resp.StatusCode) {
		persister.AssertExpectations(t)

		body, _ := ioutil.ReadAll(resp.Body)
		assert.Equal(t, "", string(body))
	}
}

func TestUserRegisterHandler_perfect_payload_twice(t *testing.T) {
	input := `{
		"name": "Diego Henrique Oliveira",
		"email": "contact@diegoholiveira.com",
		"password": "123456"
	}{
		"name": "Diego Henrique Oliveira",
		"email": "contact@diegoholiveira.com",
		"password": "123456"
	}`

	r := httptest.NewRequest("POST", "/users", strings.NewReader(input))
	w := httptest.NewRecorder()

	persister := new(mocks.UserRegister)

	h := NewUserRegisterHandler(persister)
	h.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	expected := `{
		"error": "Request body must only contain a single JSON object"
	}`

	if assert.Equal(t, http.StatusBadRequest, resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		assert.JSONEq(t, expected, string(body))
	}
}

func TestUserRegisterHandler_malformed_json(t *testing.T) {
	input := `{
		"name": "Diego Henrique Oliveira",
		"email": "contact@diegoholiveira.com",
		"password": "123456",
	}`

	r := httptest.NewRequest("POST", "/users", strings.NewReader(input))
	w := httptest.NewRecorder()

	persister := new(mocks.UserRegister)

	h := NewUserRegisterHandler(persister)
	h.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	expected := `{
		"error": "Error while decoding the JSON payload"
	}`

	if assert.Equal(t, http.StatusBadRequest, resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		assert.JSONEq(t, expected, string(body))
	}
}

func TestUserRegisterHandler_payload_without_name(t *testing.T) {
	input := `{
		"email": "contact@diegoholiveira.com",
		"password": "12345678"
	}`

	r := httptest.NewRequest("POST", "/users", strings.NewReader(input))
	w := httptest.NewRecorder()

	persister := new(mocks.UserRegister)

	h := NewUserRegisterHandler(persister)
	h.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	expected := `{
		"name": "cannot be blank"
	}`

	if assert.Equal(t, http.StatusBadRequest, resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)

		assert.JSONEq(t, expected, string(body))
	}
}

func TestUserRegisterHandler_payload_with_an_invalid_email(t *testing.T) {
	input := `{
		"name": "Diego Henrique Oliveira",
		"email": "this-is-not-an-email",
		"password": "12345678"
	}`

	r := httptest.NewRequest("POST", "/users", strings.NewReader(input))
	w := httptest.NewRecorder()

	persister := new(mocks.UserRegister)

	h := NewUserRegisterHandler(persister)
	h.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	expected := `{
		"email": "must be a valid email address"
	}`

	if assert.Equal(t, http.StatusBadRequest, resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)

		assert.JSONEq(t, expected, string(body))
	}
}

func TestUserRegisterHandler_mysql_is_down(t *testing.T) {
	input := `{
		"name": "Diego Henrique Oliveira",
		"email": "contact@diegoholiveira.com",
		"password": "12345678"
	}`

	r := httptest.NewRequest("POST", "/users", strings.NewReader(input))
	w := httptest.NewRecorder()

	u := users.User{
		Name:     "Diego Henrique Oliveira",
		Email:    "contact@diegoholiveira.com",
		Password: "12345678",
	}

	persister := new(mocks.UserRegister)
	persister.
		On("Register", mock.Anything, &u).
		Return(errors.New("MySQL is down, please try again"))

	h := NewUserRegisterHandler(persister)
	h.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	expected := `{
		"error": "MySQL is down, please try again"
	}`

	if assert.Equal(t, http.StatusInternalServerError, resp.StatusCode) {
		persister.AssertExpectations(t)

		body, _ := ioutil.ReadAll(resp.Body)
		assert.JSONEq(t, expected, string(body))
	}
}

func TestUserRegisterHandler_user_already_registered(t *testing.T) {
	input := `{
		"name": "Diego Henrique Oliveira",
		"email": "contact@diegoholiveira.com",
		"password": "12345678"
	}`

	r := httptest.NewRequest("POST", "/users", strings.NewReader(input))
	w := httptest.NewRecorder()

	u := users.User{
		Name:     "Diego Henrique Oliveira",
		Email:    "contact@diegoholiveira.com",
		Password: "12345678",
	}

	persister := new(mocks.UserRegister)
	persister.
		On("Register", mock.Anything, &u).
		Return(users.ErrEmailAlreadyInUse{Email: u.Email})

	h := NewUserRegisterHandler(persister)
	h.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	expected := `{
		"error": "The e-mail 'contact@diegoholiveira.com' is already in use"
	}`

	if assert.Equal(t, http.StatusBadRequest, resp.StatusCode) {
		persister.AssertExpectations(t)

		body, _ := ioutil.ReadAll(resp.Body)
		assert.JSONEq(t, expected, string(body))
	}
}
