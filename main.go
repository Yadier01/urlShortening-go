package main

import (
	"database/sql"
	"log"

	"github.com/Yadier01/urlshort/internal"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	server, err := internal.NewServer(db)

	log.Fatal(server.Conn.ListenAndServe())

}
