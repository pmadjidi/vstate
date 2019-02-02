package vstate


import (
	"errors"
	"fmt"
)



func ValidState(s string) (State,bool) {
	for i := Ready; i <= Nothing; i++  {
		if i.String() == s  {
			return i,true
		}
	}
	return  -1,false
}



func (s State) Print() {
	fmt.Print("\nState is " + s.String() + "\n")
}




func adminSystemSet(newState State,role URoles) (State,error) {
	var result State
	var err error = nil
	if newState < Ready || newState > Nothing {
		return Nothing, errors.New("setState: Invalid State as input: " + newState.String())
	}
	switch (role) {
	case Admins,System:
		result= newState
	default:
		result = Nothing
		errorMsg := "UnAuthorized\n"
		err = errors.New(errorMsg)
	}
	return result,err
}



func (s State) Next(newEvent Event,role URoles,state ... State) (State, error) {
	var result State
	var err error = nil
	if (newEvent < Claim || newEvent > Hours48) {
		return Nothing,errors.New("Next: Invalid Event as input: " + newEvent.Name())
	}
	switch(newEvent) {
	case Claim:
		result, err = claim(s,role)
	case DisClaim:
		result,err =  disClaim(s)
	case Hunter:
		result,err =  hunter(s,role)
	case NineThirtyPm:
		result,err =  nineThirtyPm(s,role)
	case LessThen20:
		result,err =  lessThen20(s)
	case Hours48:
		result,err = hours48(s)
	case SetState:
		if len(state) == 0 {
			result = Nothing
			err = errors.New("Argument to set state is missing....")
		} else {
			result, err = adminSystemSet(state[0], role)
		}
	default:
		result = Nothing
		err = errors.New("Next: Unkown Event")
	}
	return result,err
}

func claim(oldState State,role URoles) (State,error){
	var result State
	var err error = nil
	switch (oldState) {
	case Ready:
		result = Riding
	case Battery_low:
		if ( role == Hunters) {
			result  = Bounty
		} else {
			result = Nothing
			errorMsg := "Claim: Only Hunter can claim vehicle in state: " + oldState.String() + "\n"
			err = errors.New(errorMsg)
		}
	default:
		result = Nothing
		errorMsg := "Claim: Can not claim vehicle in state: " + oldState.String() + "\n"
		err = errors.New(errorMsg)
	}
	return result,err
}


func  disClaim(oldState State) (State,error){
	var result State
	var err error = nil
	switch (oldState) {
	case Riding:
		result = Ready
	default:
		result = Nothing
		errorMsg := "Claim: Can not disclame vehicle in state: " + oldState.String() + "\n"
		err = errors.New(errorMsg)
	}
	return result,err
}



func lessThen20(oldState State) (State,error) {
	var result State
	var err error = nil
	switch (oldState) {
	case Riding:
		result  = Battery_low
	case Battery_low,Bounty:
		result = Bounty

	default:
		result = Nothing
		errorMsg := "Can not change state from : " + oldState.String() + "\n"
		err = errors.New(errorMsg)
	}
	return result,err
}




func  hunter(oldState State,role URoles)(State,error) {

	var result State
	var err error = nil
	switch (oldState) {
	case Bounty:
		if (role == Hunters) {
			result= Collected
		} else {
			result = Nothing
			errorMsg := "Hunter: Can not Collect vehicle in state Bounty if you are not member of Hunters\n"
			err = errors.New(errorMsg)
		}
	case Collected:
		if (role == Hunters) {
			result = Dropped
		} else {
			result = Nothing
			errorMsg := "Hunter: Can not drop vehicle in state Collected if you are not member of Hunters\n"
			err = errors.New(errorMsg)
		}
	case Dropped:
		if (role == Hunters) {
			result  = Ready
		} else {
			result = Nothing
			errorMsg := "Hunter: Can not return vehicle in state ready if you are not member of Hunters\n"
			err = errors.New(errorMsg)
		}
	default:
		result = Nothing
		errorMsg := "Hunter: Can not Hunt vehicle in state: " + oldState.String() + "\n"
		err = errors.New(errorMsg)
	}
	return result,err
}

func  hours48(oldState State) (State,error){
	var result State
	var err error = nil
	switch (oldState) {
	case Ready:
		result = Unknown
	default:
		result = Nothing
		errorMsg := "nineThirtyPm state change possible only from:" + Ready.String() + "\n"
		err = errors.New(errorMsg)
	}
	return result,err
}




func  nineThirtyPm(oldState State,role URoles) (State,error){
	var result State
	var err error = nil
	switch (oldState) {
	case Ready:
		result = Bounty
	default:
		result = Nothing
		errorMsg := "nineThirtyPm state change possible only from:" + Ready.String() + "\n"
		err = errors.New(errorMsg)
	}
	return result,err
}




