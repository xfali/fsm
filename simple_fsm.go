// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package fsm

import (
	"sync"
)

type SimpleFSM struct {
	stateMap  map[State]map[Event]Action
	curState  interface{}
	listeners []Listener
	sender    MessageSender

	lock sync.Mutex
}

func NewSimpleFSM() *SimpleFSM {
	ret := &SimpleFSM{
		stateMap: map[interface{}]map[interface{}]Action{},
		curState: nil,
	}
	ret.sender = ret
	//ret.AddListener(&DefaultListener{})
	return ret
}

func (f *SimpleFSM) handleMsg(m Message) {
	m.Proc()
}

func (f *SimpleFSM) SendMessage(m Message) {
	f.handleMsg(m)
}

func (f *SimpleFSM) AddListener(listener Listener) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.listeners = append(f.listeners, listener)
}

func (f *SimpleFSM) Initial(state State) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.curState = state
	f.StateEntered(f.curState)
}

func (f *SimpleFSM) Current() *State {
	f.lock.Lock()
	defer f.lock.Unlock()

	return &f.curState
}

func (f *SimpleFSM) Start() error {
	f.sender.SendMessage(&fsmStartedMsg{
		l:   f,
		fsm: f,
	})
	return nil
}

func (f *SimpleFSM) Close() error {
	f.sender.SendMessage(&fsmStoppedMsg{
		l:   f,
		fsm: f,
	})
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
			//为了减小锁粒度，在执行action时解锁，由于DefaultFSM用chan处理event保证了互斥
			nextState, err := func() (State, error) {
				f.lock.Unlock()
				defer f.lock.Lock()

				f.sender.SendMessage(&transitionStartedMsg{
					l:      f,
					action: action,
				})
				next, err := action(param)
				f.sender.SendMessage(&transitionEndedMsg{
					l:      f,
					action: action,
				})
				if err != nil {
					f.sender.SendMessage(&fsmErrorMsg{
						l:   f,
						fsm: f,
						err: err,
					})
				}
				f.sender.SendMessage(&stateEnteredMsg{
					l:     f,
					state: next,
				})
				return next, err
			}()

			origin := f.curState
			f.curState = nextState

			func(origin, next State) {
				f.lock.Unlock()
				defer f.lock.Lock()
				f.sender.SendMessage(&stateExitedMsg{
					l:     f,
					state: origin,
				})
				f.sender.SendMessage(&stateChangedMsg{
					l:    f,
					from: origin,
					to:   next,
				})
			}(origin, nextState)

			return err
		} else {
			func() {
				f.lock.Unlock()
				defer f.lock.Lock()
				f.sender.SendMessage(&eventNotAcceptedMsg{
					l:     f,
					event: event,
				})
			}()

			//log.Printf("no action found, state: %v event: %v\n", f.curState, event)
			return nil
		}
	} else {
		func() {
			f.lock.Unlock()
			defer f.lock.Lock()
			f.sender.SendMessage(&eventNotAcceptedMsg{
				l:     f,
				event: event,
			})
		}()
	}
	return nil
}

func (f *SimpleFSM) SendEvent(event Event, param interface{}) error {
	return f.Execute(event, param)
}


func (l *SimpleFSM) StateChanged(from, to State) {
	for _, v := range l.listeners {
		v.StateChanged(from, to)
	}
}

func (l *SimpleFSM) StateEntered(state State) {
	for _, v := range l.listeners {
		v.StateEntered(state)
	}
}

func (l *SimpleFSM) StateExited(state State) {
	for _, v := range l.listeners {
		v.StateExited(state)
	}
}

func (l *SimpleFSM) EventNotAccepted(event Event) {
	for _, v := range l.listeners {
		v.EventNotAccepted(event)
	}
}

func (l *SimpleFSM) TransitionStarted(action Action) {
	for _, v := range l.listeners {
		v.TransitionStarted(action)
	}
}

func (l *SimpleFSM) TransitionEnded(action Action) {
	for _, v := range l.listeners {
		v.TransitionEnded(action)
	}
}

func (l *SimpleFSM) FSMStarted(fsm FSM) {
	for _, v := range l.listeners {
		v.FSMStarted(fsm)
	}
}

func (l *SimpleFSM) FSMStopped(fsm FSM) {
	for _, v := range l.listeners {
		v.FSMStopped(fsm)
	}
}

func (l *SimpleFSM) FSMError(fsm FSM, err error) {
	for _, v := range l.listeners {
		v.FSMError(fsm, err)
	}
}
