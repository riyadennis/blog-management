package main

import (
	"github.com/riyadennis/blog-management/api/v1/article"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/joho/godotenv"
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
	err = http.ListenAndServe(":8080", article.NewArticle(eventStore))
	if err != nil {
		log.Fatal(err)
	}
}
