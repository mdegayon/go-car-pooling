package service

import (
	"errors"
	"fmt"

	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/service/model"
)

const MaxSeats = 6

var ErrNotFound = errors.New("not found")
var ErrDuplicatedID = errors.New("duplicated ID")

type CarPool struct {
	cars     []*model.Car
	journeys []*model.Journey
	pending  []*model.Journey
}

func New_CarPool() *CarPool {
	return &CarPool{
		cars:     make([]*model.Car, 0),
		journeys: make([]*model.Journey, 0),
		pending:  make([]*model.Journey, 0),
	}
}

func (cp *CarPool) ResetCars(cars []*model.Car) error {
	cp.cars = make([]*model.Car, 0)
	cp.journeys = make([]*model.Journey, 0)
	cp.pending = make([]*model.Journey, 0)
	for _, car := range cars {
		for _, car2 := range cp.cars {
			if car.ID == car2.ID {
				cp.cars = make([]*model.Car, 0)
				cp.journeys = make([]*model.Journey, 0)
				cp.pending = make([]*model.Journey, 0)
				return ErrDuplicatedID
			}
			if car.MaxSeats < 4 || car.MaxSeats > 6 {
				cp.cars = make([]*model.Car, 0)
				cp.journeys = make([]*model.Journey, 0)
				cp.pending = make([]*model.Journey, 0)
				return errors.New("invalid seats")
			}
		}
		cp.cars = append(cp.cars, car)
	}
	return nil
}

func (cp *CarPool) NewJourney(journey *model.Journey) error {
	for _, j := range cp.journeys {
		if journey.Id == j.Id {
			return ErrDuplicatedID
		}
	}

	var selCar *model.Car
	for _, car := range cp.cars {
		if journey.Passengers <= car.AvailableSeats {
			selCar = car
			break
		}
	}

	if selCar != nil {
		journey.AssignedTo = selCar
		selCar.AvailableSeats -= journey.Passengers
	} else {
		cp.pending = append(cp.pending, journey)
	}
	cp.journeys = append(cp.journeys, journey)

	return nil
}

func (cp *CarPool) Dropoff(journey_id uint) (car *model.Car, err error) {
	var journey *model.Journey
	for _, j := range cp.journeys {
		if j.Id == journey_id {
			journey = j
		}
	}
	if journey == nil {
		return nil, ErrNotFound
	}

	car = journey.AssignedTo
	for i, j := range cp.journeys {
		if j == journey {
			cp.journeys = append(
				cp.journeys[:i],
				cp.journeys[i+1:]...,
			)
		}
	}

	if car != nil {
		car.AvailableSeats += journey.Passengers
	} else {
		for i, p := range cp.pending {
			if p.Id == journey_id {
				cp.pending = append(
					cp.pending[:i],
					cp.pending[i+1:]...,
				)
				return
			}
		}
	}
	return car, nil
}

func (cp *CarPool) Reassign(car *model.Car) {
	for i, p := range cp.pending {
		if p.Passengers <= car.AvailableSeats {
			fmt.Printf(">> Car %d reassigned to journey %d\n", car.ID, p.Id)
			p.AssignedTo = car
			car.AvailableSeats -= p.Passengers
			cp.pending = append(
				cp.pending[:i],
				cp.pending[i+1:]...,
			)
			break
		}
	}
}

func (cp *CarPool) Locate(journeyID uint) (car *model.Car, err error) {
	var journey *model.Journey
	for _, j := range cp.journeys {
		if j.Id == journeyID {
			journey = j
		}
	}
	if journey == nil {
		return nil, ErrNotFound
	}
	if journey.AssignedTo != nil {
		car = journey.AssignedTo
	}
	return car, err
}

func (cp *CarPool) findCar(seats uint) (car *model.Car) {
	for _, c := range cp.cars {
		if c.AvailableSeats >= seats {
			car = c
		}
	}
	return car
}
