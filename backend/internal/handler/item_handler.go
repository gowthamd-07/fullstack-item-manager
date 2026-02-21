package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gowthamd/go-crud-app/internal/models"
	"github.com/gowthamd/go-crud-app/internal/repository"
)

type ItemHandler struct {
	Repo *repository.ItemRepository
}

func NewItemHandler(repo *repository.ItemRepository) *ItemHandler {
	return &ItemHandler{Repo: repo}
}

func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var dto models.CreateItemDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := dto.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item, err := h.Repo.Create(r.Context(), &dto)
	if err != nil {
		http.Error(w, "Failed to create item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *ItemHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	item, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrItemNotFound) {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *ItemHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	items, total, err := h.Repo.GetAll(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, "Failed to list items", http.StatusInternalServerError)
		return
	}

	if items == nil {
		items = []models.Item{}
	}

	resp := models.PaginatedResponse{
		Items:  items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *ItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var dto models.UpdateItemDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := dto.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item, err := h.Repo.Update(r.Context(), id, &dto)
	if err != nil {
		if errors.Is(err, repository.ErrItemNotFound) {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *ItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repository.ErrItemNotFound) {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
