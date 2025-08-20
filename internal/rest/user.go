package rest

import (
	"cryptoserver/domain"
	"encoding/json"
	"net/http"
)

type UserService interface {
	Register(user *domain.User) (string, error)
	Login(user *domain.User) (string, error)
}

type UserHandler struct {
	us UserService
}

type TokenResponse struct {
	Token string `json:"token"`
}

func NewUserHandler(us UserService, mx *http.ServeMux) {
	uh := &UserHandler{
		us: us,
	}

	mx.HandleFunc("POST /auth/register", uh.Register)
	mx.HandleFunc("POST /auth/login", uh.Login)
}

func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user := &domain.User{}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		WriteError(w, err)
		return
	}

	token, err := uh.us.Register(user)
	if err != nil {
		WriteError(w, err)
		return
	}

	resp := TokenResponse{
		Token: token,
	}

	json.NewEncoder(w).Encode(resp)
}

func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user := &domain.User{}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		WriteError(w, err)
		return
	}

	token, err := uh.us.Login(user)
	if err != nil {
		WriteError(w, err)
		return
	}

	resp := TokenResponse{
		Token: token,
	}

	json.NewEncoder(w).Encode(resp)
}
