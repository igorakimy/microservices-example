package main

import (
	"context"
	"fmt"
	"log"
	"logger/data"
	"logger/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (ls *LogServer) WriteLog(
	ctx context.Context,
	req *logs.LogRequest,
) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	// write the log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	if err := ls.Models.LogEntry.Insert(logEntry); err != nil {
		resp := &logs.LogResponse{
			Result: "failed",
		}
		return resp, err
	}

	// return response
	return &logs.LogResponse{
		Result: "logged!",
	}, nil
}

func (s *Service) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	server := grpc.NewServer()

	logs.RegisterLogServiceServer(server, &LogServer{
		Models: s.Models,
	})

	log.Printf("gRPC server started on port %s", grpcPort)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
