package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80" //fail with port 8080

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	// try to connect to rabbitmq for retrieving requests from queue
	rabbitConn, err := connectToRabbitmq()
	if err != nil {
		log.Fatalln(err)
	}
	defer rabbitConn.Close()

	appli := new(Config)
	appli.Rabbit = rabbitConn

	log.Printf("Start on broker service on port %s\n", webPort)

	//define http server
	server := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		//add a route
		Handler: appli.routes(),
	}

	//start server
	err = server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// listen to the requests from the rabbitmq
func connectToRabbitmq() (*amqp.Connection, error) {
	var counts int64 // break down loop until error connection 5 times
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ is not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ")
			connection = conn
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("Backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
