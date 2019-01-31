package vstate

type Event int


const (
	Claim Event = iota // always first
	DisClaim
	Hunter
	NineThirtyPm
	LessThen20
	SetState
	Delete
    Hours48 //always last
)





func (e Event) Name() string {
	return e.String()
}

func (e Event) val() Event {
	return e;
}
