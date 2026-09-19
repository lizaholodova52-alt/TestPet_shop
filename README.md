# Go Pet Shop — V1

В версии V1 добавлены unit-тесты для Product handlers.

Протестированы:

- CreateProduct
- UpdateProduct
- DeleteProduct

Для каждого handler проверяются:

- успешный запрос — `200 OK`
- некорректный запрос — `400 Bad Request`
- внутренняя ошибка — `500 Internal Server Error`

Тесты находятся в:

`internal/handlers/product/products_test.go`

## Run

Запустить все тесты:

```bash
go test ./...
```
