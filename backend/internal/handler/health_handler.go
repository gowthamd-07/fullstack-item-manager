package handler

import (
	"net/http"

	"github.com/gowthamd/go-crud-app/internal/db"
)

type HealthHandler struct {
	DB *db.DB
}

func NewHealthHandler(db *db.DB) *HealthHandler {
	return &HealthHandler{DB: db}
}

func (h *HealthHandler) Liveness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	if err := h.DB.Ping(r.Context()); err != nil {
		http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ready"))
}
