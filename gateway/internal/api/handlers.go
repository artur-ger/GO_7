package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"ledger"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", ping)
	mux.HandleFunc("POST /api/transactions", createTransaction)
	mux.HandleFunc("GET /api/transactions", listTransactions)
	mux.HandleFunc("POST /api/budgets", createBudget)
	mux.HandleFunc("GET /api/budgets", listBudgets)
	mux.HandleFunc("GET /api/reports/summary", reportSummary)
	return mux
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("pong"))
}

func createTransaction(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "некорректный JSON в теле")
		return
	}

	tx := toTransaction(req)
	if err := tx.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	created, err := ledger.AddTransaction(tx)
	if err == ledger.ErrBudgetExceeded {
		writeError(w, http.StatusConflict, "budget exceeded")
		return
	}
	if err != nil {
		log.Println("ошибка добавления транзакции:", err)
		writeError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}

	writeJSON(w, http.StatusCreated, toTransactionResponse(created))
}

func listTransactions(w http.ResponseWriter, r *http.Request) {
	list, err := ledger.ListTransactions()
	if err != nil {
		log.Println("ошибка чтения транзакций:", err)
		writeError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}

	result := make([]TransactionResponse, 0, len(list))
	for _, tx := range list {
		result = append(result, toTransactionResponse(tx))
	}
	writeJSON(w, http.StatusOK, result)
}

func createBudget(w http.ResponseWriter, r *http.Request) {
	var req CreateBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "некорректный JSON в теле")
		return
	}

	b := toBudget(req)
	if err := b.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := ledger.SetBudget(b); err != nil {
		log.Println("ошибка сохранения бюджета:", err)
		writeError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}

	writeJSON(w, http.StatusCreated, toBudgetResponse(b))
}

func listBudgets(w http.ResponseWriter, r *http.Request) {
	list, err := ledger.ListBudgets()
	if err != nil {
		log.Println("ошибка чтения бюджетов:", err)
		writeError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}

	result := make([]BudgetResponse, 0, len(list))
	for _, b := range list {
		result = append(result, toBudgetResponse(b))
	}
	writeJSON(w, http.StatusOK, result)
}

func reportSummary(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	if _, errFrom := time.Parse("2006-01-02", from); errFrom != nil {
		writeError(w, http.StatusBadRequest, "нужны параметры from и to в формате YYYY-MM-DD")
		return
	}
	if _, errTo := time.Parse("2006-01-02", to); errTo != nil {
		writeError(w, http.StatusBadRequest, "нужны параметры from и to в формате YYYY-MM-DD")
		return
	}

	report, err := ledger.GetReportSummary(r.Context(), from, to)
	if err != nil {
		log.Println("ошибка построения отчёта:", err)
		writeError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}
	writeJSON(w, http.StatusOK, toReportResponse(report))
}
