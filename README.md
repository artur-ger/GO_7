# Домашнее задание 7: PostgreSQL и Redis

Два модуля:

- `ledger` — бюджеты и транзакции в PostgreSQL, отчёты с кешем в Redis;
- `gateway` — HTTP API: `/api/budgets`, `/api/transactions`, `/api/reports/summary`.

Эндпоинты у Gateway те же. Данные теперь лежат в БД и переживают перезапуск.

Нужны Go 1.22+, Docker (или локальные PostgreSQL/Redis) и [goose](https://github.com/pressly/goose).

## Запуск

```powershell
docker run -d --name cashcraft-postgres `
  -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=cashapp `
  -p 5432:5432 postgres:16-alpine

docker run -d --name cashcraft-redis -p 6379:6379 redis:7-alpine

go install github.com/pressly/goose/v3/cmd/goose@latest
$env:DATABASE_URL = "postgres://postgres:postgres@localhost:5432/cashapp?sslmode=disable"
goose -dir ./ledger/migrations postgres $env:DATABASE_URL up

cd gateway
$env:REDIS_ADDR = "localhost:6379"
go run .
```

Откат миграции: `goose -dir ./ledger/migrations postgres $env:DATABASE_URL down`.

Отдельно можно проверить только Ledger: `cd ledger; go run ./cmd` (подтянет `budgets.json`, если файл есть).

Подробнее: [ledger/README.md](ledger/README.md), [gateway/README.md](gateway/README.md).
