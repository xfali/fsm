// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package fsm

type Persister interface {
	Restore(FSM) error
	Save(FSM) error
}
