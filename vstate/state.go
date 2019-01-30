package vstate

type State int

const (
	Init State = iota
	Ready
	Battery_low
	Bounty
	Riding
	Collected
	Dropped
	Service_mode
	Terminated
	Unknown
	Nothing
);

