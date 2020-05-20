// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package fsm

import (
	"fmt"
	"os"
)

const (
	DefaultEventBufferSize   = 1024
	DefaultMessageBufferSize = 1024
)

type eventEntity struct {
	event Event
	param interface{}
}

type DefaultFSM struct {
	SimpleFSM

	eventChanSize int
	msgChanSize   int
	eventChan     chan eventEntity
	msgChan       chan Message
	stopChan      chan bool
}

type Opt func(f *DefaultFSM)

func New(opts ...Opt) *DefaultFSM {
	ret := &DefaultFSM{
		SimpleFSM:   *NewSimpleFSM(),
		stopChan:    make(chan bool),
		msgChanSize: DefaultMessageBufferSize,
	}
	for i := range opts {
		opts[i](ret)
	}
	if ret.eventChanSize <= 0 {
		ret.eventChanSize = DefaultEventBufferSize
	}
	ret.eventChan = make(chan eventEntity, ret.eventChanSize)
	if ret.msgChanSize > 0 {
		ret.msgChan = make(chan Message, ret.msgChanSize)
		ret.sender = ret
	}

	return ret
}

//配置状态机事件缓存大小，默认为1024
func SetEventBufferSize(size int) Opt {
	return func(f *DefaultFSM) {
		f.eventChanSize = size
	}
}

//配置状态机内部消息缓存大小，默认为1024。如果<=0则使用同步发送消息的方式
func SetMessageBufferSize(size int) Opt {
	return func(f *DefaultFSM) {
		f.msgChanSize = size
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
	if f.msgChanSize > 0 {
		go func() {
			for {
				select {
				case msg := <-f.msgChan:
					f.handleMsg(msg)
					break
				case <-f.stopChan:
					return
				}
			}
		}()
	}
	f.sender.SendMessage(&fsmStartedMsg{
		l:   f,
		fsm: f,
	})
	return nil
}

func (f *DefaultFSM) Close() error {
	close(f.stopChan)
	f.sender.SendMessage(&fsmStoppedMsg{
		l:   f,
		fsm: f,
	})
	return nil
}

func (f *DefaultFSM) SendMessage(m Message) {
	f.msgChan <- m
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
