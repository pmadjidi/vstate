package vstate

type Event int


const (
	Claim Event = iota
	DisClaim
	Hunter
	Dropp
	NineThirtyPm
	LessThen20
    Hours48
	SetState
)





func (e Event) Name() string {
	return e.String()
}

func (e Event) val() Event {
	return e;
}
