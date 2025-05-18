package server

import (
	"context"
	"log"
	"math/rand"
	"time"

	pb "grpc-demo/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WeatherServer struct {
	pb.UnimplementedWeatherServiceServer
}

func (s *WeatherServer) GetWeather(ctx context.Context, req *pb.WeatherRequest) (*pb.WeatherResponse, error) {
	log.Printf("Received request for weather in: %s", req.City)

	if req.City == "" {
		return nil, status.Error(codes.InvalidArgument, "city cannot be empty")
	}

	// Not real API call, just simulating a delay
	response := generateWeatherData(req.City)

	return response, nil
}

func (s *WeatherServer) GetWeatherStream(req *pb.WeatherRequest, stream pb.WeatherService_GetWeatherStreamServer) error {
	log.Printf("Started streaming weather for: %s", req.City)

	if req.City == "" {
		return status.Error(codes.InvalidArgument, "city cannot be empty")
	}

	for i := 0; i < 5; i++ {

		response := generateWeatherData(req.City)

		if err := stream.Send(response); err != nil {
			return err
		}

		time.Sleep(2 * time.Second)
	}

	return nil
}

func generateWeatherData(city string) *pb.WeatherResponse {
	conditions := []string{"Sunny", "Cloudy", "Rainy", "Stormy", "Snow"}
	randomCondition := conditions[rand.Intn(len(conditions))]

	now := time.Now().Format(time.RFC3339)

	return &pb.WeatherResponse{
		City:        city,
		Temperature: 10.0 + rand.Float32()*25.0,
		Conditions:  randomCondition,
		Humidity:    30.0 + rand.Float32()*70.0,
		WindSpeed:   rand.Float32() * 30.0,
		UpdatedAt:   now,
	}
}
