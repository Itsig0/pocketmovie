package database

import (
	"context"
	"database/sql"
	_ "embed"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	dburl = "./data/data.db"
)

//go:embed schema.sql
var ddl string

func Init() *Queries {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", dburl)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		log.Fatal(err)
	}

	return New(db)
}
