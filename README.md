# Go Pet Shop – V2

Версия V2 продолжает тестирование Product handlers и использует mock, сгенерированный через `mockery`

## V1

В версии V1 были добавлены unit-тесты для:

- `GetAllProducts`
- `CreateProduct`
- `UpdateProduct`
- `DeleteProduct`

Тесты создавали HTTP-запросы через `httptest`, вызывали handler и проверяли HTTP status code.

Проверялись основные сценарии:

- успешный запрос – `200 OK`
- некорректный запрос – `400 Bad Request`
- внутренняя ошибка – `500 Internal Server Error`

## V2

В версии V2 ручной mock заменён на mock, автоматически сгенерированный с помощью `mockery`.

Для интерфейса:

```go
type Products interface {
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	CreateProduct(ctx context.Context, product models.Product) (int, error)
	DeleteProduct(ctx context.Context, id int) error
	UpdateProduct(ctx context.Context, product models.Product) error
}
```

Для генерации используется директива:

```
//go:generate go run github.com/vektra/mockery/v2 --name=Products
```

## Генерация mock

```bash
go generate ./...
```

## Запуск тестов

```bash
go test ./...
```