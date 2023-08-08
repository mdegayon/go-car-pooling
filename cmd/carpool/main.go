package main

import (
	"github.com/gin-gonic/gin"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/api"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/service"
)

func main() {
	carPool := service.NewCarpool()
	controller := api.NewController(carPool)

	gin.SetMode(gin.ReleaseMode)
	controller.Run()
}
