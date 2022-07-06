package main

import (
	"auth/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn" //_ means we only needs init the package
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webport = "80"

var count int64

type Config struct {
	DB     *sql.DB //pointer to DB
	Models data.Models
}

func main() {
	log.Println("Start on authentication service...")

	//Connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Cannot connect to Postgres!")
	}

	//Set up config
	appli := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webport),
		Handler: appli.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err) //exception handling
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres is not ready...")
			count++
		} else {
			log.Println("Successfully connected to Postgres!")
			return connection
		}

		if count > 10 {
			//stop the forever loop
			log.Println(err)
			return nil
		}

		log.Println("Back off for 2 seconds...")
		time.Sleep(time.Second * 2)
		continue
	}
}
