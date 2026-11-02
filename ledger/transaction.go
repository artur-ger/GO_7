package ledger

import (
	"database/sql"
	"errors"
)

// Transaction — расход/доход по категории.
type Transaction struct {
	ID          int     `json:"id"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

func (t Transaction) Validate() error {
	if t.Amount == 0 {
		return errors.New("сумма транзакции не должна быть нулевой")
	}
	if t.Category == "" {
		return errors.New("у транзакции не указана категория")
	}
	if t.Date == "" {
		return errors.New("у транзакции не указана дата")
	}
	return nil
}

// ErrBudgetExceeded возвращается, если сумма расходов превысит лимит категории.
var ErrBudgetExceeded = errors.New("budget exceeded")

// AddTransaction проверяет лимит бюджета и вставляет запись в expenses.
func AddTransaction(tx Transaction) (Transaction, error) {
	if err := tx.Validate(); err != nil {
		return Transaction{}, err
	}

	var limit float64
	err := db.QueryRow(`SELECT limit_amount FROM budgets WHERE category = $1`, tx.Category).Scan(&limit)
	if err == sql.ErrNoRows {
		// Бюджета нет — по ТЗ допускается добавление без лимита.
		limit = 0
	} else if err != nil {
		return Transaction{}, err
	}

	if limit > 0 {
		var spent float64
		if err := db.QueryRow(
			`SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE category = $1`,
			tx.Category,
		).Scan(&spent); err != nil {
			return Transaction{}, err
		}
		if spent+tx.Amount > limit {
			return Transaction{}, ErrBudgetExceeded
		}
	}

	if err := db.QueryRow(
		`INSERT INTO expenses(amount, category, description, date)
		 VALUES($1, $2, $3, $4) RETURNING id`,
		tx.Amount, tx.Category, tx.Description, tx.Date,
	).Scan(&tx.ID); err != nil {
		return Transaction{}, err
	}
	return tx, nil
}

// ListTransactions читает все расходы из БД.
func ListTransactions() ([]Transaction, error) {
	rows, err := db.Query(
		`SELECT id, amount, category, COALESCE(description, ''), date::text
		 FROM expenses ORDER BY date DESC, id DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]Transaction, 0)
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.ID, &t.Amount, &t.Category, &t.Description, &t.Date); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
