package model

import (
	"fmt"
)

type Car struct {
	Id             uint `json:"id"`
	MaxSeats       uint `json:"seats"`
	AvailableSeats uint `json:"availableSeats"`
}

func (c *Car) String() string {
	return fmt.Sprintf(
		"#%p -> Id: %d, MaxSeats: %d, AvailableSeats: %d",
		c,
		c.Id,
		c.MaxSeats,
		c.AvailableSeats,
	)
}
