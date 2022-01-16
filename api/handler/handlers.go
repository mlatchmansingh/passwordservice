package handler

import (
	"PasswordService/app"
	"PasswordService/domain/entities"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Handlers interface {
	HashApi(w http.ResponseWriter, r *http.Request)
	GetHashedPassword(w http.ResponseWriter, r *http.Request)
}

type handlers struct {
	passwordSvc app.PasswordService
}

func NewHandlers(ps app.PasswordService) Handlers {
	return &handlers{
		passwordSvc: ps,
	}
}

func (h *handlers) HashApi(w http.ResponseWriter, r *http.Request) {

	log.Printf("HashApi handler invoked")
	if r.Header.Get("Content-Type") != "application/json" {
		h.errorResp(w, "Content-Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	var req app.PasswordDTO
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.errorResp(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := h.passwordSvc.CreatePassword(req)

	if err != nil {
		h.errorResp(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(strconv.FormatInt(data.ID, 10)))
}

func (h *handlers) GetHashedPassword(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/hash/")
	if id == "" {
		h.errorResp(w, "Missing id", http.StatusBadRequest)
	}

	eid, err := strconv.Atoi(id)
	if err != nil {
		h.errorResp(w, err.Error(), http.StatusBadRequest)
		return
	}

	password, err := h.passwordSvc.GetPassword(entities.ID(eid))
	if err != nil {
		h.errorResp(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !password.Converted {
		h.errorResp(w, "Password not ready yet", http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(password.Password))

}

func (h *handlers) errorResp(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	r := make(map[string]string)
	r["message"] = message
	j, _ := json.Marshal(r)
	w.Write(j)
}
