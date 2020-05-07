// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package fsm

import "errors"

type State = interface{}
type Event = interface{}
type Action func(interface{}) (State, error)

var NoTransitionError = errors.New("No Transition ")

type FSM interface {
	Start() error
	Close() error

	Initial(state State)
	Current() *State
	AddState(state State, event Event, action Action) error
	SendEvent(event Event, param interface{}) error
}

func IsNoTransition(err error) bool {
	return err == NoTransitionError
}
