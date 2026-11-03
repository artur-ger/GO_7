package ledger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"time"
)

const budgetsCacheKey = "budgets:all"

// Budget — лимит расходов по категории.
type Budget struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
}

func (b Budget) Validate() error {
	if b.Category == "" {
		return errors.New("у бюджета не указана категория")
	}
	if b.Limit <= 0 {
		return errors.New("лимит бюджета должен быть больше нуля")
	}
	return nil
}

// SetBudget выполняет UPSERT бюджета категории и инвалидирует кеш списка.
func SetBudget(b Budget) error {
	if err := b.Validate(); err != nil {
		return err
	}

	_, err := db.Exec(
		`INSERT INTO budgets(category, limit_amount) VALUES($1, $2)
		 ON CONFLICT(category) DO UPDATE SET limit_amount = EXCLUDED.limit_amount`,
		b.Category, b.Limit,
	)
	if err != nil {
		return err
	}

	if rdb != nil {
		if delErr := rdb.Del(context.Background(), budgetsCacheKey).Err(); delErr != nil {
			log.Println("не удалось инвалидировать кеш budgets:all:", delErr)
		}
	}
	return nil
}

// ListBudgets читает все бюджеты из БД (с коротким кешем в Redis).
func ListBudgets() ([]Budget, error) {
	ctx := context.Background()
	if rdb != nil {
		if cached, err := rdb.Get(ctx, budgetsCacheKey).Bytes(); err == nil {
			var list []Budget
			if json.Unmarshal(cached, &list) == nil {
				log.Println("cache hit:", budgetsCacheKey)
				return list, nil
			}
		}
		log.Println("cache miss:", budgetsCacheKey)
	}

	rows, err := db.Query(`SELECT category, limit_amount FROM budgets ORDER BY category`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]Budget, 0)
	for rows.Next() {
		var b Budget
		if err := rows.Scan(&b.Category, &b.Limit); err != nil {
			return nil, err
		}
		result = append(result, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if rdb != nil {
		if data, err := json.Marshal(result); err == nil {
			_ = rdb.Set(ctx, budgetsCacheKey, data, 15*time.Second).Err()
		}
	}
	return result, nil
}

// LoadBudgets читает JSON-массив бюджетов и сохраняет каждый через SetBudget.
func LoadBudgets(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("не удалось прочитать бюджеты: %w", err)
	}

	var list []Budget
	if err := json.Unmarshal(data, &list); err != nil {
		return fmt.Errorf("не удалось разобрать JSON с бюджетами: %w", err)
	}

	for _, b := range list {
		if err := SetBudget(b); err != nil {
			return err
		}
	}
	return nil
}
