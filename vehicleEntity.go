package main

import (
	"vehicles/vstate"
)


type VehicleEntity struct {
	vstate.State `json:"state"`
	Id           int           `json:"-"`
	Uid          string        `json:"id"`
	Battery      int           `json:"battery"`
	CreatedAt    int64         `json:"createdAt"`
	UpdatedAt    int64         `json:"updatedAt"`
}


func (v *VehicleEntity) id() int {
	return v.Id
}


func (v *VehicleEntity) uid() string {
	return v.Uid
}


func (v *VehicleEntity) charge() int {
	return v.Battery
}


func (v *VehicleEntity) createdAt() int64 {
	return v.CreatedAt
}

func (v *VehicleEntity) updatedAt() int64 {
	return v.UpdatedAt
}




