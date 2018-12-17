package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/service"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/service/model"
)

type Controller struct {
	service *service.CarPool
	engine  *gin.Engine
}

func NewController(service *service.CarPool) *Controller {
	c := &Controller{
		service: service,
		engine:  gin.New(),
	}
	c.engine.GET("/status", c.getStatus)
	c.engine.Any("/cars", c.putCars)
	c.engine.Any("/journey", c.postJourney)
	c.engine.Any("/dropoff", c.postDropoff)
	c.engine.Any("/locate", c.postLocate)
	return c
}

func (c *Controller) Run() {
	c.engine.Run("0.0.0.0:8080")
}

func (c *Controller) getStatus(ctx *gin.Context) {
	ctx.String(http.StatusOK, `{"status":"ok"}`)
}

func (c *Controller) putCars(ctx *gin.Context) {
	if ctx.Request.Method != "PUT" {
		ctx.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}
	if ctx.ContentType() != "application/json" {
		ctx.AbortWithStatus(http.StatusUnsupportedMediaType)
		return
	}

	var cars []*model.Car
	if err := ctx.BindJSON(&cars); err != nil {
		return
	}
	for _, car := range cars {
		car.AvailableSeats = car.MaxSeats
	}

	if err := c.service.ResetCars(cars); err != nil {
		switch err {
		case service.ErrDuplicatedID:
			ctx.Status(http.StatusBadRequest)
		default:
			ctx.Status(http.StatusInternalServerError)
		}
		return
	}
	ctx.Status(http.StatusOK)
}

func (c *Controller) postJourney(ctx *gin.Context) {
	if ctx.Request.Method != "POST" {
		ctx.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}
	if ctx.ContentType() != "application/json" {
		ctx.AbortWithStatus(http.StatusUnsupportedMediaType)
		return
	}

	var journey model.Journey
	if err := ctx.BindJSON(&journey); err != nil {
		return
	}
	if err := c.service.NewJourney(&journey); err != nil {
		switch err {
		case service.ErrDuplicatedID:
			ctx.Status(http.StatusBadRequest)
		default:
			ctx.Status(http.StatusInternalServerError)
		}
		return
	}
	ctx.Status(http.StatusOK)
}

func (c *Controller) postDropoff(ctx *gin.Context) {
	if ctx.Request.Method != "POST" {
		ctx.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}
	if ctx.ContentType() != "application/x-www-form-urlencoded" {
		ctx.AbortWithStatus(http.StatusUnsupportedMediaType)
		return
	}

	var dropoff struct {
		Id uint `form:"ID" binding:"required"`
	}
	if err := ctx.Bind(&dropoff); err != nil {
		return
	}

	car, err := c.service.Dropoff(dropoff.Id)
	if err == service.ErrNotFound {
		ctx.Status(http.StatusNotFound)
		return
	}
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	if car == nil {
		ctx.Status(http.StatusNoContent)
		return
	} else {
		c.service.Reassign(car)
	}
	ctx.Status(http.StatusOK)
}

func (c *Controller) postLocate(ctx *gin.Context) {
	if ctx.Request.Method != "POST" {
		ctx.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}
	if ctx.ContentType() != "application/x-www-form-urlencoded" {
		ctx.AbortWithStatus(http.StatusUnsupportedMediaType)
		return
	}

	var locate struct {
		Id uint `form:"ID" binding:"required"`
	}
	if err := ctx.Bind(&locate); err != nil {
		return
	}

	car, err := c.service.Locate(locate.Id)
	if err == service.ErrNotFound {
		ctx.Status(http.StatusNotFound)
		return
	}
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	if car == nil {
		ctx.Status(http.StatusNoContent)
		return
	}
	ctx.JSON(http.StatusOK, car)
}
