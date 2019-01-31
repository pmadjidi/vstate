package vstate

import (
	"fmt"
	assert2 "gotest.tools/assert"
	"testing"
)

func Test001(t *testing.T) {
	s := Ready
	s.Print()
	if (s != Ready) {
		t.Errorf("State Ready failed, got: %d, want: %d.", s, Ready)
	}
	assert2.Equal(t,s,Ready)
}

func Test002(t *testing.T) {
	s := Ready
	ns, err := s.Next(SetState, Admins, Battery_low)
	ns.Print()
	if (err != nil) {
		fmt.Print(err)
		t.Errorf("SetState failed, got: error, want: %d.", Battery_low)
	}
	assert2.Equal(t,ns,Battery_low)
}

func Test003(t *testing.T) {
	s := Ready
	ns, err := s.Next(SetState, NoAuth, Battery_low)
	ns.Print()
	if (err == nil) {
		t.Errorf("Setstate failed, got: %d, want: %d.", ns, Nothing)
	}
	assert2.Equal(t,ns,Nothing)
}

func Test004(t *testing.T) {
	s := Ready
	ns, err := s.Next(Claim, EndUser)
	s.Print()
	if (err != nil) {
		t.Errorf("Claim failed, got: %d, want: %d.", ns, Riding)
	}
	s.Print()
	assert2.Equal(t,ns,Riding)
}

func Test005(t *testing.T) {
	s := Riding
	ns, err := s.Next(Claim, EndUser)
	ns.Print()
	if (err == nil) {
		t.Errorf("Claim failed, got: %d, want: %d.", ns, Nothing)
	}
	assert2.Equal(t,ns,Nothing)
}

func Test006(t *testing.T) {
	s := Battery_low
	ns, err := s.Next(Claim, EndUser)
	ns.Print()
	if (err == nil) {
		t.Errorf("Claim failed, got: %d, want: %d.", ns, Nothing)
	}
	assert2.Equal(t,ns,Nothing)
}

func Test007(t *testing.T) {
	s := Battery_low
	ns, err := s.Next(Claim, Hunters)
	ns.Print()
	if (err != nil) {
		t.Errorf("Claim failed, got: %d, want: %d.", ns, Bounty)
	}
	assert2.Equal(t,ns,Bounty)
}

func Test008(t *testing.T) {
	s := Service_mode
	ns, err := s.Next(Claim, Hunters)
	ns.Print()
	if (err == nil) {
		t.Errorf("Claim failed, got: %d, want: %d.", ns, Nothing)
	}
	assert2.Equal(t,ns,Nothing)
}

func Test009(t *testing.T) {
	s := Terminated
	ns, err := s.Next(Claim, EndUser)
	ns.Print()
	if (err == nil) {
		t.Errorf("Claim failed, got: %d, want: %d.", ns, Nothing)
	}
	assert2.Equal(t,ns,Nothing)
}

func Test010(t *testing.T) {
	s := Bounty
	ns, err := s.Next(Claim, EndUser)
	ns.Print()
	if (err == nil) {
		t.Errorf("Claim failed, got: %d, want: %d.", ns, Nothing)
	}
	assert2.Equal(t,ns,Nothing)
}

func Test011(t *testing.T) {
	s := Bounty
	ns, err := s.Next(Hunter, EndUser)
	ns.Print()
	if (err == nil) {
		t.Errorf("Claim failed, got: %d, want: %d.", ns, Nothing)
	}
	assert2.Equal(t,ns,Nothing)
}

func Test012(t *testing.T) {
	s := Bounty
	ns, err := s.Next(Hunter, Hunters)
	ns.Print()
	if (err != nil) {
		t.Errorf("Claim failed, got: %d, want: %d.", ns, Collected)
	}
	assert2.Equal(t,ns,Collected)
}

func Test013(t *testing.T) {
	s := Collected
	ns, err := s.Next(Hunter, Hunters)
	ns.Print()
	if (err != nil) {
		t.Errorf("Claim failed, got: %d, want: %d.", ns, Dropped)
	}
	assert2.Equal(t,ns,Dropped)
}

func Test014(t *testing.T) {
	s := Dropped
	ns, err := s.Next(Hunter, Hunters)
	ns.Print()
	if (err != nil) {
		t.Errorf("Claim failed, got: %s, want: %s.", ns.String(), Ready.String())
	}
	assert2.Equal(t,ns,Ready)
}

func Test015(t *testing.T) {
	s := Unknown
	for ev := Claim; ev <= Hours48; ev++ {
		ns, err := s.Next(ev, EndUser, Nothing)
		if (err == nil) {
			fmt.Print(ev.Name())
			t.Errorf("%s, got: %s, want: %s.", ev.String(),ns.String(), Nothing.String())
			break
		}
		assert2.Equal(t,ns,Nothing)
	}
}


func Test016(t *testing.T) {
	s := Terminated
	for ev := Claim; ev <= Hours48; ev++ {
		ns, err := s.Next(ev, EndUser, Nothing)
		if (err == nil) {
			fmt.Print(ev.Name())
			t.Errorf("%s, got: %s, want: %s.",ev.String(), ns.String(), Nothing.String())
			break
		}
		assert2.Equal(t,ns,Nothing)
	}
}

func Test017(t *testing.T) {
	s := Service_mode
	for ev := Claim; ev <= Hours48; ev++ {
		ns, err := s.Next(ev, EndUser, Nothing)
		if (err == nil) {
			fmt.Print(ev.Name())
			t.Errorf("%s, got: %s, want: %s.", ev.String(),ns.String(), Nothing.String())
			break
		}
		assert2.Equal(t,ns,Nothing)
	}
}


func Test018(t *testing.T) {
	s := Battery_low
	ns, err := s.Next(LessThen20, System)
	ns.Print()
	if (err != nil) {
		t.Errorf("LessThen20, got: %s, want: %s.", ns.String(), Bounty.String())
	}
	assert2.Equal(t,ns,Bounty)
}



func Test019(t *testing.T) {
	s := Ready
	ns, err := s.Next(LessThen20, System)
	ns.Print()
	if (err == nil) {
		t.Errorf("LessThen20, got: %s, want: %s.", ns.String(), Nothing.String())
	}
	assert2.Equal(t,ns,Nothing)
}


func Test020(t *testing.T) {
	s := Riding
	ns, err := s.Next(LessThen20, System)
	ns.Print()
	if (err != nil) {
		t.Errorf("LessThen20, got: %s, want: %s.", ns.String(), Battery_low.String())
	}
	assert2.Equal(t,ns,Battery_low)
}

