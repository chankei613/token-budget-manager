// cmd/tbmserve はToken Budget Manager APIをlocalhostで提供する単体サーバー。
//
//	go run ./cmd/tbmserve -addr :8424 -db token-budget-manager.db
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/chankei613/token-budget-manager/internal/api"
	"github.com/chankei613/token-budget-manager/internal/db"
)

func main() {
	addr := flag.String("addr", ":8424", "待ち受けアドレス")
	dbPath := flag.String("db", "token-budget-manager.db", "SQLiteファイル")
	flag.Parse()

	conn, err := db.Init(*dbPath)
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}

	router := api.NewRouter(conn)
	log.Printf("token-budget-manager backend listening on %s", *addr)
	if err := http.ListenAndServe(*addr, router); err != nil {
		log.Fatal(err)
	}
}
