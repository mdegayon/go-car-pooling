package model

import "fmt"

type Journey struct {
	Id         uint `json:"id"`
	People     uint `json:"people"`
	AssignedTo *Car `json:"assignedTo"`
}

func (j *Journey) String() string {
	return fmt.Sprintf("[#%p -> Id: %d, People: %d, AssignedTo: %v]", j, j.Id, j.People, j.AssignedTo)
}
