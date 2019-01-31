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

func (v *Vehicle) presist() {
	v.RLock()
	nv := *v
	v.RUnlock()
	app.store <- &nv
}

func (v *Vehicle) clone() *Vehicle {
	v.RLock()
	nv := Vehicle{sync.RWMutex{}, v.State, v.Id, v.Uid, v.Battery, make(chan *Request), v.CreatedAt, v.UpdatedAt}
	v.RUnlock()
	return &nv
}

func (v *Vehicle) serilize() ([]byte, error) {
	v.RLock()
	b, err := json.Marshal(v)
	v.RUnlock()
	return b, err
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
	v.presist()
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
	v.presist()
	return time
}

func (v *Vehicle) ChargeLevel() int {
	v.RLock()
	res := v.Battery
	v.RUnlock()
	return res
}

func (v *Vehicle) setChargeLevel(charge int) int {
	v.Lock()
	v.Battery = charge
	v.UpdatedAt = time.Now().UnixNano()
	v.Unlock()
	v.presist()
	return charge
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
	v.UpdatedAt = time.Now().UnixNano()
	v.Unlock()
	v.presist()
	return state
}

func (v *Vehicle) Event(ev vstate.Event, r vstate.URoles, s vstate.State) (vstate.State, error) {
	v.Lock()
	defer v.Unlock()
	ns, err := v.Next(ev, r, s)
	if (err == nil) {
		v.State = ns
		return ns, err
	} else {
		fmt.Printf("Event:, Next retuns Error in error %s %s %s\n",ev.String(), v.State.String(), ns.String())
		return v.State, err
	}

}

func doNothing() {}

func (v *Vehicle) listen() {
	go func() {
		<-app.start
	Loop:
		for {
			select {
			case r := <-v.Port:
				switch r.event {
				case vstate.Delete:
					res := &Response{vstate.Nothing, nil}
					r.resp <- res
					fmt.Print("Terminating event loop for Id: " + v.Uid + "in State: " + r.event.Name())
					break
				default:
					doNothing()
				}
				_, err := v.Event(r.event, r.userRole, r.state)
				if (err != nil) {
					fmt.Print("Error: ", v.Id, err)
					res := &Response{v.getState(), err}
					r.resp <- res
				} else {
					res := &Response{v.getState(), nil}
					r.resp <- res
				}

			case <-time.After(10 * time.Second):
				charge := v.ChargeLevel()
				fmt.Printf("Id: %s in State: %s battery level %d\n", v.getUid(), v.getState(), v.ChargeLevel())
				if charge < 20 {
					_, err := v.Event(vstate.LessThen20, vstate.System, vstate.Nothing)
					if (err != nil) {
						fmt.Printf("\n%s %s %s\n", v.Uid, v.String(), err)
					}
				}
				if charge >= 19  && v.getState() == vstate.Riding {
					v.setChargeLevel(charge - 15)
				}

			case <-app.quit:
				fmt.Printf("\nStopping Vehicle Id %s event loop...\n", v.Uid)
				break Loop
			}
		}
	}()
}

func NewVehicle() *Vehicle {
	t := time.Now().UnixNano()
	v := Vehicle{State: vstate.Ready, Uid: uniqueId(), Port: make(chan *Request), Battery: 100,
		CreatedAt: t, UpdatedAt: t}
	app.store <- &v
	v.listen()
	return &v
}
