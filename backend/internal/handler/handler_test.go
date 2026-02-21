package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/gowthamd/go-crud-app/internal/repository"
)

func TestCreateItem_InvalidJSON(t *testing.T) {
	h := NewItemHandler(&repository.ItemRepository{})

	req := httptest.NewRequest(http.MethodPost, "/api/items", strings.NewReader("not-json"))
	rec := httptest.NewRecorder()

	h.CreateItem(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestCreateItem_EmptyName(t *testing.T) {
	h := NewItemHandler(&repository.ItemRepository{})

	body := `{"name":"","price":9.99}`
	req := httptest.NewRequest(http.MethodPost, "/api/items", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.CreateItem(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "name is required") {
		t.Errorf("expected validation error for name, got %s", rec.Body.String())
	}
}

func TestCreateItem_NegativePrice(t *testing.T) {
	h := NewItemHandler(&repository.ItemRepository{})

	body := `{"name":"Widget","price":-5}`
	req := httptest.NewRequest(http.MethodPost, "/api/items", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.CreateItem(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "price must be non-negative") {
		t.Errorf("expected validation error for price, got %s", rec.Body.String())
	}
}

func TestGetItem_InvalidUUID(t *testing.T) {
	h := NewItemHandler(&repository.ItemRepository{})

	req := httptest.NewRequest(http.MethodGet, "/api/items/not-a-uuid", nil)
	rec := httptest.NewRecorder()

	h.GetItem(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestUpdateItem_InvalidJSON(t *testing.T) {
	h := NewItemHandler(&repository.ItemRepository{})

	req := httptest.NewRequest(http.MethodPut, "/api/items/550e8400-e29b-41d4-a716-446655440000", strings.NewReader("not-json"))
	rec := httptest.NewRecorder()

	// Set up chi URL params
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "550e8400-e29b-41d4-a716-446655440000")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.UpdateItem(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestUpdateItem_EmptyName(t *testing.T) {
	h := NewItemHandler(&repository.ItemRepository{})

	body := `{"name":""}`
	req := httptest.NewRequest(http.MethodPut, "/api/items/550e8400-e29b-41d4-a716-446655440000", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "550e8400-e29b-41d4-a716-446655440000")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.UpdateItem(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "name cannot be empty") {
		t.Errorf("expected validation error for name, got %s", rec.Body.String())
	}
}

func TestDeleteItem_InvalidUUID(t *testing.T) {
	h := NewItemHandler(&repository.ItemRepository{})

	req := httptest.NewRequest(http.MethodDelete, "/api/items/not-a-uuid", nil)
	rec := httptest.NewRecorder()

	h.DeleteItem(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}
