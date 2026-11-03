# Ledger

Бюджеты хранятся в таблице `budgets`, транзакции — в `expenses`.
Отчёты по суммам за период кешируются в Redis на 30 секунд.

## PostgreSQL, Redis и миграции

```bash
docker run -d --name cashcraft-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=cashapp -p 5432:5432 postgres:16-alpine
docker run -d --name cashcraft-redis -p 6379:6379 redis:7-alpine

go install github.com/pressly/goose/v3/cmd/goose@latest

# PowerShell
$env:DATABASE_URL = "postgres://postgres:postgres@localhost:5432/cashapp?sslmode=disable"
goose -dir ./migrations postgres $env:DATABASE_URL up

# Bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/cashapp?sslmode=disable"
goose -dir ./migrations postgres "$DATABASE_URL" up
```

Откат одной миграции:

```bash
goose -dir ./migrations postgres "$DATABASE_URL" down
```

## Запуск Ledger

Подключение к БД — из `DATABASE_URL` (или `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_NAME`).
Redis — из `REDIS_ADDR`, `REDIS_DB`, `REDIS_PASSWORD`.

```powershell
$env:DATABASE_URL = "postgres://postgres:postgres@localhost:5432/cashapp?sslmode=disable"
$env:REDIS_ADDR = "localhost:6379"
go run ./cmd
```

Программа подключается к PostgreSQL и Redis, при наличии загружает `budgets.json`,
читает данные из БД и дважды строит отчёт (в логах будет `cache miss`, затем `cache hit`).
