package gron

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const sock = "/tmp/gron.sock"

type Job struct {
	RawCommand string
	RawPrio    int
}

func (j *Job) Encode() []byte {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	encoder.Encode(j.RawCommand)
	encoder.Encode(j.RawPrio)
	return w.Bytes()
}

func (j *Job) Decode(buf []byte) {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	decoder.Decode(&j.RawCommand)
	decoder.Decode(&j.RawPrio)
}

func Server() {
	l, err := net.Listen("unix", sock)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	log.Printf("Starting gron as a daemon")

	ksignal := make(chan os.Signal, 1)

	//Finish the application gracefully
	signal.Notify(ksignal, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		for {
			fd, err := l.Accept()
			if err != nil {
				log.Fatal("accept error:", err)
			}
			go execute(fd)
		}
	}()

	kv := <-ksignal
	log.Printf("Signal to finish : %s", kv.String())
	os.Remove(sock)
}

func execute(c net.Conn) {
	for {
		buf := make([]byte, 1024)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}
		data := buf[0:nr]
		j := new(Job)
		j.Decode(data)

		log.Printf("let me run this : %s with prio %d", j.RawCommand, j.RawPrio)
		//_, err = c.rite(data)
	}
}
