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
)

var eventNames = [...]string{
	"Claim",
	"DisClaim",
	"Hunter",
	"Dropp",
	"9.30 pm",
	"LessThen20",
	"Hours48",
}


func (e Event) Name() string {
	return eventNames[e];
}


func (e Event) val() Event {
	return e;
}
