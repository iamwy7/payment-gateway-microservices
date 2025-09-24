package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/iamwy7/payment-gateway-management/internal/application/service"
	"github.com/iamwy7/payment-gateway-management/internal/infra/http_api/server"
	"github.com/iamwy7/payment-gateway-management/internal/infra/repository"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Postgres Driver
)

// getEnv retorna variável de ambiente ou valor padrão se não definida
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Carrega variáveis de ambiente do arquivo .env
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Configura conexão com DB usando variáveis de ambiente
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "db"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "gateway"),
		getEnv("DB_SSL_MODE", "disable"),
	)

	// Inicializa conexão com o banco
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	defer db.Close()

	// Inicializa camadas da aplicação (repository -> service -> server)
	accountRepository := repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepository)

	invoiceRepository := repository.NewInvoiceRepository(db)
	invoiceService := service.NewInvoiceService(invoiceRepository, accountService)

	// Configura e inicia o servidor HTTP
	port := getEnv("HTTP_PORT", "8080")
	srv := server.NewServer(accountService, invoiceService, port)
	srv.SetupRoutes()

	if err := srv.Start(); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
