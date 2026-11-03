# Домашнее задание 7: базы данных (PostgreSQL + Redis)

В репозитории два Go-модуля:

- `ledger` — бизнес-логика, PostgreSQL (бюджеты/транзакции) и Redis (кеш отчётов);
- `gateway` — HTTP API (`/api/budgets`, `/api/transactions`, `/api/reports/summary`).

Gateway обращается к тем же HTTP-эндпоинтам, что и раньше, но Ledger хранит данные
в PostgreSQL, а отчёты кешируются в Redis на короткий TTL.

Требуется Go 1.22+, PostgreSQL, Redis, утилита [goose](https://github.com/pressly/goose).

## Быстрый старт

```powershell
# 1. Инфраструктура
docker run -d --name cashcraft-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=cashapp -p 5432:5432 postgres:16-alpine
docker run -d --name cashcraft-redis -p 6379:6379 redis:7-alpine

# 2. Миграции
go install github.com/pressly/goose/v3/cmd/goose@latest
$env:DATABASE_URL = "postgres://postgres:postgres@localhost:5432/cashapp?sslmode=disable"
goose -dir ./ledger/migrations postgres $env:DATABASE_URL up

# 3. Ledger (проверка подключения и демо-отчёта)
cd ledger
$env:REDIS_ADDR = "localhost:6379"
go run ./cmd

# 4. Gateway
cd ../gateway
go run .
```

Подробности: [ledger/README.md](ledger/README.md), [gateway/README.md](gateway/README.md).

## Что сделано

- схема `budgets` / `expenses`, миграции Goose (up/down);
- `sql.Open("pgx", …)` + `Ping`, пул `SetMaxOpenConns(10)` / `SetMaxIdleConns(5)`;
- CRUD в БД: `SetBudget` (UPSERT), `ListBudgets`, `AddTransaction` (проверка лимита → `budget exceeded`), `ListTransactions`;
- Redis: `GetReportSummary` с ключом `report:summary:<from>:<to>` и TTL 30с; опционально кеш `budgets:all`.
