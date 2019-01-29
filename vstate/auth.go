package vstate

type URoles int

const (
	NoAuth URoles = iota
	EndUser
	Hunters
	Admins
	System
)



