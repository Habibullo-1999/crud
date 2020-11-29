package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"

)

func main() {
	dsn := "postgres://app:pass@localhost:5432/db"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Println(err)
		os.Exit(1)
		return
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Println(err)
		}
	}()

	ctx := context.Background()
	_, err = db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS customers (
		id			BIGSERIAL 	PRIMARY KEY,
		name 		TEXT 		NOT NULL,
		phone 		TEXT 		NOT NULL UNIQUE,
		active      BOOLEAN 	NOT NULL DEFAULT TRUE,
		created 	TIMESTAMP	NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
`)
	if err != nil {
		log.Print(err)
		return
	}

	result, err := db.ExecContext(ctx, `
	INSET INTO customers(name, phone) VALUES ('Vasya','+992000000001') ON CONFLICT DO NOTHING;
	`)
	if err != nil {
		log.Print(err)
		return
	}
	log.Print(result.RowsAffected())
	log.Print(result.LastInsertId())

}
