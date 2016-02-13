package gron

import (
	"bytes"
	"encoding/gob"
	"net"
)

type ClientRequest struct {
	Request string
	Object  interface{}
}

func NewClientRequest() *ClientRequest {
	gob.Register(Job{})
	gob.Register(ClientRequest{})
	return new(ClientRequest)
}

func (cr *ClientRequest) Encode() []byte {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	encoder.Encode(cr)
	return w.Bytes()
}

func (cr *ClientRequest) Decode(conn net.Conn) {
	decoder := gob.NewDecoder(conn)
	decoder.Decode(&cr)
}
