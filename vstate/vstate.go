package vstate


import (
	"errors"
	"fmt"
)



func ValidState(s string) (State,bool) {
	for i := Init; i <= Nothing; i++  {
		fmt.Printf("%s,%s\n",i.String(),s)
		if i.String() == s  {
			return i,true
		}
	}
	return  -1,false
}


func (s *State) init() State {
	*s = Init;
	return *s
}

func (s State) Print() {
	fmt.Print("State is " + s.String() + "\n")
}


func (s State) Get() State {
	return s
}




func (s *State) Set(newState State,role URoles) (State,error) {
	var result State
	var err error = nil
	if newState < Init || newState > Nothing {
		return Nothing, errors.New("setState: Invalid State as input: " + newState.String())
	}
	switch (role) {
	case Admins,System:
		result,*s = newState,newState
	default:
		result = Nothing
		errorMsg := "UnAuthorized\n"
		err = errors.New(errorMsg)
	}
	return result,err
}



func (s *State) Next(newEvent Event,role URoles, state ... State) (State, error) {
	var result State
	var err error = nil
	if (newEvent < Claim || newEvent > Hours48) {
		return Nothing,errors.New("Next: Invalid Event as input: " + newEvent.Name())
	}
	switch(newEvent) {
	case Claim:
		result, err = s.claim(role)
	case DisClaim:
		result,err = s.disClaim()
	case Hunter:
		result,err = s.hunter(role)
	case Dropp:
		result,err = s.dropp(role)
	case NineThirtyPm:
		result,err = s.nineThirtyPm(role)
	case LessThen20:
		result,err = s.lessThen20()
	case Hours48:
		result,err = s.hours48()
	case SetState:
		result,err = s.Set(state[0],role)
	default:
		result = Nothing
		err = errors.New("Next: Unkown Event")
	}
	if (err == nil) {
		*s = result
		return *s,nil
	}
	return result,err
}

func (s *State) claim(role URoles) (State,error){
	var result State
	var err error = nil
	switch (*s) {
	case Ready:
		result,*s = Riding,Riding
	case Battery_low:
		if ( role == Hunters) {
			result,*s = Riding,Riding
		} else {
			result = Nothing
			errorMsg := "Claim: Only Hunter can claim vehicle in state: " + s.String() + "\n"
			err = errors.New(errorMsg)
		}
	case Init:
		if ( role == Admins) {
			result,*s = Riding,Riding
		} else {
			result = Nothing
			errorMsg := "Claim: Only Admins can claim vehicle in state: " + s.String() + "\n"
			err = errors.New(errorMsg)
		}
	default:
		result = Nothing
		errorMsg := "Claim: Can not claim vehicle in state: " + s.String() + "\n"
		err = errors.New(errorMsg)
	}
	return result,err
}


func (s *State) disClaim() (State,error){
	var result State
	var err error = nil
	switch (*s) {
	case Riding:
		result,*s = Ready,Ready
	case Battery_low:
		result,*s = Bounty,Bounty
	default:
		result = Nothing
		errorMsg := "Claim: Can not disclame vehicle in state: " + s.String() + "\n"
		err = errors.New(errorMsg)
	}
	return result,err
}



func (s *State) lessThen20() (State,error) {
	var result State
	var err error = nil
	switch (*s) {
	case Unknown,Terminated,Service_mode:
		result = Nothing
		errorMsg := "Can not change state from : " + s.String() + "\n"
		err = errors.New(errorMsg)
	default:
		result,*s = Bounty,Bounty
	}
	return result,err
}




func (s *State) hunter(role URoles)(State,error) {

	var result State
	var err error = nil
	switch (*s) {
	case Bounty:
		if (role == Hunters) {
			result, *s = Collected, Collected
		} else {
			result = Nothing
			errorMsg := "Hunter: Can not Collect vehicle in state Bounty if you are not member of Hunters\n"
			err = errors.New(errorMsg)
		}
	case Collected:
		if (role == Hunters) {
			result, *s = Dropped, Dropped
		} else {
			result = Nothing
			errorMsg := "Hunter: Can not drop vehicle in state Collected if you are not member of Hunters\n"
			err = errors.New(errorMsg)
		}
	case Dropped:
		if (role == Hunters) {
			result, *s = Ready, Ready
		} else {
			result = Nothing
			errorMsg := "Hunter: Can not return vehicle in state ready if you are not member of Hunters\n"
			err = errors.New(errorMsg)
		}
	default:
		result = Nothing
		errorMsg := "Hunter: Can not Hunt vehicle in state: " + s.String() + "\n"
		err = errors.New(errorMsg)
	}
	return result,err
}

func (s *State) hours48() (State,error){
	var result State
	var err error = nil
	switch (*s) {
	case Unknown,Terminated,Service_mode:
		result = Nothing
		errorMsg := "Can not change state from : " + s.String() + "\n"
		err = errors.New(errorMsg)
	default:
		result,*s = Unknown,Unknown
	}
	return result,err
}


func (s *State) dropp(role URoles) (State,error){
	var result State
	var err error = nil
	switch (*s) {
	case Collected:
		if (role == Hunters) {
			*s,result = Ready,Ready
		} else {
			result = Nothing
			errorMsg := "Dropp: Can can only be done by Hunters: " + s.String() + "\n"
			err = errors.New(errorMsg)
		}
	default:
		result = Nothing
		errorMsg := "Dropp: Can only be performed  on collected: " + s.String() + "\n"
		err = errors.New(errorMsg)

	}
	return result,err
}


func (s *State) nineThirtyPm(role URoles) (State,error){
	var result State
	var err error = nil
	switch (*s) {
	case Unknown,Terminated,Service_mode:
		result = Nothing
		errorMsg := "Can not change state from : " + s.String() + "\n"
		err = errors.New(errorMsg)
	default:
		result,*s = Bounty,Bounty
	}
	return result,err
}




