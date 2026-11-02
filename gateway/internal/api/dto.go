package api

import "ledger"

type CreateTransactionRequest struct {
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

type TransactionResponse struct {
	ID          int     `json:"id"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

type CreateBudgetRequest struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
}

type BudgetResponse struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
}

type ReportResponse struct {
	From       string             `json:"from"`
	To         string             `json:"to"`
	ByCategory map[string]float64 `json:"by_category"`
	Total      float64            `json:"total"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func toTransaction(req CreateTransactionRequest) ledger.Transaction {
	return ledger.Transaction{
		Amount:      req.Amount,
		Category:    req.Category,
		Description: req.Description,
		Date:        req.Date,
	}
}

func toTransactionResponse(tx ledger.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:          tx.ID,
		Amount:      tx.Amount,
		Category:    tx.Category,
		Description: tx.Description,
		Date:        tx.Date,
	}
}

func toBudget(req CreateBudgetRequest) ledger.Budget {
	return ledger.Budget{
		Category: req.Category,
		Limit:    req.Limit,
	}
}

func toBudgetResponse(b ledger.Budget) BudgetResponse {
	return BudgetResponse{
		Category: b.Category,
		Limit:    b.Limit,
	}
}

func toReportResponse(r ledger.ReportSummary) ReportResponse {
	return ReportResponse{
		From:       r.From,
		To:         r.To,
		ByCategory: r.ByCategory,
		Total:      r.Total,
	}
}
