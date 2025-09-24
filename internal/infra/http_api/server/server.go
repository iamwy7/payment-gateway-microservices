package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iamwy7/payment-gateway-management/internal/application/service"
	"github.com/iamwy7/payment-gateway-management/internal/infra/http_api/handlers"
	"github.com/iamwy7/payment-gateway-management/internal/infra/http_api/middleware"
)

type Server struct {
	router         *chi.Mux // Usaremos o Chi para poder habilitar o Middleware de forma seamless e para agrupar as rotas no mesmo recurso.
	server         *http.Server
	accountService *service.AccountService
	invoiceService *service.InvoiceService
	port           string
}

func NewServer(accountService *service.AccountService, invoiceService *service.InvoiceService,
	port string) *Server {
	return &Server{
		router:         chi.NewRouter(),
		accountService: accountService,
		invoiceService: invoiceService,
		port:           port,
	}
}

func (s *Server) SetupRoutes() {
	authMiddleware := middleware.NewAuthMiddleware(s.accountService)
	accountHandler := handlers.NewAccountHandler(s.accountService)
	invoiceHandler := handlers.NewInvoiceHandler(s.invoiceService)

	// Account routes
	s.router.Post("/accounts", accountHandler.Create)
	s.router.Get("/accounts", accountHandler.Get)

	// Invoices authenticated routes
	s.router.Group(func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)
		s.router.Post("/invoices", invoiceHandler.Create)
		s.router.Get("/invoices/{invoiceId}", invoiceHandler.GetById)
		s.router.Get("/invoices", invoiceHandler.ListByAccount)
	})
}

func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:    ":" + s.port,
		Handler: s.router,
	}
	return s.server.ListenAndServe()
}
