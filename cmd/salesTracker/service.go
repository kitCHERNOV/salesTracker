package main

import (
	"database/sql"
	"fmt"
	"net/http"
	config "salesTracker/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

type SalesService struct {
	Database *sql.DB
}

func MustOpenDataBaseConnection(connStr, connDriver string) *sql.DB {
	db, err := sql.Open(connDriver, connStr)
	if err != nil {
		panic(err)
	}

	return db
}

func NewSalesService() SalesService {
	const op = "NewSalesService"

	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	cfg, err := config.MustLoad()
	if err != nil {
		panic(fmt.Errorf("%s,%v", op, err))
	}

	return SalesService{
		Database: MustOpenDataBaseConnection(cfg.Database.DSN(), cfg.Database.Driver),
	}
}

func (s *SalesService) Run(server string) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/items", createItem)
	r.Get("/items", getItems)
	r.Put("/items/{id}", updateItem)
	r.Delete("/items/{id}", deleteItem)

	if err := http.ListenAndServe(server, r); err != nil {
		panic(err)
	}
}
