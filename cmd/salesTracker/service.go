package main

import (
	"database/sql"
	"fmt"
	"net/http"
	config "salesTracker/internal/config"
	"salesTracker/internal/handlers"
	"salesTracker/internal/handlers/analytics"
	postgresql "salesTracker/internal/storage/postgresql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func MustOpenDataBaseConnection(connStr, connDriver string) *sql.DB {
	db, err := sql.Open(connDriver, connStr)
	if err != nil {
		panic(err)
	}

	return db
}

func NewSalesService() *postgresql.Storage {
	const op = "NewSalesService"

	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	cfg, err := config.MustLoad()
	if err != nil {
		panic(fmt.Errorf("%s: %w", op, err))
	}

	return &postgresql.Storage{
		DB: MustOpenDataBaseConnection(cfg.Database.DSN(), cfg.Database.Driver),
	}
}

// setupRoutes - настраивает все роуты приложения
func setupRoutes(r *chi.Mux, storage *postgresql.Storage) {
	// ====================================================================
	// API v1 - Основные CRUD операции
	// ====================================================================
	r.Route("/api/v1", func(r chi.Router) {
		// CATEGORIES - Категории товаров
		r.Route("/categories", func(r chi.Router) {
			r.Post("/", handlers.CreateCategory(storage))
			r.Get("/", handlers.ListCategories(storage))
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handlers.GetCategory(storage))
				r.Put("/", handlers.UpdateCategory(storage))
				r.Delete("/", handlers.DeleteCategory(storage))
				// Товары категории
				r.Get("/products", handlers.ListProductsByCategory(storage))
			})
		})

		// PRODUCTS - Товары
		r.Route("/products", func(r chi.Router) {
			r.Post("/", handlers.CreateProduct(storage))
			r.Get("/", handlers.ListProducts(storage))
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handlers.GetProduct(storage))
				r.Put("/", handlers.UpdateProduct(storage))
				r.Delete("/", handlers.DeleteProduct(storage))
			})
		})

		// CUSTOMERS - Покупатели
		r.Route("/customers", func(r chi.Router) {
			r.Post("/", handlers.CreateCustomer(storage))
			r.Get("/", handlers.ListCustomers(storage))
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handlers.GetCustomer(storage))
				r.Put("/", handlers.UpdateCustomer(storage))
				r.Delete("/", handlers.DeleteCustomer(storage))
				// Заказы покупателя
				r.Get("/orders", handlers.ListOrdersByCustomer(storage))
			})
		})

		// ORDERS - Заказы
		r.Route("/orders", func(r chi.Router) {
			r.Post("/", handlers.CreateOrder(storage))
			r.Get("/", handlers.ListOrders(storage))
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handlers.GetOrder(storage))
				r.Put("/", handlers.UpdateOrder(storage))
				r.Delete("/", handlers.DeleteOrder(storage))
				// Позиции заказа
				r.Get("/items", handlers.ListOrderItems(storage))
			})
		})

		// ORDER ITEMS - Позиции в заказах
		r.Route("/order-items", func(r chi.Router) {
			r.Post("/", handlers.CreateOrderItem(storage))
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handlers.GetOrderItem(storage))
				r.Put("/", handlers.UpdateOrderItem(storage))
				r.Delete("/", handlers.DeleteOrderItem(storage))
			})
		})
	})

	// ====================================================================
	// ANALYTICS - Аналитика
	// ====================================================================
	r.Route("/analytics", func(r chi.Router) {
		r.Get("/revenue", analytics.TotalRevenueByPeriod(storage))
		r.Get("/daily-orders", analytics.OrdersPerDay(storage))
		r.Get("/average-check", analytics.AverageCheckByPeriod(storage))
		r.Get("/orders-median", analytics.OrdersMedian(storage))
		r.Get("/customer-median", analytics.CustomerSpendingMedian(storage))
		r.Get("/orders-percentile", analytics.OrdersPercentile(storage))
		r.Get("/customer-percentile", analytics.CustomerSpendingPercentile(storage))
		r.Get("/sales-report", analytics.GenerateSalesReport(storage))
	})
}

// Run запускает HTTP сервер с всеми роутами
func Run(server string, storage *postgresql.Storage) {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Настраиваем роуты
	setupRoutes(r, storage)

	if err := http.ListenAndServe(server, r); err != nil {
		panic(err)
	}
}
