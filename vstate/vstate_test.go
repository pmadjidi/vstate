package vstate

import (
	"fmt"
	"testing"
)

func Test001(t *testing.T) {
	var s State
	s.init()
	s.Print()
	if (s.Get() != Init) {
		t.Errorf("State init failed, got: %d, want: %d.", s.Get(), Init)
	}
}

func Test002(t *testing.T) {
	var s State
	s.Set(Ready, Admins)
	s.Print()
	if (s.Get() != Ready) {
		t.Errorf("State init failed, got: %d, want: %d.", s.Get(), Ready)
	}
}

func Test003(t *testing.T) {
	var s State
	_, err := s.Set(Ready, 14)
	s.Print()
	if (err != nil) {
		fmt.Print(err)
	} else {
		t.Errorf("Expecting error UnAuthorized")
	}
}

func Test004(t *testing.T) {
	var s State
	s.Set(Ready, Admins)
	s.Next(Claim, EndUser)
	s.Print()
	if (s.Get() != Riding) {
		t.Errorf("Claim failed, got: %d, want: %d.", s.Get(), Riding)
	}
}

func Test005(t *testing.T) {
	var s State
	s.Set(Riding, Admins)
	_, err := s.Next(Claim, EndUser)
	s.Print()
	fmt.Print(err)
	if (err == nil) {
		t.Errorf("Claim failed, got: %d, want: %d.", s.Get(), Riding)
	}
}

func Test006(t *testing.T) {
	var s State
	s.Set(Battery_low, Admins)
	_, err := s.Next(Claim, EndUser)
	s.Print()
	fmt.Print(err)
	if (err == nil) {
		t.Errorf("Claim failed, got: %d, want: %d.", s.Get(), Battery_low)
	}
}

func Test007(t *testing.T) {
	var s State
	s.Set(Battery_low, Admins)
	_, err := s.Next(Claim, Hunters)
	s.Print()
	if (err != nil) {
		t.Errorf("Claim failed, got: %d, want: %d.", s.Get(), Riding)
	}
}


func Test008(t *testing.T) {
	var s State
	s.Set(Service_mode, Admins)
	_, err := s.Next(Claim, Hunters)
	s.Print()
	fmt.Print(err)
	if (err == nil) {
		t.Errorf("Claim failed, got: %d, want: %d.", s.Get(), Service_mode)
	}
}

func Test009(t *testing.T) {
	var s State
	s.Set(Terminated, Admins)
	_, err := s.Next(Claim, EndUser)
	s.Print()
	fmt.Print(err)
	if (err == nil) {
		t.Errorf("Claim failed, got: %d, want: %d.", s.Get(), Terminated)
	}
}

func Test010(t *testing.T) {
	var s State
	s.Set(Bounty, Admins)
	_, err := s.Next(Claim, EndUser)
	s.Print()
	fmt.Print(err)
	if (err == nil) {
		t.Errorf("Claim failed, got: %d, want: %d.", s.Get(), Bounty)
	}
}

func Test011(t *testing.T) {
	var s State
	s.Set(Bounty, Admins)
	_, err := s.Next(Hunter, EndUser)
	s.Print()
	fmt.Print(err)
	if (err == nil) {
		t.Errorf("Claim failed, got: %d, want: %d.", s.Get(), Bounty)
	}
}

func Test012(t *testing.T) {
	var s State
	s.Set(Bounty, Admins)
	_, err := s.Next(Hunter, Hunters)
	s.Print()
	if (err != nil) {
		t.Errorf("Claim failed, got: %d, want: %d.", s.Get(), Collected)
	}
}



func Test013(t *testing.T) {
	var s State
	s.Set(Collected, Admins)
	_, err := s.Next(Hunter, Hunters)
	s.Print()
	if (err != nil) {
		t.Errorf("Claim failed, got: %d, want: %d.", s.Get(), Dropped)
	}
}


func Test014(t *testing.T) {
	var s State
	s.Set(Dropped, Admins)
	_, err := s.Next(Hunter, Hunters)
	s.Print()
	if (err != nil) {
		t.Errorf("Claim failed, got: %d, want: %d.", s.Get(), Ready)
	}
}



func Test015(t *testing.T) {
	var s State
	s.Set(Unknown, Admins)
	for ns := Claim; ns <= Hours48; ns++ {
		_, err := s.Next(ns, EndUser,Nothing)
		fmt.Print(err)
		if (err == nil) {
			fmt.Print(ns.Name())
			t.Errorf("Claim failed, got: %d, want: %d.", s.Get(), Unknown)
			break
		}
	}
}

func Test016(t *testing.T) {
	var s State
	s.Set(Terminated, Admins)
	for ns := Claim; ns <= Hours48; ns++ {
		_, err := s.Next(ns, EndUser,Nothing)
		fmt.Print(err)
		if (err == nil) {
			fmt.Print(ns.Name())
			t.Errorf("Claim failed, got: %d, want: %d.", s.Get(), Terminated)
			break
		}
	}
}

func Test017(t *testing.T) {
	var s State
	s.Set(Service_mode, Admins)
	for ns := Claim; ns <= Hours48; ns++ {
		_, err := s.Next(ns, EndUser,Nothing)
		fmt.Print(err)
		if (err == nil) {
			fmt.Print(ns.Name())
			t.Errorf("Claim failed, got: %d, want: %d.", s.Get(), Service_mode)
			break
		}
	}
}

func Test18(t *testing.T) {
	var s State
	s.Set(Battery_low, System)
	fmt.Printf("*******************")
	for ns := Claim; ns <=  Hours48; ns++ {
		_, err := s.Next(ns, EndUser,Nothing)
		if (err == nil ) {
			fmt.Printf("**Battery low and %s  %s\n",ns.String(),"Battery_low")
		}
		fmt.Printf("%s",err)
	}
}

func Test019(t *testing.T) {
	var s State
	s.Set(Bounty, Admins)
	_, err := s.Next(Hunter, Hunters)
	s.Print()
	if (err != nil) {
		t.Errorf("Claim failed, got: %d, want: %d.", s.Get(), Bounty)
	}
}

func Test020(t *testing.T) {
	var s State
	s.Set(Collected, Admins)
	_, err := s.Next(Hunter, Hunters)
	s.Print()
	if (err != nil) {
		t.Errorf("Claim failed, got: %d, want: %d.", s.Get(), Bounty)
	}
}


func Test021(t *testing.T) {
	var s State
	s.Set(Dropped, Admins)
	_, err := s.Next(Hunter, Hunters)
	s.Print()
	if (err != nil) {
		t.Errorf("Claim failed, got: %d, want: %d.", s.Get(), Bounty)
	}
}
