package main

import (
	"log"
	"net/http"

	"github.com/chankei613/token-budget-manager/internal/api"
	"github.com/chankei613/token-budget-manager/internal/db"
)

func main() {
	conn, err := db.Init("token-budget-manager.db")
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}

	router := api.NewRouter(conn)
	log.Println("token-budget-manager backend listening on :8424")
	if err := http.ListenAndServe(":8424", router); err != nil {
		log.Fatal(err)
	}
}
