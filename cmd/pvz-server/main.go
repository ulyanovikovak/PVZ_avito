package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/ulyanovikovak/PVZ_avito/internal/handler"
	"github.com/ulyanovikovak/PVZ_avito/internal/middleware"
	"github.com/ulyanovikovak/PVZ_avito/internal/storage"

	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("DB not reachable: %v", err)
	}

	pvzStore := storage.NewPVZStorage(db)
	receptionStore := storage.NewReceptionStorage(db)
	productStore := storage.NewProductStorage(db)

	r := mux.NewRouter()

	r.HandleFunc("/dummyLogin", handler.DummyLoginHandler).Methods("POST")

	api := r.PathPrefix("/").Subrouter()
	api.Use(middleware.JWTAuthMiddleware)

	api.Handle("/pvz", handler.CreatePVZHandler(pvzStore)).Methods("POST")
	api.Handle("/pvz", handler.GetPVZsHandler(pvzStore)).Methods("GET")

	api.Handle("/receptions", handler.CreateReceptionHandler(receptionStore)).Methods("POST")
	api.Handle("/pvz/{pvzId}/close_last_reception", handler.CloseLastReceptionHandler(receptionStore)).Methods("POST")

	api.Handle("/products", handler.AddProductHandler(productStore)).Methods("POST")
	api.Handle("/pvz/{pvzId}/delete_last_product", handler.DeleteLastProductHandler(productStore)).Methods("POST")

	// Старт сервера
	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
