package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/riyadennis/blog-management/commands"
	"os"
)


// Article holds structure of the blob in Article
type Article struct {
	CreatedBy string
	Name string
	Heading string
	Body string
}


func Store(ctx context.Context, e *commands.CreateCommand) error{
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("MYSQL_USERNAME"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)

	conn, err := sql.Open("mysql", dsn)
	if err != nil{
		return err
	}

	query, err := conn.Prepare("INSERT INTO events_store(id,version,state,data) values(?,?,?,?)")
	if err != nil{
		return err
	}

	_, err = query.Exec(e.ID,e.EventVersion,e.State,e.Article)
	if err != nil{
		return err
	}
	defer conn.Close()

	return nil
}

