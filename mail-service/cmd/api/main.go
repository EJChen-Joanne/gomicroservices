package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Mailer Mail
}

const webPort = "80"

func main() {
	appli := Config{
		Mailer: createMail(),
	}

	log.Println("Start on mail service on port", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: appli.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	var m Mail
	m.Domain = os.Getenv("MAIL_DOMAIN")
	m.Host = os.Getenv("MAIL_HOST")
	m.Port = port
	m.Username = os.Getenv("MAIL_USERNAME")
	m.Password = os.Getenv("MAIL_PASSWORD")
	m.Encryption = os.Getenv("MAIL_ENCRYPTION")
	m.FromName = os.Getenv("FROM_NAME")
	m.FromAddress = os.Getenv("FROM_ADDRESS")

	return m
}
