package product

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	ProductMocks "TEST_PetShop/internal/handlers/product/mocks"
	"TEST_PetShop/internal/models"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"
)

// Get Product - Ready
func TestGetAllProducts_Success(t *testing.T) {
	// Мокаем storage с помощью сгенерированного мока
	mockStorage := ProductMocks.NewProducts(t)
	mockStorage.On("GetAllProducts", mock.Anything).Return([]models.Product{{ID: 1, Name: "Dog Food"}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mockStorage)

	handler.GetAllProducts(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestGetAllProducts_Error(t *testing.T) {
	mockStorage := ProductMocks.NewProducts(t)
	mockStorage.On("GetAllProducts", mock.Anything).Return(nil, errors.New("DB error"))

	
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mockStorage)
	handler.GetAllProducts(w, req)

	
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// =======================
// Create Product
// =======================

func TestCreateProduct_Success(t *testing.T) {
	
	mockStorage := ProductMocks.NewProducts(t)
	mockStorage.On("CreateProduct", mock.Anything, mock.Anything).Return(1, nil)

	
	body := strings.NewReader(`{"Name": "Dog Food", "Price": 15.5, "Stock": 10}`)
	req := httptest.NewRequest(http.MethodPost, "/products", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mockStorage)
	handler.CreateProduct(w, req)

	
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestCreateProduct_BadRequest(t *testing.T) {
	mockStorage := ProductMocks.NewProducts(t)

	body := strings.NewReader(`{"Name": invalid_json}`)
	req := httptest.NewRequest(http.MethodPost, "/products", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mockStorage)
	handler.CreateProduct(w, req)

	
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestCreateProduct_Fail(t *testing.T) {
	// Мокаем storage — возвращаем ошибку базы данных
	mockStorage := ProductMocks.NewProducts(t)
	mockStorage.On("CreateProduct", mock.Anything, mock.Anything).Return(0, errors.New("DB error"))

	body := strings.NewReader(`{"Name": "Dog Food", "Price": 15.5, "Stock": 10}`)
	req := httptest.NewRequest(http.MethodPost, "/products", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mockStorage)
	handler.CreateProduct(w, req)

	
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

// =======================
// Update Product
// =======================

func TestUpdateProduct_Success(t *testing.T) {
	// Мокаем storage — успешное обновление
	mockStorage := ProductMocks.NewProducts(t)
	mockStorage.On("UpdateProduct", mock.Anything, mock.Anything).Return(nil)

	
	body := strings.NewReader(`{"Name": "Cat Food", "Price": 20.0, "Stock": 5}`)
	req := httptest.NewRequest(http.MethodPut, "/products/1", body)
	req.Header.Set("Content-Type", "application/json")

	
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	handler := New(slog.Default(), mockStorage)
	handler.UpdateProduct(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestUpdateProduct_BadRequest(t *testing.T) {
	mockStorage := ProductMocks.NewProducts(t)

	
	body := strings.NewReader(`{"Name": invalid_json}`)
	req := httptest.NewRequest(http.MethodPut, "/products/1", body)
	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	handler := New(slog.Default(), mockStorage)
	handler.UpdateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestUpdateProduct_Fail(t *testing.T) {
	// Мокаем storage — ошибка при обновлении
	mockStorage := ProductMocks.NewProducts(t)
	mockStorage.On("UpdateProduct", mock.Anything, mock.Anything).Return(errors.New("DB error"))

	body := strings.NewReader(`{"Name": "Cat Food", "Price": 20.0, "Stock": 5}`)
	req := httptest.NewRequest(http.MethodPut, "/products/1", body)
	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	handler := New(slog.Default(), mockStorage)
	handler.UpdateProduct(w, req)

	
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

// =======================
// Delete Product
// =======================

func TestDeleteProduct_Success(t *testing.T) {
	// Мокаем storage — успешное удаление
	mockStorage := ProductMocks.NewProducts(t)
	mockStorage.On("DeleteProduct", mock.Anything, 1).Return(nil)

	
	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	handler := New(slog.Default(), mockStorage)
	handler.DeleteProduct(w, req)

	
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestDeleteProduct_BadRequest(t *testing.T) {
	mockStorage := ProductMocks.NewProducts(t)

	
	req := httptest.NewRequest(http.MethodDelete, "/products/", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mockStorage)
	handler.DeleteProduct(w, req)

	
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestDeleteProduct_Fail(t *testing.T) {
	// Мокаем storage — ошибка при удалении
	mockStorage := ProductMocks.NewProducts(t)
	mockStorage.On("DeleteProduct", mock.Anything, 1).Return(errors.New("DB error"))

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	handler := New(slog.Default(), mockStorage)
	handler.DeleteProduct(w, req)

	
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}