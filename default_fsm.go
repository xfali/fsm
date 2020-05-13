// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package fsm

import (
	"fmt"
	"os"
	"sync"
)

type fsmEvent interface {
	Type() int
}

const (
	UnknownEvent = iota
	stateChangedEvent
	stateEnteredEvent
	stateExitedEvent
	eventNotAcceptedEvent
	transitionEvent
	transitionStartedEvent
	transitionEndedEvent
	fsmStartedEvent
	fsmStoppedEvent
	fsmErrorEvent

	CustomerEvent = 10000
)

type SimpleFSM struct {
	stateMap map[State]map[Event]Action
	curState interface{}
	listener Listener

	lock sync.Mutex
}

func NewSimpleFSM() *SimpleFSM {
	ret := &SimpleFSM{
		stateMap: map[interface{}]map[interface{}]Action{},
		curState: nil,
		listener: &DefaultListener{},
	}
	return ret
}

func (f *SimpleFSM) HandlerFsmEvent(e fsmEvent) {
	switch e.(type) {

	}
}

func (f *SimpleFSM) SetListener(listener Listener) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.listener = listener
}

func (f *SimpleFSM) Initial(state State) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.curState = state
	f.listener.StateEntered(f.curState)
}

func (f *SimpleFSM) Current() *State {
	f.lock.Lock()
	defer f.lock.Unlock()

	return &f.curState
}

func (f *SimpleFSM) Start() error {
	f.listener.FSMStarted(f)
	return nil
}

func (f *SimpleFSM) Close() error {
	f.listener.FSMStopped(f)
	return nil
}

func (f *SimpleFSM) AddState(state State, event Event, action Action) error {
	f.lock.Lock()
	defer f.lock.Unlock()

	if evs, ok := f.stateMap[state]; ok {
		evs[event] = action
	} else {
		f.stateMap[state] = map[interface{}]Action{
			event: action,
		}
	}
	return nil
}

func (f *SimpleFSM) Execute(event Event, param interface{}) error {
	f.lock.Lock()
	defer f.lock.Unlock()

	if evs, ok := f.stateMap[f.curState]; ok {
		if action, ok := evs[event]; ok {
			f.lock.Unlock()
			f.listener.TransitionStarted(action)
			nextState, err := action(param)
			f.listener.TransitionEnded(action)
			if err != nil {
				f.listener.FSMError(f, err)
			}
			f.listener.StateEntered(nextState)
			f.lock.Lock()
			origin := f.curState
			f.curState = nextState
			f.listener.StateExited(origin)
			f.listener.StateChanged(origin, nextState)
			return err
		} else {
			f.listener.EventNotAccepted(event)
			//log.Printf("no action found, state: %v event: %v\n", f.curState, event)
			return nil
		}
	}
	f.listener.EventNotAccepted(event)
	return nil
}

func (f *SimpleFSM) SendEvent(event Event, param interface{}) error {
	return f.Execute(event, param)
}

type eventEntity struct {
	event Event
	param interface{}
}

type DefaultFSM struct {
	SimpleFSM

	eventChanSize int
	eventChan     chan eventEntity
	stopChan      chan bool
}

type Opt func(f *DefaultFSM)

func New(opts ...Opt) *DefaultFSM {
	ret := &DefaultFSM{
		SimpleFSM: *NewSimpleFSM(),
		stopChan:  make(chan bool),
	}
	for i := range opts {
		opts[i](ret)
	}
	if ret.eventChanSize <= 0 {
		ret.eventChanSize = 1024
	}
	ret.eventChan = make(chan eventEntity, 1024)

	return ret
}

func SetEventBufferSize(size int) Opt {
	return func(f *DefaultFSM) {
		f.eventChanSize = size
	}
}

func (f *DefaultFSM) Start() error {
	go func() {
		for {
			select {
			case entity := <-f.eventChan:
				f.Execute(entity.event, entity.param)
				break
			case <-f.stopChan:
				return
			}
		}
	}()
	f.listener.FSMStarted(f)
	return nil
}

func (f *DefaultFSM) Close() error {
	close(f.stopChan)
	f.listener.FSMStopped(f)
	return nil
}

func (f *DefaultFSM) SendEvent(event Event, param interface{}) error {
	f.eventChan <- eventEntity{event: event, param: param}
	return nil
}

type DefaultListener struct{ Silent bool }

func (l *DefaultListener) StateChanged(from, to State) {
	if l.Silent {
		return
	}
	fmt.Fprintf(os.Stdout, "StateChanged: from %v to %v\n", from, to)
}

func (l *DefaultListener) StateEntered(state State) {
	if l.Silent {
		return
	}
	fmt.Fprintf(os.Stdout, "StateEntered: %v\n", state)
}

func (l *DefaultListener) StateExited(state State) {
	if l.Silent {
		return
	}
	fmt.Fprintf(os.Stdout, "StateExited: %v\n", state)
}

func (l *DefaultListener) EventNotAccepted(event Event) {
	if l.Silent {
		return
	}
	fmt.Fprintf(os.Stdout, "EventNotAccepted: %v\n", event)
}

func (l *DefaultListener) Transition(action Action) {
	//fmt.Fprintf(os.Stdout, "Transition\n")
}

func (l *DefaultListener) TransitionStarted(action Action) {
	//fmt.Fprintf(os.Stdout, "TransitionStarted\n")
}

func (l *DefaultListener) TransitionEnded(action Action) {
	//fmt.Fprintf(os.Stdout, "TransitionEnded\n")
}

func (l *DefaultListener) FSMStarted(fsm FSM) {
	if l.Silent {
		return
	}
	fmt.Fprintf(os.Stdout, "FSMStarted with state: %v\n", *fsm.Current())
}

func (l *DefaultListener) FSMStopped(fsm FSM) {
	if l.Silent {
		return
	}
	fmt.Fprintf(os.Stdout, "FSMStopped with state: %v\n", *fsm.Current())
}

func (l *DefaultListener) FSMError(fsm FSM, err error) {
	if l.Silent {
		return
	}
	fmt.Fprintf(os.Stderr, "FSMError with state: %v and error: %v\n", *fsm.Current(), err)
}
