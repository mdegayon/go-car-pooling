package model

type Journey struct {
	Id         uint `json:"id"`
	Passengers uint `json:"passengers"`
	AssignedTo *Car `json:"assignedTo"`
}
