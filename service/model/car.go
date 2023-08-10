package model

import (
	"fmt"
)

type Car struct {
	Id             uint `json:"id"`
	Seats          uint `json:"seats"`
	AvailableSeats uint `json:"availableSeats"`
}

func (c *Car) String() string {
	return fmt.Sprintf(
		"#%p -> Id: %d, Seats: %d, AvailableSeats: %d",
		c,
		c.Id,
		c.Seats,
		c.AvailableSeats,
	)
}
