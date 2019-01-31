package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"vehicles/vstate"
)

type Vehicle struct {
	sync.RWMutex
	vstate.State `json:"state"`
	Id           int           `json:"-"`
	Uid          string        `json:"id"`
	Battery      int           `json:"battery"`
	Port         chan *Request `json:"-"`
	CreatedAt    int64         `json:"createdAt"`
	UpdatedAt    int64         `json:"updatedAt"`
}

func (v *Vehicle) readValues() (uid string,state vstate.State,battery int,createdAt int64,updatedAt int64){
	v.RLock()
	uid = v.Uid
	state = v.State
	battery = v.Battery
	createdAt = v.CreatedAt
	updatedAt = v.UpdatedAt
	v.RUnlock()
	return uid,state,battery,createdAt,updatedAt
}

func (v *Vehicle) serilize() ([]byte,error) {
	v.RLock()
	b, err := json.Marshal(v)
	v.RUnlock()
	return b,err
}

func (v *Vehicle) print() {
	v.RLock()
	fmt.Printf("%#v\n", v)
	v.RUnlock()
}

func (v *Vehicle) stamp() vstate.State {
	v.Lock()
	v.UpdatedAt = time.Now().UnixNano()
	v.Unlock()
	return v.State
}

func (v *Vehicle) inState() string {
	v.RLock()
	res := v.State.String()
	v.RUnlock()
	return res
}

func (v *Vehicle) getLastMod() int64 {
	v.RLock()
	res := v.UpdatedAt
	v.RUnlock()
	return res
}

func (v *Vehicle) getUid() string {
	v.RLock()
	res := v.Uid
	v.RUnlock()
	return res
}

func (v *Vehicle) getCreatedAt() int64 {
	v.RLock()
	res := v.CreatedAt
	v.RUnlock()
	return res
}


func (v *Vehicle) settLastMod(time int64) int64 {
	v.Lock()
	v.UpdatedAt = time
	v.Unlock()
	return time
}

func (v *Vehicle) setCreatedAt(time int64) int64 {
	v.Lock()
	v.CreatedAt = time
	v.Unlock()
	return time
}

func (v *Vehicle) ChargeLevel() int {
	v.RLock()
	res :=  v.Battery
	v.RUnlock()
	return res
}

func (v *Vehicle) getState() vstate.State {
	v.RLock()
	res := v.State
	v.RUnlock()
	return res
}

func (v *Vehicle) setState(state vstate.State) vstate.State {
	v.Lock()
	v.State = state
	v.Unlock()
	return state
}

func (v *Vehicle) listen() {
	go func() {
		<-app.start
		Loop:
		for {
			select {
			case r := <-v.Port:
				if (r.event == vstate.Delete) {
					res := &Response{vstate.Nothing, nil}
					r.resp <- res
					fmt.Print("Terminating event loop for Id: " + v.Uid + "in State: " + r.event.Name())
					break
				}
				fmt.Print("Id: " + v.Uid + "in State: " + r.event.Name())
				oldState := v.getState()
				nextState, err := v.Next(r.event, r.userRole, r.state)
				if (err != nil) {
					fmt.Print("Error: ", v.Id, err)
					res := &Response{v.Get(), err}
					r.resp <- res
				} else {
					res := &Response{v.Get(), nil}
					r.resp <- res
					if (oldState != nextState.Get()) {
						v.stamp()
						app.store <- v
					}

				}
			case <-time.After(10 * time.Second):
				fmt.Print("Vehicle id: " + v.Uid + " in State: " + v.String() + "  is Alive....\n")
				if v.Battery < 20 {
					_,err := v.Next(vstate.LessThen20, vstate.System,vstate.Nothing)
					if (err != nil) {
						fmt.Printf("\n###%s %s %s\n",v.Uid,v.String(),err)
					}
				}
			case <-app.quit:
				fmt.Printf("\nStopping Vehicle Id %s event loop...\n",v.Uid)
				break Loop
			}
		}
	}()
}

func NewVehicle() *Vehicle {
	t :=  time.Now().UnixNano()
	v := Vehicle{State:vstate.Init,Uid: uniqueId(), Port: make(chan *Request), Battery: 100,
	CreatedAt: t ,UpdatedAt:t}
	app.store <- &v
	v.listen()
	return &v
}
