package main

import "vehicles/vstate"

type Request struct {
	event vstate.Event
	userRole vstate.URoles
}
