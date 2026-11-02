package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"ledger"
)

func main() {
	if err := ledger.InitDB(); err != nil {
		fmt.Println("ошибка подключения к БД:", err)
		os.Exit(1)
	}
	defer ledger.Close()

	if err := ledger.InitCache(); err != nil {
		fmt.Println("ошибка подключения к Redis:", err)
		os.Exit(1)
	}

	file, err := os.Open("budgets.json")
	if err != nil {
		log.Println("budgets.json не найден, пропускаем загрузку:", err)
	} else {
		defer file.Close()
		if err := ledger.LoadBudgets(file); err != nil {
			fmt.Println("ошибка загрузки бюджетов:", err)
			os.Exit(1)
		}
	}

	budgets, err := ledger.ListBudgets()
	if err != nil {
		fmt.Println("ошибка чтения бюджетов:", err)
		os.Exit(1)
	}
	fmt.Println("Бюджеты в БД:")
	for _, b := range budgets {
		fmt.Printf("  %s: лимит %.2f\n", b.Category, b.Limit)
	}

	list, err := ledger.ListTransactions()
	if err != nil {
		fmt.Println("ошибка чтения транзакций:", err)
		os.Exit(1)
	}
	fmt.Println("Транзакций в БД:", len(list))

	today := time.Now().Format("2006-01-02")
	monthStart := time.Now().Format("2006-01") + "-01"
	fmt.Println("Отчёт дважды (ожидаем cache miss, затем cache hit):")
	for i := 0; i < 2; i++ {
		report, err := ledger.GetReportSummary(context.Background(), monthStart, today)
		if err != nil {
			fmt.Println("ошибка отчёта:", err)
			os.Exit(1)
		}
		fmt.Printf("  всего %.2f, по категориям: %v\n", report.Total, report.ByCategory)
	}

	fmt.Println("Ledger готов (данные в PostgreSQL, отчёты кешируются в Redis)")
}
