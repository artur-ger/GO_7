# Gateway

Эндпоинты те же, что раньше, но Ledger хранит данные в PostgreSQL.
`GET /api/reports/summary` кешируется в Redis (TTL 30с).

Сначала поднимите PostgreSQL/Redis и примените миграции — см. [../ledger/README.md](../ledger/README.md).

## Запуск

```powershell
$env:DATABASE_URL = "postgres://postgres:postgres@localhost:5432/cashapp?sslmode=disable"
$env:REDIS_ADDR = "localhost:6379"
go run .
```

## Проверка через cURL

```bash
curl -X POST http://localhost:8080/api/budgets ^
  -H "Content-Type: application/json" ^
  -d "{\"category\":\"еда\",\"limit\":1000}"
```

```bash
curl -X POST http://localhost:8080/api/transactions ^
  -H "Content-Type: application/json" ^
  -d "{\"amount\":450,\"category\":\"еда\",\"description\":\"обед\",\"date\":\"2026-11-01\"}"
```

Сумма 700 при лимите 1000 и уже потраченных 450 → **409**:

```bash
curl -X POST http://localhost:8080/api/transactions ^
  -H "Content-Type: application/json" ^
  -d "{\"amount\":700,\"category\":\"еда\",\"description\":\"ужин\",\"date\":\"2026-11-02\"}"
```

```bash
curl http://localhost:8080/api/transactions
curl http://localhost:8080/api/budgets
```

Отчёт дважды подряд — в логах `cache miss`, затем `cache hit`:

```bash
curl "http://localhost:8080/api/reports/summary?from=2026-11-01&to=2026-11-30"
curl "http://localhost:8080/api/reports/summary?from=2026-11-01&to=2026-11-30"
```

После перезапуска Gateway данные остаются в PostgreSQL.
