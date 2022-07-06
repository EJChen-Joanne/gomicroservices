package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"time"
)

//Methods that take this as a receiver type are available
//over RPC, as long as they are exported.
type RPCServer struct{}

//RPCPayload is the type with the data received from RPC Server
type RPCPayload struct {
	Name string
	Data string
}

//LogInfo: writes payload into mongoDB
func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {

	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error writing to mongo", err)
		return err
	}

	*resp = fmt.Sprintf("Processed payload via RPC: %s\n", payload.Name)
	return nil
}
