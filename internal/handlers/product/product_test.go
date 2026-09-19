package product


import (
    "context"
    "errors"
    "log/slog"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"


    "TEST_PetShop/internal/models"


    "github.com/go-chi/chi"
)


// Get Product - Ready
func TestGetAllProducts_Success(t *testing.T) {
    // Мокаем storage — он вернёт один продукт.
    mock := &ProductsMock{
        GetAllProductsFunc: func(ctx context.Context) ([]models.Product, error) {
            return []models.Product{
                {ID: 1, Name: "Dog Food"},
            }, nil
        },
    }


    // Создаем HTTP-запрос GET /products
    req := httptest.NewRequest(http.MethodGet, "/products", nil)
    w := httptest.NewRecorder()


    // Создаем хендлер с мок-хранилищем
    handler := New(slog.Default(), mock)


    // Вызываем метод GetAllProducts, который является http.HandlerFunc
    handler.GetAllProducts(w, req)


    // Проверяем HTTP-код
    if w.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", w.Code)
    }
}
func TestGetAllProducts_Error(t *testing.T) {
    // Мокаем storage — он будет возвращать ошибку
    mock := &ProductsMock{
        GetAllProductsFunc: func(ctx context.Context) ([]models.Product, error) {
            return nil, errors.New("DB error")
        },
    }


    // Создаем запрос
    req := httptest.NewRequest(http.MethodGet, "/products", nil)
    w := httptest.NewRecorder()


    handler := New(slog.Default(), mock)
    handler.GetAllProducts(w, req)


    // Ожидаем HTTP 500
    if w.Code != http.StatusInternalServerError {
        t.Fatalf("expected 500, got %d", w.Code)
    }
}


// =======================
// Create Product
// =======================


func TestCreateProduct_Success(t *testing.T) {
    // Мокаем storage — он вернет ID созданного продукта
    mock := &ProductsMock{
        CreateProductFunc: func(ctx context.Context, product models.Product) (int, error) {
            return 1, nil
        },
    }


    // Создаем HTTP-запрос POST /products с валидным JSON
    body := strings.NewReader(`{"Name": "Dog Food", "Price": 15.5, "Stock": 10}`)
    req := httptest.NewRequest(http.MethodPost, "/products", body)
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()


    handler := New(slog.Default(), mock)
    handler.CreateProduct(w, req)


    // Проверяем HTTP-код 200 OK
    if w.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", w.Code)
    }
}


func TestCreateProduct_BadRequest(t *testing.T) {
    mock := &ProductsMock{}


    // Создаем HTTP-запрос с невалидным JSON
    body := strings.NewReader(`{"Name": invalid_json}`)
    req := httptest.NewRequest(http.MethodPost, "/products", body)
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()


    handler := New(slog.Default(), mock)
    handler.CreateProduct(w, req)


    // Ожидаем HTTP 400 Bad Request
    if w.Code != http.StatusBadRequest {
        t.Fatalf("expected status 400, got %d", w.Code)
    }
}


func TestCreateProduct_Fail(t *testing.T) {
    // Мокаем storage — возвращаем ошибку базы данных
    mock := &ProductsMock{
        CreateProductFunc: func(ctx context.Context, product models.Product) (int, error) {
            return 0, errors.New("DB error")
        },
    }


    body := strings.NewReader(`{"Name": "Dog Food", "Price": 15.5, "Stock": 10}`)
    req := httptest.NewRequest(http.MethodPost, "/products", body)
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()


    handler := New(slog.Default(), mock)
    handler.CreateProduct(w, req)


    // Ожидаем HTTP 500 Internal Server Error
    if w.Code != http.StatusInternalServerError {
        t.Fatalf("expected status 500, got %d", w.Code)
    }
}


// =======================
// Update Product
// =======================


func TestUpdateProduct_Success(t *testing.T) {
    // Мокаем storage — успешное обновление
    mock := &ProductsMock{
        UpdateProductFunc: func(ctx context.Context, product models.Product) error {
            return nil
        },
    }


    // Создаем HTTP-запрос PUT /products/1
    body := strings.NewReader(`{"Name": "Cat Food", "Price": 20.0, "Stock": 5}`)
    req := httptest.NewRequest(http.MethodPut, "/products/1", body)
    req.Header.Set("Content-Type", "application/json")


    // Передаем URL-параметр id через chi context
    rctx := chi.NewRouteContext()
    rctx.URLParams.Add("id", "1")
    req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))


    w := httptest.NewRecorder()


    handler := New(slog.Default(), mock)
    handler.UpdateProduct(w, req)


    // Проверяем HTTP-код 200 OK
    if w.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", w.Code)
    }
}


func TestUpdateProduct_BadRequest(t *testing.T) {
    mock := &ProductsMock{}


    // Создаем HTTP-запрос с невалидным JSON
    body := strings.NewReader(`{"Name": invalid_json}`)
    req := httptest.NewRequest(http.MethodPut, "/products/1", body)
    req.Header.Set("Content-Type", "application/json")


    rctx := chi.NewRouteContext()
    rctx.URLParams.Add("id", "1")
    req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))


    w := httptest.NewRecorder()


    handler := New(slog.Default(), mock)
    handler.UpdateProduct(w, req)


    // Ожидаем HTTP 400 Bad Request
    if w.Code != http.StatusBadRequest {
        t.Fatalf("expected status 400, got %d", w.Code)
    }
}


func TestUpdateProduct_Fail(t *testing.T) {
    // Мокаем storage — ошибка при обновлении
    mock := &ProductsMock{
        UpdateProductFunc: func(ctx context.Context, product models.Product) error {
            return errors.New("DB error")
        },
    }


    body := strings.NewReader(`{"Name": "Cat Food", "Price": 20.0, "Stock": 5}`)
    req := httptest.NewRequest(http.MethodPut, "/products/1", body)
    req.Header.Set("Content-Type", "application/json")


    rctx := chi.NewRouteContext()
    rctx.URLParams.Add("id", "1")
    req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))


    w := httptest.NewRecorder()


    handler := New(slog.Default(), mock)
    handler.UpdateProduct(w, req)


    // Ожидаем HTTP 500 Internal Server Error
    if w.Code != http.StatusInternalServerError {
        t.Fatalf("expected status 500, got %d", w.Code)
    }
}


// =======================
// Delete Product
// =======================


func TestDeleteProduct_Success(t *testing.T) {
    // Мокаем storage — успешное удаление
    mock := &ProductsMock{
        DeleteProductFunc: func(ctx context.Context, id int) error {
            return nil
        },
    }


    // Создаем HTTP-запрос DELETE /products/1
    req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)


    rctx := chi.NewRouteContext()
    rctx.URLParams.Add("id", "1")
    req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))


    w := httptest.NewRecorder()


    handler := New(slog.Default(), mock)
    handler.DeleteProduct(w, req)


    // Проверяем HTTP-код 200 OK
    if w.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", w.Code)
    }
}


func TestDeleteProduct_BadRequest(t *testing.T) {
    mock := &ProductsMock{}


    // Создаем HTTP-запрос без параметра id (пустой id)
    req := httptest.NewRequest(http.MethodDelete, "/products/", nil)
    w := httptest.NewRecorder()


    handler := New(slog.Default(), mock)
    handler.DeleteProduct(w, req)


    // Ожидаем HTTP 400 Bad Request
    if w.Code != http.StatusBadRequest {
        t.Fatalf("expected status 400, got %d", w.Code)
    }
}


func TestDeleteProduct_Fail(t *testing.T) {
    // Мокаем storage — ошибка при удалении
    mock := &ProductsMock{
        DeleteProductFunc: func(ctx context.Context, id int) error {
            return errors.New("DB error")
        },
    }


    // Создаем HTTP-запрос DELETE /products/1
    req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)


    rctx := chi.NewRouteContext()
    rctx.URLParams.Add("id", "1")
    req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))


    w := httptest.NewRecorder()


    handler := New(slog.Default(), mock)
    handler.DeleteProduct(w, req)


    // Ожидаем HTTP 500 Internal Server Error
    if w.Code != http.StatusInternalServerError {
        t.Fatalf("expected status 500, got %d", w.Code)
	}
}