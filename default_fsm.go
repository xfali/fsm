// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package fsm

import "log"

type SimpleFSM struct {
	stateMap map[State]map[Event]Action

	curState interface{}
}

func NewSimpleFSM() *SimpleFSM {
	ret := &SimpleFSM{
		stateMap: map[interface{}]map[interface{}]Action{},
		curState: nil,
	}
	return ret
}

func (f *SimpleFSM) Initial(state State) {
	f.curState = state
}

func (f *SimpleFSM) Current() *State {
	return &f.curState
}

func (f *SimpleFSM) Start() error {
	return nil
}

func (f *SimpleFSM) Close() error {
	return nil
}

func (f *SimpleFSM) AddState(state State, event Event, action Action) error {
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
	if evs, ok := f.stateMap[f.curState]; ok {
		if action, ok := evs[event]; ok {
			nextState, err := action(param)
			f.curState = nextState
			return err
		} else {
			log.Printf("no action found, state: %v event: %v\n", f.curState, event)
			return NoTransitionError
		}
	}
	return NoTransitionError
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
	return nil
}

func (f *DefaultFSM) Close() error {
	close(f.stopChan)
	return nil
}

func (f *DefaultFSM) SendEvent(event Event, param interface{}) error {
	f.eventChan <- eventEntity{event: event, param: param}
	return nil
}
