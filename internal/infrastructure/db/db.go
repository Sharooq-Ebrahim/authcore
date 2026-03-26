package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectDB(dbURL string) (*sql.DB, error) {

	db, err := sql.Open("pgx", dbURL)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {

		log.Println("Database connection failed", err)
		return nil, err
	}

	log.Println("Database connected successfully")

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	return db, nil

}
