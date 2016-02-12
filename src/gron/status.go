package gron

import (
	"bytes"
	"encoding/gob"
	"net"
)


type Status struct {
	Process 	int
	Running		int
	Sequence	int
	Waiting     interface{}
}

func NewStatus() *Status {
	gob.Register([]*Job{})
	gob.Register(Status{})
	return new(Status)
}

func (s *Status) Encode() []byte {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	encoder.Encode(s)
	return w.Bytes()
}

func (s *Status) Decode(conn net.Conn) {
	decoder := gob.NewDecoder(conn)
	decoder.Decode(&s)
}