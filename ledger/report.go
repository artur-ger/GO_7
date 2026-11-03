package ledger

import (
	"context"
	"encoding/json"
	"log"
	"time"
)

// ReportSummary — суммы расходов по категориям за период.
type ReportSummary struct {
	From       string             `json:"from"`
	To         string             `json:"to"`
	ByCategory map[string]float64 `json:"by_category"`
	Total      float64            `json:"total"`
}

// GetReportSummary считает отчёт через БД с кешированием в Redis (TTL 30с).
func GetReportSummary(ctx context.Context, from, to string) (ReportSummary, error) {
	key := "report:summary:" + from + ":" + to

	if rdb != nil {
		if cached, err := rdb.Get(ctx, key).Bytes(); err == nil {
			var report ReportSummary
			if json.Unmarshal(cached, &report) == nil {
				log.Println("cache hit:", key)
				return report, nil
			}
		}
		log.Println("cache miss:", key)
	}

	rows, err := db.QueryContext(ctx,
		`SELECT category, SUM(amount)
		 FROM expenses
		 WHERE date BETWEEN $1 AND $2
		 GROUP BY category`,
		from, to,
	)
	if err != nil {
		return ReportSummary{}, err
	}
	defer rows.Close()

	report := ReportSummary{
		From:       from,
		To:         to,
		ByCategory: map[string]float64{},
	}
	for rows.Next() {
		var category string
		var sum float64
		if err := rows.Scan(&category, &sum); err != nil {
			return ReportSummary{}, err
		}
		report.ByCategory[category] = sum
		report.Total += sum
	}
	if err := rows.Err(); err != nil {
		return ReportSummary{}, err
	}

	if rdb != nil {
		if data, err := json.Marshal(report); err == nil {
			_ = rdb.Set(ctx, key, data, 30*time.Second).Err()
		}
	}
	return report, nil
}
