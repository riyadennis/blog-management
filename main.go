package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/riyadennis/blog-management/commands"
	"github.com/riyadennis/blog-management/db"
	"log"
	"net/http"
)

var (
	conn db.EventStore
	err  error
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	conn, err = db.NewConn()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	command := commands.NewCommand()
	http.HandleFunc("/api/v1/CreateArticle", func(w http.ResponseWriter, r *http.Request) {
		command.CreateArticle(conn, w, r)
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
