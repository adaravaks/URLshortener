package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/adaravaks/URLshortener/internal/handler"
	"github.com/adaravaks/URLshortener/internal/repository"
	_ "github.com/lib/pq"
)

func main() {
	dsn := "postgres://urlshortener:urlshortener@localhost:5433/urlshortener?sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to open db connection: ", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("failed to close db connection: ", err)
		}
	}(db)

	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping db: ", err)
	}
	log.Println("connected to db successfully")

	repo := repository.NewPostgresLinkRepository(db)
	h := handler.NewHandler(repo)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Health checked")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("POST /links", h.CreateLink)
	mux.HandleFunc("GET /links/{code}/stats", h.Stats)
	mux.HandleFunc("GET /{code}", h.Redirect)
	wrapped := handler.LoggingMiddleware(handler.RecoveryMiddleware(mux))

	log.Println("Server started on :8080")
	err = http.ListenAndServe(":8080", wrapped)
	if err != nil {
		log.Fatal(err)
	}
}
