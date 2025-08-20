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
