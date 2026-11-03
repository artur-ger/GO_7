# Ledger

Таблицы:

- `budgets` — лимит по категории;
- `expenses` — транзакции.

Отчёт `GetReportSummary` пишется в Redis на 30 секунд (ключ `report:summary:<from>:<to>`).
Список бюджетов тоже можно закешировать (`budgets:all`, короткий TTL); после `SetBudget` ключ сбрасывается.

## Миграции

```powershell
docker run -d --name cashcraft-postgres `
  -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=cashapp `
  -p 5432:5432 postgres:16-alpine

docker run -d --name cashcraft-redis -p 6379:6379 redis:7-alpine

go install github.com/pressly/goose/v3/cmd/goose@latest
$env:DATABASE_URL = "postgres://postgres:postgres@localhost:5432/cashapp?sslmode=disable"
goose -dir ./migrations postgres $env:DATABASE_URL up
```

Откат: `goose -dir ./migrations postgres $env:DATABASE_URL down`.

## Переменные

| Переменная | Зачем |
|---|---|
| `DATABASE_URL` | строка подключения (главный вариант) |
| `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_NAME` | если `DATABASE_URL` не задан (по умолчанию `localhost:5432`, `postgres`/`postgres`, БД `cashapp`) |
| `REDIS_ADDR` | по умолчанию `localhost:6379` |
| `REDIS_DB` | по умолчанию `0` |
| `REDIS_PASSWORD` | если нужен пароль |

## Запуск

```powershell
$env:DATABASE_URL = "postgres://postgres:postgres@localhost:5432/cashapp?sslmode=disable"
$env:REDIS_ADDR = "localhost:6379"
go run ./cmd
```

При старте: `Ping` к Postgres и Redis, пул `MaxOpenConns=10`, `MaxIdleConns=5`.
Если есть `budgets.json`, он загрузится в БД. Отчёт вызывается два раза: в логе сначала `cache miss`, потом `cache hit`.
