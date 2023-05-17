package dao

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DAO struct {
	dbClient *sql.DB
}

func NewDAO() *DAO {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("PGSQL_HOST"), os.Getenv("PGSQL_PORT"),
		os.Getenv("PGSQL_USER"), os.Getenv("PGSQL_PASS"),
		os.Getenv("PGSQL_DB"))

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	dao := DAO{
		dbClient: db,
	}
	return &dao
}
