package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/riyadennis/blog-management/commands"
	"github.com/riyadennis/blog-management/handlers"
	"log"
	"net/http"
)

var (
	conn commands.EventStore
	err  error
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	conn, err = commands.NewConn()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/api/v1/article", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateEvent(conn, w, r)
	})
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
