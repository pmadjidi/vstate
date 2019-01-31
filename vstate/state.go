package vstate

type State int

const (
	Ready State = iota
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

