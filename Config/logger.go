package config

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
    connStr := "host=localhost port=5432 user=projectuas password=12345678 dbname=postgres sslmode=disable"
	var err error
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Gagal koneksi ke database:", err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal("Gagal ping database:", err)
    }

    return db
}

func LoggerMiddleware(c *fiber.Ctx) error {
	fmt.Println("Request:", c.Method(), c.Path())
	return c.Next()
}