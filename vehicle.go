package main

import (
	"fmt"
	"time"
	"vehicles/vstate"
)

type Vehicle struct {
	vstate.State  `json:"state"`
	Id         int  `json:"-"`
	Uid        string `json:"id"`
	Battery   int `json:"battery"`
	Port      chan Request  `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (v *Vehicle) print() {
	fmt.Printf("%#v\n", v)
}

func (v *Vehicle) stamp() vstate.State {
	v.UpdatedAt = time.Now()
	return v.State
}

func (v *Vehicle) inState() string {
	return v.State.Name()
}

func (v *Vehicle) lastMod() time.Time {
	return v.UpdatedAt
}

func (v *Vehicle) createdAt() time.Time {
	return v.CreatedAt
}

func (v *Vehicle) ChargeLevel() int {
	return v.Battery
}

func (v *Vehicle) listen() {
	app.wg.Add(1)
	go func() {
		<- app.start
		for {
			select {
			case r := <-v.Port:
				fmt.Print( "Id: " + v.Uid + "in State: " + r.event.Name())
				oldState := v.Get()
				res, err := v.Next(r.event, r.userRole)
				if (err != nil) {
					fmt.Print("Error: ", v.Id, err)
				} else {
					if (oldState != res.Get()) {
						v.stamp()
						app.store <- v
					}
				}
			case <-time.After(10 * time.Second):
				fmt.Print( "Vehicle id: " + v.Uid + " in State: " + v.Name() + "  is Alive....\n")
			case <- app.quit:
				break
			}
		}
	}()
}

func NewVehicle() *Vehicle {
	v := Vehicle{Uid: uniqueId(), Port: make(chan Request), Battery: 100, CreatedAt: time.Now()}
	v.Set(vstate.Init, vstate.System)
	v.stamp()
	app.store <- &v
	v.listen()
	return &v
}
