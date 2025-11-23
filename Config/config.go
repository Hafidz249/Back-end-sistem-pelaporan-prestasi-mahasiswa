package Config

import (
	"database/sql"
	"log"
	"os"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB


func ConnectDB() *sql.DB {
	dsn := os.Getenv("DB_DSN")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Koneksi database gagal: ", err)
	}

	DB = db

	return db
}

func GetDB() *sql.DB {
	return DB
}

func Ping() error {
    if DB == nil {
        return fmt.Errorf("database connection is not initialized")
    }
    return DB.Ping()
}