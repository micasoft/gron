package gron

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"log"
	"net"
)

type ClientRequest struct {
	Request string
	Job     Job
}

func (cr *ClientRequest) Encode() []byte {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	encoder.Encode(cr.Request)
	encoder.Encode(cr.Job)
	return w.Bytes()
}

func (cr *ClientRequest) Decode(buf []byte) {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	decoder.Decode(&cr.Request)
	decoder.Decode(&cr.Job)
}

func connect() net.Conn {
	c, err := net.Dial("unix", sock)
	if err != nil {
		panic(err)
	}
	return c
}

func Client(cmd *string, prio *int) {
	c := connect()
	defer c.Close()
	bcr := ClientRequest{Request: "job", Job: Job{RawCommand: *cmd, RawPrio: *prio}}
	_, err := c.Write(bcr.Encode())
	if err != nil {
		log.Fatal("write error:", err)
	}
}

func Status() {
	c := connect()
	defer c.Close()
	bcr := ClientRequest{Request: "status"}
	_, err := c.Write(bcr.Encode())
	if err != nil {
		log.Fatal("write error:", err)
	} else {
		result, _ := ioutil.ReadAll(c)
		log.Println(result)
	}
}
