package service

import (
	"errors"
	"fmt"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/service/model"
	"log"
)

const (
	JourneyMaxValidSeats uint = 6
	JourneyMinValidSeats uint = 1
	HopIn                uint = 0
	DropOff              uint = 1
)

var ErrNotFound = errors.New("not found")
var ErrDuplicatedID = errors.New("duplicated Id")
var ErrInvalidPassengersNumber = errors.New("wrong number of passengers")
var ErrInvalidSeatsNumber = errors.New(fmt.Sprintf(
	"invalid number of seats. Expected: (%d - %d))",
	JourneyMinValidSeats, JourneyMaxValidSeats,
))

type CarPool struct {
	cars                 map[uint]*model.Car
	Journeys             map[uint]*model.Journey
	pending              []*model.Journey
	carsByAvailableSeats [JourneyMaxValidSeats + 1][]*model.Car
}

func NewCarpool() *CarPool {
	return &CarPool{
		cars:     make(map[uint]*model.Car),
		Journeys: make(map[uint]*model.Journey),
		pending:  make([]*model.Journey, 0),
	}
}

func (cp *CarPool) ResetCars(cars []*model.Car) error {

	cp.initCarpool()

	carMap := make(map[uint]*model.Car, len(cars))
	carsBySeats := new([JourneyMaxValidSeats + 1][]*model.Car)

	for _, currentCar := range cars {

		if currentCar.Seats <= 0 {
			return ErrInvalidSeatsNumber
		}

		//Check whether currentCar's id was already used
		if _, duplicatedCar := carMap[currentCar.Id]; duplicatedCar {

			return ErrDuplicatedID
		}

		carMap[currentCar.Id] = currentCar
		carsBySeats[currentCar.AvailableSeats] = append(carsBySeats[currentCar.AvailableSeats], currentCar)
	}

	cp.cars = carMap
	cp.carsByAvailableSeats = *carsBySeats

	return nil
}

func (cp *CarPool) NewJourney(journey *model.Journey) error {

	if _, duplicatedJourney := cp.Journeys[journey.Id]; duplicatedJourney {
		return ErrDuplicatedID
	}

	if journey.People < JourneyMinValidSeats || journey.People > JourneyMaxValidSeats {
		return ErrInvalidPassengersNumber
	}

	var selectedCar *model.Car
	for _, availableSeats := range cp.getWorstFitForSeats(journey.People) {

		if len(cp.carsByAvailableSeats[availableSeats]) > 0 {
			selectedCar = cp.popCarFromSeatsArray(availableSeats)
			break
		}

	}

	if selectedCar != nil {
		journey.AssignedTo = selectedCar
		err := cp.updateCarAvailableSeats(selectedCar, journey.People, HopIn)
		if err != nil {
			return err
		}
	} else {
		cp.pending = append(cp.pending, journey)
	}

	cp.Journeys[journey.Id] = journey

	return nil
}

func (cp *CarPool) Dropoff(journeyId uint) (car *model.Car, err error) {

	journey := cp.Journeys[journeyId]
	if journey == nil {
		return nil, ErrNotFound
	}

	car = journey.AssignedTo
	delete(cp.Journeys, journeyId)

	//Journey assigned to a car
	if car != nil {

		cp.removeCarFromSeatsArray(car.Id, car.AvailableSeats)
		err := cp.updateCarAvailableSeats(car, journey.People, DropOff)
		if err != nil {
			return nil, err
		}

		//Pending journey
	} else {
		cp.removePendingJourney(journeyId)
	}
	return car, nil
}

func (cp *CarPool) Reassign(car *model.Car) {

	startingSeats := car.AvailableSeats

	i := 0
	for car.AvailableSeats > 0 && i < len(cp.pending) {

		pendingJourney := cp.pending[i]

		if pendingJourney.People <= car.AvailableSeats {

			pendingJourney.AssignedTo = car
			car.AvailableSeats -= pendingJourney.People

			cp.pending = append(
				cp.pending[:i],
				cp.pending[i+1:]...,
			)

		} else {
			i++
		}
	}

	if startingSeats == car.AvailableSeats {
		return
	}

	cp.removeCarFromSeatsArray(car.Id, startingSeats)

	cp.carsByAvailableSeats[car.AvailableSeats] = append(cp.carsByAvailableSeats[car.AvailableSeats], car)
}

func (cp *CarPool) Locate(journeyId uint) (car *model.Car, err error) {

	journey := cp.Journeys[journeyId]
	if journey == nil {
		return nil, ErrNotFound
	}

	car = nil
	if journey.AssignedTo != nil {
		car = journey.AssignedTo
	}

	return car, nil
}

func (cp *CarPool) initCarpool() {

	cp.cars = make(map[uint]*model.Car)
	cp.Journeys = make(map[uint]*model.Journey)
	cp.pending = make([]*model.Journey, 0)
	for i := 0; i < len(cp.carsByAvailableSeats); i++ {
		cp.carsByAvailableSeats[i] = make([]*model.Car, 0)
	}
}

func (cp *CarPool) getWorstFitForSeats(perfectFitSeats uint) []uint {

	worstFits := make([]uint, JourneyMaxValidSeats-perfectFitSeats+1)

	worstFits[0] = perfectFitSeats

	for seats, i := JourneyMaxValidSeats, 1; seats > perfectFitSeats; seats, i = seats-1, i+1 {
		worstFits[i] = seats
	}

	return worstFits
}

func (cp *CarPool) getBestFitSeats(perfectFitSeats uint) []uint {

	bestFits := make([]uint, JourneyMaxValidSeats-perfectFitSeats+1)

	for seats, i := perfectFitSeats, 0; seats <= JourneyMaxValidSeats; seats, i = seats+1, i+1 {
		bestFits[i] = seats
	}

	return bestFits
}

func (cp *CarPool) updateCarAvailableSeats(car *model.Car, passengers uint, op uint) error {

	var updatedSeats uint
	if op == DropOff {
		updatedSeats = car.AvailableSeats + passengers
	} else if op == HopIn {
		updatedSeats = car.AvailableSeats - passengers
	} else {
		return errors.New("Unknown operation")
	}

	car.AvailableSeats = updatedSeats

	cp.carsByAvailableSeats[updatedSeats] = append(cp.carsByAvailableSeats[updatedSeats], car)

	return nil
}

func (cp *CarPool) removePendingJourney(journeyId uint) {

	for i, p := range cp.pending {
		if p.Id == journeyId {
			cp.pending = append(
				cp.pending[:i],
				cp.pending[i+1:]...,
			)
			return
		}
	}
}

func (cp *CarPool) removeCarFromSeatsArray(carId uint, seats uint) {

	for i, currentCar := range cp.carsByAvailableSeats[seats] {
		if currentCar.Id == carId {

			carsLen := len(cp.carsByAvailableSeats[seats])
			cp.carsByAvailableSeats[seats][i] = cp.carsByAvailableSeats[seats][carsLen-1]
			cp.carsByAvailableSeats[seats] = cp.carsByAvailableSeats[seats][:carsLen-1]
			return
		}
	}

}

func (cp *CarPool) popCarFromSeatsArray(seats uint) *model.Car {

	carsCount := len(cp.carsByAvailableSeats[seats])

	selectedCar := cp.carsByAvailableSeats[seats][carsCount-1]

	cp.carsByAvailableSeats[seats] = cp.carsByAvailableSeats[seats][:carsCount-1]

	return selectedCar
}

func (cp *CarPool) LogService() {

	log.Println("\n\n~~~~ CarPool ~~~~ ")

	log.Println(" -> Journeys: ")
	if len(cp.Journeys) == 0 {
		log.Print("      -")
	} else {
		for _, journey := range cp.Journeys {
			log.Println("      ", journey)
		}
	}

	log.Println(" -> Cars: ")
	if len(cp.cars) == 0 {
		log.Println("      -")
	} else {
		for _, car := range cp.cars {
			log.Println("      ", car)
		}
	}

	log.Println(" -> Pending: ")
	if len(cp.pending) == 0 {
		log.Println("      -")
	} else {
		for _, p := range cp.pending {
			log.Println("      ", p)
		}
	}

	log.Println(" -> Cars by Available Seats:")
	for i, cars := range cp.carsByAvailableSeats {
		log.Printf("       '%d' Seats:\n", i)
		for _, car := range cars {
			log.Println("         ", car)
		}
	}

	log.Println("~~~~ /CarPool ~~~~ ")
}
