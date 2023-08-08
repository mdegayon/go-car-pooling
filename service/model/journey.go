package model

import "fmt"

type Journey struct {
	Id         uint `json:"id"`
	Passengers uint `json:"people"`
	AssignedTo *Car `json:"assignedTo"`
}

func (j *Journey) String() string {
	return fmt.Sprintf("[#%p -> Id: %d, Passengers: %d, AssignedTo: %v]", j, j.Id, j.Passengers, j.AssignedTo)
}
