package api

import (
	"context"
	"net/http"
	"time"

	pb "grpc-demo/proto"

	"github.com/gin-gonic/gin"
)

type WeatherHandler struct {
	GrpcClient pb.WeatherServiceClient
}

func (h *WeatherHandler) GetWeather(c *gin.Context) {
	city := c.Param("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "şehir parametresi gerekli"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	response, err := h.GrpcClient.GetWeather(ctx, &pb.WeatherRequest{City: city})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"city":        response.City,
		"temperature": response.Temperature,
		"conditions":  response.Conditions,
		"humidity":    response.Humidity,
		"wind_speed":  response.WindSpeed,
		"updated_at":  response.UpdatedAt,
	})
}

func SetupRouter(grpcClient pb.WeatherServiceClient) *gin.Engine {
	router := gin.Default()

	handler := &WeatherHandler{GrpcClient: grpcClient}

	router.GET("/weather/:city", handler.GetWeather)

	return router
}
