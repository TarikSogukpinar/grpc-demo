package main

import (
	"log"
	"net"

	"grpc-demo/api"
	pb "grpc-demo/proto"
	"grpc-demo/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	go startGrpcServer()

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	grpcClient := pb.NewWeatherServiceClient(conn)

	router := api.SetupRouter(grpcClient)

	log.Println("Starting REST API server on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func startGrpcServer() {

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterWeatherServiceServer(grpcServer, &server.WeatherServer{})

	log.Println("gRPC Server running on port 50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
