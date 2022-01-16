package handler

import (
	"PasswordService/app"
	"PasswordService/domain/entities"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Handlers interface {
	HashApi(w http.ResponseWriter, r *http.Request)
	GetHashedPassword(w http.ResponseWriter, r *http.Request)
	GetStats(w http.ResponseWriter, r *http.Request)
}

type handlers struct {
	passwordSvc app.PasswordService
	posts       int64
	total       int64
	statsLock   sync.Mutex
}

func NewHandlers(ps app.PasswordService) Handlers {
	return &handlers{
		passwordSvc: ps,
		posts:       0,
		total:       0,
	}
}

func (h *handlers) HashApi(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	log.Printf("HashApi handler invoked")
	if r.Header.Get("Content-Type") != "application/json" {
		h.errorResp(w, "Content-Type is not application/json", http.StatusUnsupportedMediaType)
		h.calcElapsed(start)
		return
	}

	var req app.PasswordDTO
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.errorResp(w, err.Error(), http.StatusBadRequest)
		h.calcElapsed(start)
		return
	}

	data, err := h.passwordSvc.CreatePassword(req)

	if err != nil {
		h.errorResp(w, err.Error(), http.StatusInternalServerError)
		h.calcElapsed(start)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(strconv.FormatInt(data.ID, 10)))
	h.calcElapsed(start)
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

func (h *handlers) GetStats(w http.ResponseWriter, r *http.Request) {
	var resp app.PasswordStatsDTO

	h.statsLock.Lock()
	resp.NumPosts = h.posts
	if resp.NumPosts == 0 {
		resp.Average = 0
	} else {
		resp.Average = h.total / h.posts
	}
	h.statsLock.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *handlers) errorResp(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	r := make(map[string]string)
	r["message"] = message
	j, _ := json.Marshal(r)
	w.Write(j)
}

func (h *handlers) calcElapsed(start time.Time) {
	h.statsLock.Lock()
	defer h.statsLock.Unlock()

	elapsed := time.Since(start)
	h.total += elapsed.Microseconds()
	h.posts++
}
