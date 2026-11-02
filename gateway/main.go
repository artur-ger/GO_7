package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gateway/internal/api"
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           api.Logging(api.NewRouter()),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Gateway запущен на :%s (данные Ledger в PostgreSQL, отчёты в Redis)", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
