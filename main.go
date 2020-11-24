package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/riyadennis/blog-management/api"
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

	eventsource.Set(eventStore)
}

func main() {
	err = http.ListenAndServe(":8080", api.NewAPIv1())
	if err != nil {
		log.Fatal(err)
	}
}
