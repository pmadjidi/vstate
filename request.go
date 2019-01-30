package main

import "vehicles/vstate"

type Response struct {
	state vstate.State
	err error
}

type Request struct {
	event vstate.Event
	userRole vstate.URoles
	state vstate.State
	resp chan *Response
}

