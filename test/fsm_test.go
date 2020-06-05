// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package test

import (
	"github.com/xfali/fsm"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

const (
	state_a = "a"
	state_b = "b"
	state_c = "c"
	state_d = "d"
)

const (
	event_a = "1"
	event_b = "2"
	event_c = "3"
	event_d = "4"
)

func NewFSM(t *testing.T) fsm.FSM {
	m := fsm.New()
	m.AddListener(&fsm.DefaultListener{})
	//m := fsm.New(fsm.SetMessageBufferSize(0))
	m.AddState(state_a, event_b, func(i interface{}) (state fsm.State, e error) {
		t.Log("a -> 2")
		return state_b, nil
	})
	m.AddState(state_a, event_c, func(i interface{}) (state fsm.State, e error) {
		t.Log("a -> 3")
		return state_c, nil
	})
	m.AddState(state_a, event_d, func(i interface{}) (state fsm.State, e error) {
		t.Log("a -> 4")
		return state_d, nil
	})
	m.AddState(state_b, event_a, func(i interface{}) (state fsm.State, e error) {
		t.Log("b -> 1")
		return state_a, nil
	})
	m.AddState(state_b, event_c, func(i interface{}) (state fsm.State, e error) {
		t.Log("b -> 3")
		return state_c, nil
	})
	m.AddState(state_b, event_d, func(i interface{}) (state fsm.State, e error) {
		t.Log("b -> 4")
		return state_d, nil
	})
	m.AddState(state_c, event_a, func(i interface{}) (state fsm.State, e error) {
		t.Log("c -> 1")
		return state_a, nil
	})
	m.AddState(state_c, event_b, func(i interface{}) (state fsm.State, e error) {
		t.Log("c -> 2")
		return state_b, nil
	})
	m.AddState(state_c, event_d, func(i interface{}) (state fsm.State, e error) {
		t.Log("c -> 4")
		return state_d, nil
	})
	m.AddState(state_d, event_a, func(i interface{}) (state fsm.State, e error) {
		t.Log("d -> 1")
		return state_a, nil
	})
	m.AddState(state_d, event_b, func(i interface{}) (state fsm.State, e error) {
		t.Log("d -> 2")
		return state_b, nil
	})
	m.AddState(state_d, event_c, func(i interface{}) (state fsm.State, e error) {
		t.Log("c -> 3")
		return state_c, nil
	})

	m.Initial(state_a)
	return m
}

func TestDefaultFsm(t *testing.T) {
	m := NewFSM(t)

	m.AddListener(&fsm.DefaultListener{})
	m.Start()
	defer m.Close()

	m.SendEvent(event_a, nil)
	m.SendEvent(event_b, nil)
	m.SendEvent(event_c, nil)
	m.SendEvent(event_d, nil)

	m.SendEvent(event_a, nil)
	m.SendEvent(event_b, nil)
	m.SendEvent(event_c, nil)
	m.SendEvent(event_d, nil)

	<-time.NewTimer(10 * time.Second).C
}

func TestDefaultFsmPresister(t *testing.T) {
	m := NewFSM(t)

	m.Start()
	defer m.Close()

	p := fsm.NewFilePersister("./store")
	m.SendEvent(event_b, nil)
	m.SendEvent(event_c, nil)
	m.SendEvent(event_d, nil)

	time.Sleep(time.Second)

	err := p.Save(m)
	if err != nil {
		t.Fatal(err)
	}

	err = p.Restore(m)
	if err != nil {
		t.Fatal(err)
	}

	if *m.Current() != state_d {
		t.Fatal("m. not state_d")
	}
}

func TestDefaultFsmAsync(t *testing.T) {
	m := NewFSM(t)

	m.Start()
	defer m.Close()

	random := rand.New(rand.NewSource(time.Now().Unix()))
	wait := sync.WaitGroup{}
	wait.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wait.Done()

			for j:=0; j<100; j++ {
				event := strconv.Itoa(random.Intn(5))
				m.SendEvent(event, nil)
			}
		}()
	}

	wait.Wait()

	time.Sleep(time.Second)

}
