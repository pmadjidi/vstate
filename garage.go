package main

import "sync"

type Garage struct {
	sync.RWMutex
	internal map[string]*Vehicle
}

func NewGarage() *Garage {
	return &Garage{
		internal: make(map[string]*Vehicle),
	}
}


func (g *Garage) getVehicleById(key string) (v *Vehicle, ok bool) {
	g.RLock()
	v, ok = g.internal[key]
	g.RUnlock()
	return v, ok
}

func (g *Garage) getMap() map[string]*Vehicle {
	copy := make( map[string]*Vehicle)
	g.RLock()
	for k,v := range g.internal {
		copy[k] = v
	}
	g.RUnlock()
	return copy
}

func (g *Garage) Delete(key string) {
	g.Lock()
	delete(g.internal, key)
	g.Unlock()
}

func (g *Garage) Set(key string, v *Vehicle) {
	g.Lock()
	g.internal[key] = v
	g.Unlock()
}

