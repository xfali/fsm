// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package fsm

import (
	"encoding/json"
	"io/ioutil"
)

type FilePersister struct {
	path string
}

func NewFilePersister(path string) *FilePersister {
	return &FilePersister{
		path: path,
	}
}

func (p *FilePersister) Restore(fsm FSM) error {
	d, err := ioutil.ReadFile(p.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(d, fsm.Current())
}

func (p *FilePersister) Save(fsm FSM) error {
	d, err := json.Marshal(fsm.Current())
	if err != nil {
		return err
	}
	return ioutil.WriteFile(p.path, d, 0660)
}
