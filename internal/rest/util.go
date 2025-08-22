package rest

import (
	"cryptoserver/domain"
	"encoding/json"
	"net/http"
)

type ErrorStruct struct {
	Message string `json:"message"`
}

func WriteError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrNotFound:
		_writeError(w, err.Error(), http.StatusNotFound)
	case domain.ErrInvalidToken:
		_writeError(w, err.Error(), 400)
	case domain.ErrAlreadyExist:
		_writeError(w, err.Error(), http.StatusBadRequest)
	default:
		_writeError(w, err.Error(), http.StatusInternalServerError)
	}
}

func _writeError(w http.ResponseWriter, msg string, statusCode int) {
	es := ErrorStruct{Message: msg}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(es)
}

type Middleware func(http.Handler) http.Handler

type Router struct {
	mx *http.ServeMux
	md []Middleware
}

func NewRouter(mx *http.ServeMux) *Router {
	return &Router{
		mx: mx,
		md: []Middleware{},
	}
}

func (r *Router) ApplyMiddleware(mds ...Middleware) {
	r.md = append(r.md, mds...)
}

func (r *Router) Handle(pattern string, fn http.Handler) {
	for _, middleware := range r.md {
		fn = middleware(fn)
	}
	r.mx.Handle(pattern, fn)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mx.ServeHTTP(w, req)
}
