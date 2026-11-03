# Gateway

HTTP-шлюз на `:8080`. Бизнес-логика и хранение в пакете `ledger` (Postgres + Redis).

Сначала подними БД и Redis, прогони миграции: [../ledger/README.md](../ledger/README.md).

## Запуск

```powershell
$env:DATABASE_URL = "postgres://postgres:postgres@localhost:5432/cashapp?sslmode=disable"
$env:REDIS_ADDR = "localhost:6379"
go run .
```

## Примеры запросов

Бюджет:

```powershell
curl.exe -X POST http://localhost:8080/api/budgets `
  -H "Content-Type: application/json" `
  -d "{\"category\":\"food\",\"limit\":1000}"
```

Транзакция в пределах лимита (ожидаем 201):

```powershell
curl.exe -X POST http://localhost:8080/api/transactions `
  -H "Content-Type: application/json" `
  -d "{\"amount\":450,\"category\":\"food\",\"description\":\"lunch\",\"date\":\"2026-11-01\"}"
```

Превышение лимита (ожидаем 409, в `expenses` записи нет):

```powershell
curl.exe -X POST http://localhost:8080/api/transactions `
  -H "Content-Type: application/json" `
  -d "{\"amount\":700,\"category\":\"food\",\"description\":\"dinner\",\"date\":\"2026-11-02\"}"
```

Списки:

```powershell
curl.exe http://localhost:8080/api/budgets
curl.exe http://localhost:8080/api/transactions
```

Отчёт два раза подряд. В логах Gateway: `cache miss`, затем `cache hit`:

```powershell
curl.exe "http://localhost:8080/api/reports/summary?from=2026-11-01&to=2026-11-30"
curl.exe "http://localhost:8080/api/reports/summary?from=2026-11-01&to=2026-11-30"
```

После перезапуска Gateway данные остаются в PostgreSQL.
