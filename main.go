package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/joho/godotenv"
	"github.com/riyadennis/blog-management/api/v1/CreateArticle"
	"github.com/riyadennis/blog-management/pkg/api/eventsource"
)

var (
	eventStore eventsource.EventStore
	err        error
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	eventStore, err = eventsource.NewConn()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/api/v1/CreateArticle", CreateArticle.NewCommand(eventStore).ServeHTTP)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
