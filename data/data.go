package data

// defines configuration for connecting to external mysql db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// Initiates connection to postgres database
func Connect() *sql.DB {
	var port = os.Getenv("POSTGRES_PORT")
	var host = os.Getenv("POSTGRES_HOST")
	var user = os.Getenv("POSTGRES_USER")
	var pass = os.Getenv("POSTGRES_PASSWORD")
	var name = os.Getenv("POSTGRES_NAME")

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, name)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Print("[ERROR] Problem connecting to database")
		panic(err)
	}

	return db
}

func executeQuery(query string) *sql.Rows {
	db := Connect()
	defer db.Close()

	res, err := db.Query(query)
	if err != nil {
		log.Print(err.Error())
	}
	return res
}
