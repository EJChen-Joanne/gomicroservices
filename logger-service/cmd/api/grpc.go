package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"logger-service/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer //ensure backwards compatibility
	Models                             data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	//write the log

	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{
			Result: "failed to log in",
		}
		return res, err
	}

	//return response
	res := &logs.LogResponse{
		Result: "logged in",
	}
	return res, nil
}

//listens to the gRPC server
func (appli *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen to gRPC: %v\n", err)
	}
	defer lis.Close() //?//

	srv := grpc.NewServer()

	logs.RegisterLogServiceServer(srv, &LogServer{Models: appli.Models})

	log.Printf("gRPC Server started on %s\n", gRpcPort)

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to listen to gRPC: %v\n", err)
	}
}
