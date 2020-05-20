// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package fsm

import "errors"

type State = interface{}
type Event = interface{}
type Action func(interface{}) (State, error)

type Listener interface {
	StateChanged(from, to State)

	StateEntered(state State)

	StateExited(state State)

	EventNotAccepted(event Event)

	TransitionStarted(action Action)

	TransitionEnded(action Action)

	FSMStarted(fsm FSM)

	FSMStopped(fsm FSM)

	FSMError(fsm FSM, err error)
}

var NoTransitionError = errors.New("No Transition ")

type FSM interface {
	Start() error

	Close() error

	AddListener(listener Listener)

	Initial(state State)

	Current() *State

	AddState(state State, event Event, action Action) error

	SendEvent(event Event, param interface{}) error
}

func IsNoTransition(err error) bool {
	return err == NoTransitionError
}
