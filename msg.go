// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package fsm

const (
	UnknownMsgType = iota
	stateChangedMsgType
	stateEnteredMsgType
	stateExitedMsgType
	eventNotAcceptedMsgType
	transitionMsgType
	transitionStartedMsgType
	transitionEndedMsgType
	fsmStartedMsgType
	fsmStoppedMsgType
	fsmErrorMsgType

	CustomerMsgType = 10000
)

type Message interface {
	Type() int

	Proc() error
}

type MessageSender interface {
	SendMessage(m Message)
}

type stateChangedMsg struct {
	l        Listener
	from, to State
}

func (m *stateChangedMsg) Type() int {
	return stateChangedMsgType
}

func (m *stateChangedMsg) Proc() error {
	m.l.StateChanged(m.from, m.to)
	return nil
}

type stateEnteredMsg struct {
	l     Listener
	state State
}

func (m *stateEnteredMsg) Type() int {
	return stateEnteredMsgType
}

func (m *stateEnteredMsg) Proc() error {
	m.l.StateEntered(m.state)
	return nil
}

type stateExitedMsg struct {
	l     Listener
	state State
}

func (m *stateExitedMsg) Type() int {
	return stateExitedMsgType
}

func (m *stateExitedMsg) Proc() error {
	m.l.StateExited(m.state)
	return nil
}

type eventNotAcceptedMsg struct {
	l     Listener
	event Event
}

func (m *eventNotAcceptedMsg) Type() int {
	return eventNotAcceptedMsgType
}

func (m *eventNotAcceptedMsg) Proc() error {
	m.l.EventNotAccepted(m.event)
	return nil
}

type transitionMsg struct {
	l      Listener
	action Action
}

func (m *transitionMsg) Type() int {
	return transitionMsgType
}

func (m *transitionMsg) Proc() error {
	panic("not support")
	return nil
}

type transitionStartedMsg struct {
	l      Listener
	action Action
}

func (m *transitionStartedMsg) Type() int {
	return transitionStartedMsgType
}

func (m *transitionStartedMsg) Proc() error {
	m.l.TransitionStarted(m.action)
	return nil
}

type transitionEndedMsg struct {
	l      Listener
	action Action
}

func (m *transitionEndedMsg) Type() int {
	return transitionEndedMsgType
}

func (m *transitionEndedMsg) Proc() error {
	m.l.TransitionEnded(m.action)
	return nil
}

type fsmStartedMsg struct {
	l   Listener
	fsm FSM
}

func (m *fsmStartedMsg) Type() int {
	return fsmStartedMsgType
}

func (m *fsmStartedMsg) Proc() error {
	m.l.FSMStarted(m.fsm)
	return nil
}

type fsmStoppedMsg struct {
	l   Listener
	fsm FSM
}

func (m *fsmStoppedMsg) Type() int {
	return fsmStoppedMsgType
}

func (m *fsmStoppedMsg) Proc() error {
	m.l.FSMStopped(m.fsm)
	return nil
}

type fsmErrorMsg struct {
	l   Listener
	fsm FSM
	err error
}

func (m *fsmErrorMsg) Type() int {
	return fsmErrorMsgType
}

func (m *fsmErrorMsg) Proc() error {
	m.l.FSMError(m.fsm, m.err)
	return nil
}
