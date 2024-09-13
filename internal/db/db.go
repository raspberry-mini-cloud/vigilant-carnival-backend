package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)
var Db *sql.DB

func InitDB() {
	postgres_password := os.Getenv("POSTGRES_PASSWORD")

	connStr := fmt.Sprintf("postgres://postgres:%s@10.40.125.129:5432/coolkeeper_data?sslmode=disable", postgres_password)
    Db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Failed to init to the database:", err)
    }

    // Verify the connection is valid
    err = Db.Ping()
    if err != nil {
        log.Fatal("Failed to connect to the database:", err)
    }

}