package model

type Car struct {
	ID             uint `json:"id"`
	MaxSeats       uint `json:"maxSeats"`
	AvailableSeats uint `json:"availableSeats"`
}
